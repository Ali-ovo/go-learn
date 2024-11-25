package data

import "gorm.io/gorm"

type DataFactory interface {
	Inventory() InventoryStore
	Begin() *gorm.DB
}
