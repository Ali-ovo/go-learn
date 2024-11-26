package srv

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type brandService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (bs *brandService) List(ctx context.Context, request *goods_pb.BrandFilterRequest) (*dto.BrandsDTOList, error) {
	var ret dto.BrandsDTOList
	page := metav1.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums),
	}

	brandList, err := bs.data.Brands().List(ctx, page, request.Orderby)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.TotalCount = brandList.TotalCount
	for _, value := range brandList.Items {
		ret.Items = append(ret.Items, &dto.BrandsDTO{
			BrandsDO: *value,
		})
	}
	return &ret, nil
}

func (bs *brandService) Create(ctx context.Context, brand *dto.BrandsDTO) (int64, error) {
	if result := bs.data.Brands().Create(ctx, nil, &brand.BrandsDO); result.RowsAffected == 0 {
		log.Errorf("data.Create err: %v", result.Error)
		if result.Error != nil {
			return 0, errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return 0, errors.WithCode(code2.ErrBrandsNotFound, "Create Brands failure")
	}
	return brand.ID, nil
}

func (bs *brandService) Update(ctx context.Context, brand *dto.BrandsDTO) error {
	var err error
	var result *gorm.DB
	var brandDO *do.BrandsDO

	if brandDO, err = bs.data.Brands().Get(ctx, brand.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrBrandsNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	brandDO.Name = brand.Name
	brandDO.Logo = brand.Logo

	if result = bs.data.Brands().Update(ctx, nil, brandDO); result.RowsAffected == 0 {
		log.Errorf("data.Update err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrBrandsNotFound, "Update Brands failure")
	}
	return nil
}

func (bs *brandService) Delete(ctx context.Context, id int64) error {
	if result := bs.data.Brands().Delete(ctx, nil, id); result.RowsAffected == 0 {
		log.Errorf("data.Delete err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrBrandsNotFound, "Delete Brands failure")
	}
	return nil
}

func newBrand(srv *serviceFactory) service.BrandSrv {
	return &brandService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
