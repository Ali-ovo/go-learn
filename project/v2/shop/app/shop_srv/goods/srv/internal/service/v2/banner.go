package srv

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type bannerService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (bs *bannerService) List(ctx context.Context) (*dto.BannerDTOList, error) {
	var ret dto.BannerDTOList

	bannerList, err := bs.data.Banner().List(ctx)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
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
	if result := bs.data.Banner().Create(ctx, nil, &branner.BannerDO); result.RowsAffected == 0 {
		log.Errorf("data.Create err: %v", result.Error)
		if result.Error != nil {
			return 0, errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return 0, errors.WithCode(code2.ErrBannerNotFound, "Create banner failure")
	}
	return branner.ID, nil
}

func (bs *bannerService) Update(ctx context.Context, branner *dto.BannerDTO) error {
	var err error
	var result *gorm.DB
	var brannerDO *do.BannerDO

	if brannerDO, err = bs.data.Banner().Get(ctx, branner.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrBannerNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	brannerDO.Image = branner.Image
	brannerDO.Url = branner.Url
	brannerDO.Index = branner.Index

	if result = bs.data.Banner().Update(ctx, nil, brannerDO); result.RowsAffected == 0 {
		log.Errorf("data.Update err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrBannerNotFound, "Update banner failure")
	}
	return nil
}

func (bs *bannerService) Delete(ctx context.Context, id int64) error {
	if result := bs.data.Banner().Delete(ctx, nil, id); result.RowsAffected == 0 {
		log.Errorf("data.Delete err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrBannerNotFound, "Delete banner failure")
	}
	return nil
}

func newBanner(srv *serviceFactory) service.BannerSrv {
	return &bannerService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
