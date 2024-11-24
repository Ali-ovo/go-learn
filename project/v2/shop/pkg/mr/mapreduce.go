package mr

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

/*
	此代码复制于 go-zero https://go-zero.dev/cn/docs/blog/concurrency/mapreduce/
	用于 并发处理以降低服务响应时间
*/

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was cancelled with nil.
	ErrCancelWithNil = errors.New("mapreduce cancelled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	// ForEachFunc is used to do element processing, but no output.
	ForEachFunc[T any] func(item T)
	// GenerateFunc is used to let callers send elements into source.
	// source: 一个容量为0的channel 逐渐收集我们所需要的数据
	GenerateFunc[T any] func(source chan<- T)
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc[T, U any] func(item T, writer Writer[U])
	// MapperFunc is used to do element processing and write the output to writer,
	// use cancel func to cancel the processing.
	// 根据 item数据 开启多个协程 并将数据写道 writer 中  如果执行过程中错误 将 直接返回报错
	// item: 逐渐接收这个 channel 所给的数据
	// writer: 执行的结果会写入到此结构体中
	// cancal: 执行结果后 or 执行错误 此函数会被执行到
	MapperFunc[T, U any] func(item T, writer Writer[U], cancel func(error))
	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	// pipe: 函数二的 Writer guardedWriter.channel 的数据会存储 pipe 中
	ReducerFunc[U, V any] func(pipe <-chan U, writer Writer[V], cancel func(error))
	// VoidReducerFunc is used to reduce all the mapping output, but no output.
	// Use cancel func to cancel the processing.
	VoidReducerFunc[U any] func(pipe <-chan U, cancel func(error))
	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapperContext[T, U any] struct {
		ctx       context.Context
		mapper    MapFunc[T, U]
		source    <-chan T
		panicChan *onceChan
		collector chan<- U
		doneChan  <-chan struct{}
		workers   int
	}

	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}

	// Writer interface wraps Write method.
	Writer[T any] interface {
		Write(v T)
	}
)

// Finish 并行运行 fns 一旦出现 errors 就取消
func Finish(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return MapReduceVoid(
		func(source chan<- func() error) {
			for _, fn := range fns {
				source <- fn
			}
		},
		func(fn func() error, writer Writer[any], cancel func(error)) {
			if err := fn(); err != nil {
				cancel(err)
			}
		},
		func(pipe <-chan any, cancel func(error)) {
		},
		WithWorkers(len(fns)), // 设置 最大并发数
	)
}

// FinishVoid runs fns parallelly.
func FinishVoid(fns ...func()) {
	if len(fns) == 0 {
		return
	}

	ForEach(func(source chan<- func()) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(fn func()) {
		fn()
	}, WithWorkers(len(fns)))
}

