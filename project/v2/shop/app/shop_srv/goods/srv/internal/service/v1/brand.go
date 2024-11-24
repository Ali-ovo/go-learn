package srv

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/log"
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
		return nil, err
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
	if err := bs.data.Brands().Create(ctx, &brand.BrandsDO); err != nil {
		log.Errorf("data.Create err: %v", err)
		return 0, err
	}
	return brand.ID, nil
}

func (bs *brandService) Update(ctx context.Context, brand *dto.BrandsDTO) error {
	var err error
	var brandDO *do.BrandsDO

	if brandDO, err = bs.data.Brands().Get(ctx, int64(brand.ID)); err != nil {
		return err
	}

	brandDO.Name = brand.Name
	brandDO.Logo = brand.Logo

	if err = bs.data.Brands().Update(ctx, brandDO); err != nil {
		//log.Errorf("data.Update err: %v", err)
		return err
	}
	return nil
}

func (bs *brandService) Delete(ctx context.Context, id int64) error {
	if err := bs.data.Brands().Delete(ctx, id); err != nil {
		log.Errorf("data.Delete err: %v", err)
		return err
	}
	return nil
}

func newBrand(srv *serviceFactory) service.BrandSrv {
	return &brandService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
