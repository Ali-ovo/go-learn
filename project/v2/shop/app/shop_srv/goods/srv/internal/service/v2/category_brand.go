package srv

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/log"
)

type CategoryBrandService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (cbs *CategoryBrandService) List(ctx context.Context, request *goods_pb.CategoryBrandFilterRequest) (*dto.CategoryBrandDTOList, error) {
	var ret dto.CategoryBrandDTOList
	page := metav1.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums),
	}

	categoryBrandList, err := cbs.data.CategoryBrands().List(ctx, page, request.Orderby)
	if err != nil {
		return nil, err
	}
	ret.TotalCount = categoryBrandList.TotalCount
	for _, value := range categoryBrandList.Items {
		ret.Items = append(ret.Items, &dto.CategoryBrandDTO{
			CategoryBrandDO: *value,
		})
	}
	return &ret, nil
}

func (cbs *CategoryBrandService) Get(ctx context.Context, categoryID int64) (*dto.CategoryBrandDTOList, error) {
	var ret dto.CategoryBrandDTOList
	var err error

	if _, err = cbs.data.Category().Get(ctx, categoryID); err != nil {
		return nil, err
	}

	categoryBrandList, err := cbs.data.CategoryBrands().GetBrandList(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	ret.TotalCount = categoryBrandList.TotalCount
	for _, value := range categoryBrandList.Items {
		ret.Items = append(ret.Items, &dto.CategoryBrandDTO{
			CategoryBrandDO: *value,
		})
	}
	return &ret, nil
}

func (cbs *CategoryBrandService) Create(ctx context.Context, categoryBrand *dto.CategoryBrandDTO) (int64, error) {
	var err error
	if _, err = cbs.data.Category().Get(ctx, int64(categoryBrand.CategoryID)); err != nil {
		return 0, err
	}
	if _, err = cbs.data.Brands().Get(ctx, int64(categoryBrand.BrandsID)); err != nil {
		return 0, err
	}

	if err = cbs.data.CategoryBrands().Create(ctx, nil, &categoryBrand.CategoryBrandDO); err != nil {
		log.Errorf("data.Create err: %v", err)
		return 0, err
	}
	return int64(categoryBrand.ID), nil
}

func (cbs *CategoryBrandService) Update(ctx context.Context, categoryBrand *dto.CategoryBrandDTO) error {
	var err error
	var categoryBrandDO *do.CategoryBrandDO

	if categoryBrandDO, err = cbs.data.CategoryBrands().Get(ctx, int64(categoryBrand.ID)); err != nil {
		return err
	}
	if _, err = cbs.data.Category().Get(ctx, int64(categoryBrand.CategoryID)); err != nil {
		return err
	}
	if _, err = cbs.data.Brands().Get(ctx, int64(categoryBrand.BrandsID)); err != nil {
		return err
	}

	categoryBrandDO.CategoryID = categoryBrand.CategoryID
	categoryBrandDO.BrandsID = categoryBrand.BrandsID

	if err = cbs.data.CategoryBrands().Update(ctx, nil, categoryBrandDO); err != nil {
		//log.Errorf("data.Update err: %v", err)
		return err
	}
	return nil
}

func (cbs *CategoryBrandService) Delete(ctx context.Context, id int64) error {
	if err := cbs.data.Category().Delete(ctx, nil, id); err != nil {
		log.Errorf("data.Delete err: %v", err)
		return err
	}
	return nil
}

func newCategoryBrand(srv *serviceFactory) service.CategoryBrandSrv {
	return &CategoryBrandService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
