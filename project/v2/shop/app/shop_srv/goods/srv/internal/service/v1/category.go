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

type CategoryService struct {
	data      data.DataFactory
	seachData data_search.SearchFactory
}

func (cs *CategoryService) AllList(ctx context.Context) (*dto.CategoryDTOList, error) {
	var ret dto.CategoryDTOList

	categoryList, err := cs.data.Category().ListAll(ctx, []string{})
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	ret.CategoryDO = *category
	return &ret, nil
}

func (cs *CategoryService) Create(ctx context.Context, category *dto.CategoryDTO) (int64, error) {
	if err := cs.data.Category().Create(ctx, &category.CategoryDO); err != nil {
		log.Errorf("data.Create err: %v", err)
		return 0, err
	}
	return int64(category.ID), nil
}

func (cs *CategoryService) Update(ctx context.Context, category *dto.CategoryDTO) error {
	var err error
	var categoryDO *do.CategoryDO

	if categoryDO, err = cs.data.Category().Get(ctx, int64(category.ID)); err != nil {
		return err
	}

	categoryDO.Name = category.Name
	categoryDO.ParentCategoryID = category.ParentCategoryID
	categoryDO.Level = category.Level
	categoryDO.IsTab = category.IsTab

	if err = cs.data.Category().Update(ctx, &category.CategoryDO); err != nil {
		//log.Errorf("data.Update err: %v", err)
		return err
	}
	return nil
}

func (cs *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := cs.data.Category().Delete(ctx, id); err != nil {
		log.Errorf("data.Delete err: %v", err)
		return err
	}
	return nil
}

func newCategory(srv *serviceFactory) service.CategorySrv {
	return &CategoryService{
		data:      srv.data,
		seachData: srv.dataSearch,
	}
}
