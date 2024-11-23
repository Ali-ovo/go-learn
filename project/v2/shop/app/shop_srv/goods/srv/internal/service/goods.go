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
	Update(ctx context.Context, id int64, goods *dto.GoodsDTO) error
	// Delete 删除商品
	Delete(ctx context.Context, id int64) error
	// BatchGet 批量查询商品
	BatchGet(ctx context.Context, ids []int64) ([]*dto.GoodsDTO, error)
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
	goodsIDs := []uint32{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint32(value.ID))
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

func (gs *goodsService) Get(ctx context.Context, id uint64) (*dto.GoodsDTO, error) {
	good, err := gs.goodsData.Get(ctx, id)
	if err != nil {
		log.ErrorfC(ctx, "data.Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *good,
	}, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	if _, err := gs.brandsData.Get(ctx, goods.BrandsID); err == nil {
		return err
	}
	if _, err := gs.categoryData.Get(ctx, goods.CategoryID); err != nil {
		return err
	}

	if err := gs.goodsData.Create(ctx, &goods.GoodsDO); err != nil {
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
	if err := gs.goodsSerachData.Create(ctx, &goodsSearchDo); err != nil {
		return err
	}
	return nil
}

func (gs *goodsService) Update(ctx context.Context, id int64, goods *dto.GoodsDTO) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []int64) ([]*dto.GoodsDTO, error) {
	//TODO implement me
	panic("implement me")
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
