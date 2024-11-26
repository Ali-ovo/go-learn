package srv

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	"shop/gmicro/pkg/log"
)

type bannerService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (bs *bannerService) List(ctx context.Context) (*dto.BannerDTOList, error) {
	var ret dto.BannerDTOList

	bannerList, err := bs.data.Banner().List(ctx)
	if err != nil {
		return nil, err
	}
	ret.TotalCount = bannerList.TotalCount
	for _, value := range bannerList.Items {
		ret.Items = append(ret.Items, &dto.BannerDTO{
			BannerDO: *value,
		})
	}
	return &ret, nil
}

func (bs *bannerService) Create(ctx context.Context, branner *dto.BannerDTO) (int64, error) {
	if err := bs.data.Banner().Create(ctx, nil, &branner.BannerDO); err != nil {
		log.Errorf("data.Create err: %v", err)
		return 0, err
	}
	return branner.ID, nil
}

func (bs *bannerService) Update(ctx context.Context, branner *dto.BannerDTO) error {
	var err error
	var brannerDO *do.BannerDO

	if brannerDO, err = bs.data.Banner().Get(ctx, int64(branner.ID)); err != nil {
		return err
	}

	brannerDO.Image = branner.Image
	brannerDO.Url = branner.Url
	brannerDO.Index = branner.Index

	if err = bs.data.Banner().Update(ctx, brannerDO); err != nil {
		//log.Errorf("data.Update err: %v", err)
		return err
	}
	return nil
}

func (bs *bannerService) Delete(ctx context.Context, id int64) error {
	if err := bs.data.Banner().Delete(ctx, id); err != nil {
		log.Errorf("data.Delete err: %v", err)
		return err
	}
	return nil
}

func newBanner(srv *serviceFactory) service.BannerSrv {
	return &bannerService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
