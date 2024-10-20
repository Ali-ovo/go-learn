package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

type Account struct {
	AccountNumber int32  `json:"account_number"`
	FirstName     string `json:"firstname"`
}

const goodsMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
			"properties":{
				"name":{
					"type":"text",
					"analyzer":"ik_max_word"
				},
				"id":{
					"type":"integer"
				}
			}
		}
}`

func main() {
	host := "http://192.168.189.128:9200"
	logger := log.New(os.Stdout, "APP", log.LstdFlags)
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
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
		account := Account{}
		_ = json.Unmarshal(value.Source, &account)

		// fmt.Printf("account_number: %d, firstname: %s\n", account.AccountNumber, account.FirstName)
	}

	// account := Account{AccountNumber: 123456, FirstName: "Ali233"}
	// put1, err := client.Index().
	// 	Index("myuser").
	// 	BodyJson(account).
	// 	Do(context.Background())
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	createIndex, err := client.CreateIndex("mygoods").BodyString(goodsMapping).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}
