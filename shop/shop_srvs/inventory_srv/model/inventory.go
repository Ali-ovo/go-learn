package model

// type Stock struct {
// 	BaseModel
// 	Name    string
// 	Address string
// }

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式 乐观锁
}

// type InventoryHistory struct {
// 	user   int32
// 	goods  int32
// 	nums   int32
// 	order  int32
// 	status int32 // 1 表示与扣减， 2 表示支付成功
// }
