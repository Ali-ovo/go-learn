package es

import (
	"context"
	"encoding/json"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"
	"shop/pkg/es"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type goods struct {
	esClient *elasticsearch.Client
}

func (g *goods) Search(ctx context.Context, req *data_search.GoodsFilterRequest) (*do.GoodsSearchDOList, error) {
	var ret do.GoodsSearchDOList

	boolQuery, query := (*do.GoodsSearchDO)(nil).GetSearchBool()
	if req.KeyWords != "" {
		multiMatchQuery := map[string]any{
			"multi_match": map[string]any{
				"query":  req.KeyWords,
				"fields": []string{"name", "goods_brief"},
			},
		}
		boolQuery["must"] = append(boolQuery["must"].([]any), multiMatchQuery)
	}
	if req.IsHot {
		termQuery := map[string]any{
			"term": map[string]any{
				"is_hot": req.IsHot,
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]any), termQuery)
	}
	if req.IsNew {
		termQuery := map[string]any{
			"term": map[string]any{
				"is_new": req.IsNew,
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]any), termQuery)
	}
	if req.PriceMin > 0 {
		rangeQuery := map[string]any{
			"range": map[string]any{
				"shop_price": map[string]any{
					"gte": req.PriceMin,
				},
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]any), rangeQuery)
	}
	if req.PriceMax > 0 {
		rangeQuery := map[string]any{
			"range": map[string]any{
				"shop_price": map[string]any{
					"lte": req.PriceMax,
				},
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]any), rangeQuery)
	}
	if req.Brand > 0 {
		termQuery := map[string]any{
			"term": map[string]any{
				"brand_id": req.Brand,
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]any), termQuery)
	}

	// 通过 category 去查询商品
	if req.TopCategory > 0 {
		termsQuery := map[string]any{
			"terms": map[string]any{
				"category_id": req.CategoryIDs,
			},
		}
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), termsQuery)
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	query["from"] = (req.Pages - 1) * req.PagePerNums
	query["size"] = req.PagePerNums

	res, err := esapi.SearchRequest{
		Index: []string{(*do.GoodsSearchDO)(nil).GetIndexName()},
		Body:  esutil.NewJSONReader(query),
	}.Do(context.Background(), g.esClient)
	if err != nil {
		return nil, errors.WithCode(code2.ErrEsDatabase, err.Error())
	}
	defer res.Body.Close()

	// 处理查询结果
	var result es.SearchResult
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.WithCode(code.ErrDecodingFailed, err.Error())
	}
	ret.TotalCount = result.Hits.TotalHits.Value
	// 获取商品id
	for _, value := range result.Hits.Hits {
		goodsSearch := do.GoodsSearchDO{}
		err = json.Unmarshal(value.Source, &goodsSearch)
		if err != nil {
			return nil, errors.WithCode(code2.ErrEsUnmarshal, err.Error())
		}
		ret.Items = append(ret.Items, &goodsSearch)
	}
	return &ret, nil
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsSearchDO) error {
	_, err := esapi.IndexRequest{
		Index:      goods.GetIndexName(),
		DocumentID: strconv.Itoa(int(goods.ID)),
		Body:       esutil.NewJSONReader(goods),
	}.Do(ctx, g.esClient)
	if err != nil {
		return errors.WithCode(code2.ErrEsDatabase, err.Error())
	}
	return nil
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsSearchDO) error {
	_, err := esapi.UpdateRequest{
		Index:      goods.GetIndexName(),
		DocumentID: strconv.Itoa(int(goods.ID)),
		Body:       esutil.NewJSONReader(goods),
	}.Do(ctx, g.esClient)
	if err != nil {
		return errors.WithCode(code2.ErrEsDatabase, err.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	var goodsSearch = (*do.GoodsSearchDO)(nil)

	_, err := esapi.DeleteRequest{
		Index:      goodsSearch.GetIndexName(),
		DocumentID: strconv.Itoa(int(ID)),
	}.Do(ctx, g.esClient)
	if err != nil {
		return errors.WithCode(code2.ErrEsDatabase, err.Error())
	}
	return nil
}

var _ data_search.GoodsStore = (*goods)(nil)
