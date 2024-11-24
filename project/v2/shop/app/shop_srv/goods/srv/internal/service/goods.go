package service

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/gmicro/pkg/log"
)

type GoodsSrv interface {
	// List 商品列表
	List(ctx context.Context, req *goods_pb.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)
	// Get 商品详情
	Get(ctx context.Context, id uint64) (*dto.GoodsDTO, error)
	// Create 创建商品
	Create(ctx context.Context, goods *dto.GoodsDTO) error
	// Update 更新商品
	Update(ctx context.Context, goods *dto.GoodsDTO) error
	// Delete 删除商品
	Delete(ctx context.Context, id uint64) error
	// BatchGet 批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

type goodsService struct {
	goodsData       data.GoodsStore
	goodsSerachData data_search.GoodsStore
	categoryData    data.CategoryStore
	brandsData      data.BrandsStore
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

func (gs *goodsService) List(ctx context.Context, req *goods_pb.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	var ret *dto.GoodsDTOList

	searchReq := data_search.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	if req.TopCategory > 0 {
		category, err := gs.categoryData.Get(ctx, req.TopCategory)
		if err != nil {
			log.ErrorfC(ctx, "categoryData.Get err: %v", err)
			return nil, err
		}
		for _, value := range retrieveIDs(category) {
			searchReq.CategoryIDs = append(searchReq.CategoryIDs, value)
		}
	}
	goodsList, err := gs.goodsSerachData.Search(ctx, &searchReq)
	if err != nil {
		log.ErrorfC(ctx, "goodsSerachData.Search err: %v", err)
		return nil, err
	}

	// 通过 id 批量查询 mysql 数据
	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}
	goodsData, err := gs.goodsData.ListByIDs(ctx, goodsIDs, orderby)
	if err != nil {
		return nil, err
	}
	ret.TotalCount = goodsList.TotalCount
	for _, value := range goodsData.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return ret, nil
}

func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	good, err := gs.goodsData.Get(ctx, ID)
	if err != nil {
		log.ErrorfC(ctx, "data.Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *good,
	}, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	/*
		方案一: 基于可靠消息实现最终一致性 消息队列[事务消息] (更好 代码侵入性强)
		方案二: 基于mysql事务消息 (有一定的风险 ES 超时但是执行了的情况)
		方案三: 基于 阿里云开源的 canal
			读取 mysql binlog文件 并将 数据 分发到 kafka 、 rocketmq 或 hbase 等中间件
		方案四: 创建一张表 专门记录 执行 ES数据失败的信息 (比如记录 数据库的id)
			然后设置一个定时任务, 定时向ES执行 失败的数据的数据 (根据 记录的 ID 从mysql中 读取并且再次执行)
	*/
	var err error
	if _, err = gs.brandsData.Get(ctx, goods.BrandsID); err == nil {
		return err
	}
	if _, err = gs.categoryData.Get(ctx, goods.CategoryID); err != nil {
		return err
	}

	txn := gs.goodsData.Begin(ctx)
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

	if err = gs.goodsData.CreateInTxn(ctx, txn, &goods.GoodsDO); err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
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
	if err = gs.goodsSerachData.Create(ctx, &goodsSearchDo); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	var err error

	txn := gs.goodsData.Begin(ctx)
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

	if err = gs.goodsData.UpdateInTxn(ctx, txn, &goods.GoodsDO); err != nil {
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
	if err = gs.goodsSerachData.Update(ctx, &goodsSearchDo); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) Delete(ctx context.Context, id uint64) error {
	var err error

	txn := gs.goodsData.Begin(ctx)
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

	if err = gs.goodsData.DeleteInTxn(ctx, txn, id); err != nil {
		return err
	}

	if err = gs.goodsSerachData.Delete(ctx, id); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	// 并发去请求数据 一次性启动多个 goroutine 需要去监听
	var ret []*dto.GoodsDTO

	ds, err := gs.goodsData.ListByIDs(ctx, ids, []string{})
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

func NewGoodsService(gs data.GoodsStore, sgs data_search.GoodsStore, cs data.CategoryStore, bs data.BrandsStore) GoodsSrv {
	return &goodsService{
		goodsData:       gs,
		goodsSerachData: sgs,
		categoryData:    cs,
		brandsData:      bs,
	}
}

var _ GoodsSrv = (*goodsService)(nil)
