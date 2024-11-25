package do

import (
	"database/sql/driver"
	"encoding/json"
	"shop/pkg/gorm"
)

type GoodsDetail struct {
	Goods int64
	Num   int32
}

type GoodsDetailList []GoodsDetail

func (gdl GoodsDetailList) Len() int {
	return len(gdl)
}
func (gdl GoodsDetailList) Less(i, j int) bool {
	return gdl[i].Goods < gdl[j].Goods
}
func (gdl GoodsDetailList) Swap(i, j int) {
	gdl[i], gdl[j] = gdl[j], gdl[i]
}

func (gdl GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(gdl)
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (gdl *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &gdl)
}

type InventoryDO struct {
	gorm.BaseModel
	Goods   int64 `gorm:"type:int;index"` // 商品id
	Stocks  int32 `gorm:"type:int"`       // 库存数量
	Version int32 `gorm:"type:int"`       // 乐观锁
}

func (id *InventoryDO) TableName() string {
	return "inventory"
}

type StockSellDetailDO struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"` // 订单编号
	Status  int32           `gorm:"type:varchar(200)"`                            //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`                            // 订单中的 商品id 和 商品数量
}

func (ssd *StockSellDetailDO) TableName() string {
	return "stockselldetail"
}
