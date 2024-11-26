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

type CategoryService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (cs *CategoryService) AllList(ctx context.Context) (*dto.CategoryDTOList, error) {
	var ret dto.CategoryDTOList

	categoryList, err := cs.data.Category().ListAll(ctx, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.TotalCount = categoryList.TotalCount
	for _, value := range categoryList.Items {
		ret.Items = append(ret.Items, &dto.CategoryDTO{
			CategoryDO: *value,
		})
	}
	return &ret, nil
}

func (cs *CategoryService) List(ctx context.Context, level int32) (*dto.CategoryDTOList, error) {
	var ret dto.CategoryDTOList

	categoryList, err := cs.data.Category().List(ctx, level)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.TotalCount = categoryList.TotalCount
	for _, value := range categoryList.Items {
		ret.Items = append(ret.Items, &dto.CategoryDTO{
			CategoryDO: *value,
		})
	}
	return &ret, nil
}

func (cs *CategoryService) Get(ctx context.Context, id int64) (*dto.CategoryDTO, error) {
	var ret dto.CategoryDTO

	category, err := cs.data.Category().Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.CategoryDO = *category
	return &ret, nil
}

func (cs *CategoryService) Create(ctx context.Context, category *dto.CategoryDTO) (int64, error) {
	result := cs.data.Category().Create(ctx, nil, &category.CategoryDO)
	if result.RowsAffected == 0 {
		log.Errorf("data.Ceate err: %v", result.Error)
		if result.Error != nil {
			return 0, errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return 0, errors.WithCode(code2.ErrCategoryNotFound, "Create category failure")
	}
	return category.ID, nil
}

func (cs *CategoryService) Update(ctx context.Context, category *dto.CategoryDTO) error {
	var err error
	var result *gorm.DB
	var categoryDO *do.CategoryDO

	if categoryDO, err = cs.data.Category().Get(ctx, category.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	categoryDO.Name = category.Name
	categoryDO.ParentCategoryID = category.ParentCategoryID
	categoryDO.Level = category.Level
	categoryDO.IsTab = category.IsTab

	if result = cs.data.Category().Update(ctx, nil, &category.CategoryDO); result.RowsAffected == 0 {
		//log.Errorf("data.Update err: %v", err)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrCategoryNotFound, "Update category failure")
	}
	return nil
}

func (cs *CategoryService) Delete(ctx context.Context, id int64) error {
	if result := cs.data.Category().Delete(ctx, nil, id); result.RowsAffected == 0 {
		log.Errorf("data.Delete err: %v", result.Error)
		if result.Error != nil {
			return errors.WithCode(code.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code2.ErrCategoryNotFound, "Delete category failure")
	}
	return nil
}

func newCategory(srv *serviceFactory) service.CategorySrv {
	return &CategoryService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
