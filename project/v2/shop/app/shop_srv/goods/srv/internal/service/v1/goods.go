package srv

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	"shop/gmicro/pkg/log"
	"shop/pkg/mr"
	"sync"
)

type goodsService struct {
	// 工厂
	data      data.DataFactory
	seachData data_search.SearchFactory
}

// 遍历树结构
func retrieveIDs(category *do.CategoryDO) []uint32 {
	ids := []uint32{}
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint32(category.ID))
	for _, child := range category.SubCategory {
		subids := retrieveIDs(child)
		ids = append(ids, subids...)
	}
	return ids
}

func (gs *goodsService) List(ctx context.Context, req *goods_pb.GoodsFilterRequest) (*dto.GoodsDTOList, error) {
	var ret dto.GoodsDTOList

	searchReq := data_search.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	if req.TopCategory > 0 {
		category, err := gs.data.Category().Get(ctx, int64(req.TopCategory))
		if err != nil {
			log.ErrorfC(ctx, "categoryData.Get err: %v", err)
			return nil, err
		}
		for _, value := range retrieveIDs(category) {
			searchReq.CategoryIDs = append(searchReq.CategoryIDs, value)
		}
	}
	goodsList, err := gs.seachData.Goods().Search(ctx, &searchReq)
	if err != nil {
		log.ErrorfC(ctx, "goodsSerachData.Search err: %v", err)
		return nil, err
	}

	// 通过 id 批量查询 mysql 数据
	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}
	goodsData, err := gs.data.Goods().ListByIDs(ctx, goodsIDs, req.Orderby)
	if err != nil {
		return nil, err
	}
	ret.TotalCount = goodsList.TotalCount
	for _, value := range goodsData.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return &ret, nil
}

func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	good, err := gs.data.Goods().Get(ctx, ID)
	if err != nil {
		log.ErrorfC(ctx, "data.Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *good,
	}, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) (int64, error) {
	/*
		方案一: 基于可靠消息实现最终一致性 消息队列[事务消息] (更好 代码侵入性强)
		方案二: 基于mysql事务消息 (有一定的风险 ES 超时但是执行了的情况)
		方案三: 基于 阿里云开源的 canal
			读取 mysql binlog文件 并将 数据 分发到 kafka 、 rocketmq 或 hbase 等中间件
		方案四: 创建一张表 专门记录 执行 ES数据失败的信息 (比如记录 数据库的id)
			然后设置一个定时任务, 定时向ES执行 失败的数据的数据 (根据 记录的 ID 从mysql中 读取并且再次执行)
	*/
	var err error
	if _, err = gs.data.Brands().Get(ctx, int64(goods.BrandsID)); err != nil {
		return 0, err
	}
	if _, err = gs.data.Category().Get(ctx, int64(goods.CategoryID)); err != nil {
		return 0, err
	}

	txn := gs.data.Begin()
	defer func() { // 异常处理
		if p := recover(); p != nil {
			txn.Rollback()
			log.ErrorfC(ctx, "goodsService.Create panic: %v", p)
			return
		} else if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	if err = gs.data.Goods().Create(ctx, txn, &goods.GoodsDO); err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
		return 0, err
	}

	goodsSearchDo := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief:  goods.GoodsBrief,
		ShopPrice:   goods.ShopPrice,
	}
	if err = gs.seachData.Goods().Create(ctx, &goodsSearchDo); err != nil {
		return 0, err
	}
	txn.Commit()
	return goods.ID, nil
}