// ForEach maps all elements from given generate but no output.
func ForEach[T any](generate GenerateFunc[T], mapper ForEachFunc[T], opts ...Option) {
	options := buildOptions(opts...)
	panicChan := &onceChan{channel: make(chan any)}
	source := buildSource(generate, panicChan)
	collector := make(chan any)
	done := make(chan struct{})

	go executeMappers(mapperContext[T, any]{
		ctx: options.ctx,
		mapper: func(item T, _ Writer[any]) {
			mapper(item)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	for {
		select {
		case v := <-panicChan.channel:
			panic(v)
		case _, ok := <-collector:
			if !ok {
				return
			}
		}
	}
}

// MapReduce maps all elements generated from given generate func,
// and reduces the output elements with given reducer.
func MapReduce[T, U, V any](generate GenerateFunc[T], mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (V, error) {
	// 用于 存储 错误数据
	panicChan := &onceChan{channel: make(chan any)}
	// 遍历需要被执行的函数 传递到 source 中
	source := buildSource(generate, panicChan)
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// MapReduceChan maps all elements from source, and reduce the output elements with given reducer.
func MapReduceChan[T, U, V any](source <-chan T, mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (V, error) {
	// 用于 存储 错误数据
	panicChan := &onceChan{channel: make(chan any)}
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// mapReduceWithPanicChan maps all elements from source, and reduce the output elements with given reducer.
func mapReduceWithPanicChan[T, U, V any](source <-chan T, panicChan *onceChan, mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (val V, err error) {
	options := buildOptions(opts...)
	// output 写入最终的数据
	output := make(chan V)
	defer func() {
		// reducer can only write once, if more, panic
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	// 收集器用于从 mapper 收集数据，并在 reducer 中使用
	// 创建 channel 根据 设置的最大来定
	collector := make(chan U, options.workers)
	// if done is closed, all mappers and reducer should stop processing
	done := make(chan struct{})
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once
	// use atomic type to avoid data race
	var retErr AtomicError
	// 定义 一个 用于关闭 output done channel 的函数
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	// 创建一个用于关闭 的函数
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}
		// 用于清空 channel 的下水道 避免内存泄露
		drain(source)
		finish()
	})

	// PS 为什么第三个函数reducer 需要比 第二个函数mapper 先被执行
	// 因为 这个函数需要接收 第二个函数传递 channel 数据  这个channel是无容量的 不这样写 会被hold住
	go func() {
		defer func() {
			drain(collector)
			if r := recover(); r != nil {
				// 写入错误信息
				panicChan.write(r)
			}
			finish()
		}()
		// collector: channel列表根据最大并发量来计算
		// writer: 被写入的数据流
		// cancel: 调用结束或失败 时候 关闭用
		reducer(collector, writer, cancel)
	}()

	go executeMappers(mapperContext[T, U]{
		ctx: options.ctx,
		mapper: func(item T, w Writer[U]) {
			mapper(item, w, cancel)
		},
		source:    source,    // 已经请求到的数据
		panicChan: panicChan, // 错误流
		collector: collector, // channel列表根据最大并发量来计算
		doneChan:  done,
		workers:   options.workers, // 最大并发量
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		err = context.DeadlineExceeded
	case v := <-panicChan.channel:
		// drain output here, otherwise for loop panic in defer
		drain(output)
		panic(v)
	case v, ok := <-output:
		if e := retErr.Load(); e != nil {
			err = e
		} else if ok {
			val = v
		} else {
			err = ErrReduceNoOutput
		}
	}

	return
}

// MapReduceVoid  maps all elements generated from given generate, and reduce the output elements with given
//
//	@Description:
//	@param generate
//	@param mapper
//	@param U]
//	@param reducer
//	@param opts
//	@return error
func MapReduceVoid[T, U any](generate GenerateFunc[T], mapper MapperFunc[T, U], reducer VoidReducerFunc[U], opts ...Option) error {
	_, err := MapReduce(
		generate,
		mapper,
		func(input <-chan U, writer Writer[any], cancel func(error)) {
			reducer(input, cancel)
		},
		opts...)
	if errors.Is(err, ErrReduceNoOutput) {
		return nil
	}

	return err
}

// WithContext customizes a mapreduce processing accepts a given ctx.
func WithContext(ctx context.Context) Option {
	return func(opts *mapReduceOptions) {
		opts.ctx = ctx
	}
}

// WithWorkers 设置最大并发量 默认最大并发数 16
//
//	@Description:
//	@param workers
//	@return Option
func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

// 生成 mapReduceOptions 结构体
func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// 开启协程 执行逻辑 并且 recover 错误, 如果有错误 panicChan.write(r) 写入到 panicChan 中
func buildSource[T any](generate GenerateFunc[T], panicChan *onceChan) chan T {
	// 执行的数据 需要存储再 source 中
	source := make(chan T)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()

		// 将 数据 一个个传递到 source 中
		generate(source)
	}()

	return source
}

// 用于清空 channel 的下水道 避免内存泄露
func drain[T any](channel <-chan T) {
	// drain the channel
	for range channel {
	}
}

func executeMappers[T, U any](mCtx mapperContext[T, U]) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(mCtx.collector)
		// 清空 通道的数据
		drain(mCtx.source)
	}()

	var failed int32
	pool := make(chan struct{}, mCtx.workers)
	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	// 如果 failed == 0 循环执行 报错退出
	for atomic.LoadInt32(&failed) == 0 {
		select {
		case <-mCtx.ctx.Done():
			return
		case <-mCtx.doneChan:
			return
		case pool <- struct{}{}: // 存入 空结构
			item, ok := <-mCtx.source // 逐个获取的 source channel 中的数据
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			}()
		}
	}
}

// 创建 默认 mapReduceOptions
func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		ctx:     context.Background(),
		workers: defaultWorkers,
	}
}

// 封装了原生 sync.Once 用于执行一次的方法
func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

type guardedWriter[T any] struct {
	ctx     context.Context
	channel chan<- T
	done    <-chan struct{}
}

func newGuardedWriter[T any](ctx context.Context, channel chan<- T, done <-chan struct{}) guardedWriter[T] {
	return guardedWriter[T]{
		ctx:     ctx,
		channel: channel,
		done:    done,
	}
}

// 写入数据到 gw.channel 中
func (gw guardedWriter[T]) Write(v T) {
	select {
	case <-gw.ctx.Done():
		return
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}

type onceChan struct {
	channel chan any
	wrote   int32
}

// 写入 数据 只执行一次
func (oc *onceChan) write(val any) {
	/*
		判断是否 等于 旧值 0
		true retun true 并 set 新值
		false return false
	*/
	if atomic.CompareAndSwapInt32(&oc.wrote, 0, 1) {
		oc.channel <- val
	}
}
