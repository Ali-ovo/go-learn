package mapstructure

import (
	"github.com/mitchellh/mapstructure"
)

type DecoderOption func(*mapstructure.DecoderConfig)

func NewDecoder(opts ...DecoderOption) (*mapstructure.Decoder, error) {
	decoderConfig := &mapstructure.DecoderConfig{}

	for _, opt := range opts {
		opt(decoderConfig)
	}

	return mapstructure.NewDecoder(decoderConfig)
}

// WithResult 用于指定解码结果的类型，通常是一个结构体或者一个指向结构体的指针
func WithResult(result interface{}) DecoderOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.Result = result
	}
}

// WithDecodeHookFunc 用于定义一个钩子函数，可以在解码时对字段值进行转换或校验。
func WithDecodeHookFunc(fun *mapstructure.DecoderConfig) DecoderOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.DecodeHook = fun
	}
}

// WithTagName 用于定义 struct 标签的名称，默认为 mapstructure
func WithTagName(tagName string) DecoderOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.TagName = tagName
	}
}

// WithWeaklyTypedInput 用于控制是否允许弱类型输入, 如果启用则会将输入值尽可能转换为目标类型，否则会尝试匹配完全相同的类型
func WithWeaklyTypedInput() DecoderOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.WeaklyTypedInput = true
	}
}

// WithZeroFields 用于控制在解码时是否对零值字段进行赋值。如果启用，则会将输入中的零值赋给结构体字段，否则会保留结构体定义时的默认值
func WithZeroFields() DecoderOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.ZeroFields = true
	}
}

func Decode(input interface{}, output interface{}) error {
	decoder, err := NewDecoder(
		WithResult(output),
		WithWeaklyTypedInput(),
	)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
