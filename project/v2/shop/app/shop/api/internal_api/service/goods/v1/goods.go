package serviceGoods

import (
	"context"
	"encoding/json"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_api/data/v1"
	dtoGoods "shop/app/shop/api/internal_api/domain/dto/goods"
	srvGoods "shop/app/shop/api/internal_api/service/goods"
	code2 "shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
)

type GoodsService struct {
	data data.DataFactory
}

func (gs *GoodsService) GoodsList(ctx context.Context, req *dtoGoods.GoodsFilter) (*dtoGoods.GoodDTOList, error) {
	var dtoGoodsList dtoGoods.GoodDTOList
	var doGoodsFilter goods_pb.GoodsFilterRequest

	byteReq, _ := json.Marshal(req)
	if err := json.Unmarshal(byteReq, &doGoodsFilter); err != nil {
		return nil, errors.WithCode(code2.ErrEncodingJSON, "编码 json数据 失败")
	}

	doGoods, err := gs.data.Goods().GoodsList(ctx, &doGoodsFilter)
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	byteDoGoods, _ := json.Marshal(doGoods)
	if err = json.Unmarshal(byteDoGoods, &dtoGoodsList); err != nil {
		return nil, errors.WithCode(code2.ErrDecodingJSON, "解码 json数据 失败")
	}
	return &dtoGoodsList, nil
}

func (gs *GoodsService) BatchGetGoods(ctx context.Context, req *goods_pb.BatchGoodsIdInfo) (*goods_pb.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *GoodsService) CreateGoods(ctx context.Context, req *goods_pb.CreateGoodsInfo) (*goods_pb.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *GoodsService) DeleteGoods(ctx context.Context, req *goods_pb.DeleteGoodsInfo) error {
	//TODO implement me
	panic("implement me")
}

func (gs *GoodsService) UpdateGoods(ctx context.Context, req *goods_pb.CreateGoodsInfo) error {
	//TODO implement me
	panic("implement me")
}

func (gs *GoodsService) GetGoodsDetail(ctx context.Context, req *goods_pb.GoodInfoRequest) (*goods_pb.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewGoodsService(data data.DataFactory) srvGoods.GoodsSrv {
	return &GoodsService{data: data}
}
