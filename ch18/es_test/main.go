package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	host := "http://192.168.189.128:9200"
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	q := elastic.NewMatchQuery("address", "street")

	result, err := client.Search().Index("user").Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}

	total := result.TotalHits()
	fmt.Printf("Found 搜索数量：  %d\n", total)

	for _, value := range result.Hits.Hits {
		if jsonData, err := value.Source.MarshalJSON(); err == nil {
			fmt.Println(string(jsonData))
		}
	}
}
