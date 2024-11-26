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
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	categoryBrandList, err := cbs.data.CategoryBrands().GetBrandList(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryBrandsNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
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
	var result *gorm.DB
	if _, err = cbs.data.Category().Get(ctx, categoryBrand.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}
	if _, err = cbs.data.Brands().Get(ctx, categoryBrand.BrandsID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.WithCode(code2.ErrBrandsNotFound, err.Error())
		}
		return 0, errors.WithCode(code.ErrDatabase, err.Error())
	}

	if result = cbs.data.CategoryBrands().Create(ctx, nil, &categoryBrand.CategoryBrandDO); result.RowsAffected == 0 {
		log.Errorf("data.Create err: %v", result.Error)
		if result.Error != nil {
			return 0, errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return 0, errors.WithCode(code2.ErrCategoryBrandsNotFound, "Create CategoryBrands failed")
	}
	return categoryBrand.ID, nil
}

func (cbs *CategoryBrandService) Update(ctx context.Context, categoryBrand *dto.CategoryBrandDTO) error {
	var err error
	var result *gorm.DB
	var categoryBrandDO *do.CategoryBrandDO

	if categoryBrandDO, err = cbs.data.CategoryBrands().Get(ctx, categoryBrand.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrCategoryBrandsNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	if _, err = cbs.data.Category().Get(ctx, categoryBrand.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	if _, err = cbs.data.Brands().Get(ctx, categoryBrand.BrandsID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrBrandsNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	categoryBrandDO.CategoryID = categoryBrand.CategoryID
	categoryBrandDO.BrandsID = categoryBrand.BrandsID

	if result = cbs.data.CategoryBrands().Update(ctx, nil, categoryBrandDO); result.RowsAffected == 0 {
		log.Errorf("data.Update err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrCategoryBrandsNotFound, "Update CategoryBrands failed")
	}
	return nil
}

func (cbs *CategoryBrandService) Delete(ctx context.Context, id int64) error {
	if result := cbs.data.CategoryBrands().Delete(ctx, nil, id); result.RowsAffected == 0 {
		log.Errorf("data.Delete err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrCategoryBrandsNotFound, "Delete CategoryBrands failed")
	}
	return nil
}

func newCategoryBrand(srv *serviceFactory) service.CategoryBrandSrv {
	return &CategoryBrandService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