func (gs *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	var err error
	var goodDO *do.GoodsDO

	if goodDO, err = gs.data.Goods().Get(ctx, uint64(goods.ID)); err != nil {
		return err
	}

	if goods.BrandsID == 0 {
		goodDO.IsNew = goods.IsNew
		goodDO.IsHot = goods.IsHot
		goodDO.OnSale = goods.OnSale
	} else {
		if _, err = gs.data.Brands().Get(ctx, int64(goods.BrandsID)); err != nil {
			return err
		}
		if _, err = gs.data.Category().Get(ctx, int64(goods.CategoryID)); err != nil {
			return err
		}
		goodDO.BrandsID = goods.BrandsID
		goodDO.CategoryID = goods.CategoryID
		goodDO.Name = goods.Name
		goodDO.GoodsSn = goods.GoodsSn
		goodDO.MarketPrice = goods.MarketPrice
		goodDO.ShopPrice = goods.ShopPrice
		goodDO.GoodsBrief = goods.GoodsBrief
		goodDO.ShipFree = goods.ShipFree
		goodDO.Images = goods.Images
		goodDO.DescImages = goods.DescImages
		goodDO.GoodsFrontImage = goods.GoodsFrontImage
	}

	txn := gs.data.Begin()
	defer func() { // 异常处理
		if p := recover(); p != nil {
			txn.Rollback()
			log.ErrorfC(ctx, "goodsService.Create panic: %v", p)
			return
		} else if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	if err = gs.data.Goods().Update(ctx, txn, goodDO); err != nil {
		return err
	}

	goodsSearchDo := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief:  goods.GoodsBrief,
		ShopPrice:   goods.ShopPrice,
	}
	if err = gs.seachData.Goods().Update(ctx, &goodsSearchDo); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) Delete(ctx context.Context, id uint64) error {
	var err error

	txn := gs.data.Begin()
	defer func() { // 异常处理
		if p := recover(); p != nil {
			txn.Rollback()
			log.ErrorfC(ctx, "goodsService.Create panic: %v", p)
			return
		} else if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	if err = gs.data.Goods().Delete(ctx, txn, id); err != nil {
		return err
	}

	if err = gs.seachData.Goods().Delete(ctx, id); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []int64) ([]*dto.GoodsDTO, error) {
	/*
		go-zero 有个工具MapReduce 可以并发去请求数据 一次性启动多个 goroutine 需要去监听
		方案1: 在底层有 List 的接口前提下优先使用 List 接口 (优先)
		方案2: 在底层没有 List 前提下, 可以通过并发去请求数据 PS:这里实现 并发去请求数据 但不是最优解
	*/
	var callFuncs []func() (*dto.GoodsDTO, error)
	for _, value := range ids {
		// 坑!!! 因为 这里的函数不是在此函数内执行的, 所以 value 会被动态的修改 需要需要在此函数内声明一个全新的变量来使用
		tmp := value
		callFuncs = append(callFuncs, func() (*dto.GoodsDTO, error) {
			goodsDTO, err := gs.Get(ctx, uint64(tmp))
			if err != nil {
				return goodsDTO, err
			}
			return goodsDTO, nil
		})
	}

	ret, err := mr.MapReduce(
		func(source chan<- func() (*dto.GoodsDTO, error)) {
			for _, fn := range callFuncs {
				source <- fn
			}
		},
		func(fn func() (*dto.GoodsDTO, error), writer mr.Writer[*dto.GoodsDTO], cancel func(error)) {
			if goodsDTO, err := fn(); err != nil {
				cancel(err)
			} else {
				writer.Write(goodsDTO)
			}
		},
		func(pipe <-chan *dto.GoodsDTO, writer mr.Writer[[]*dto.GoodsDTO], cancel func(error)) {
			var ret []*dto.GoodsDTO
			for goodsDTO := range pipe {
				ret = append(ret, goodsDTO)
			}
			writer.Write(ret)
		},
		mr.WithWorkers(len(callFuncs)), // 设置 最大并发数
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	//var ret []*dto.GoodsDTO
	//
	//ds, err := gs.goodsData.ListByIDs(ctx, ids, []string{})
	//if err != nil {
	//	return nil, err
	//}
	//for _, value := range ds.Items {
	//	ret = append(ret, &dto.GoodsDTO{
	//		GoodsDO: *value,
	//	})
	//}
	return ret, nil
}

func (gs *goodsService) BatchGetTwe(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	var ret []*dto.GoodsDTO

	ds, err := gs.data.Goods().ListByIDs(ctx, ids, []string{})
	if err != nil {
		return nil, err
	}
	for _, value := range ds.Items {
		ret = append(ret, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return ret, nil
}

func (gs *goodsService) BatchGetThree(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	var ret []*dto.GoodsDTO
	var callFuncs []func() error
	var mu sync.Mutex
	for _, value := range ids {
		// 坑!!! 因为 这里的函数不是在此函数内执行的, 所以 value 会被动态的修改 需要需要在此函数内声明一个全新的变量来使用
		tmp := value
		callFuncs = append(callFuncs, func() error {
			goodsDTO, err := gs.Get(ctx, tmp)
			mu.Lock()
			ret = append(ret, goodsDTO)
			mu.Unlock()
			return err
		})
	}
	err := mr.Finish(callFuncs...)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func newGoods(srv *serviceFactory) service.GoodsSrv {
	return &goodsService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}

var _ service.GoodsSrv = (*goodsService)(nil)
