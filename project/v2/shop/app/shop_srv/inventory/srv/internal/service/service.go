package service

type ServiceFactory interface {
	Inventory() InventorySrv
}
