package mapstructure

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Product struct {
	Name  string  `mapstructure:"name"`
	Price float32 `mapstructure:"price"`
}

func TestDecode(t *testing.T) {
	jsonStr := `{"name": "Product A", "price": "12.34"}`
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		t.Fatal(err)
	}

	var product Product
	err = Decode(m, &product)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(product.Name, product.Price)
}
