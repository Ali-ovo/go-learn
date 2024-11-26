package service

type ServiceFactory interface {
	Orders() OrderSrv
	ShopCart() ShopCartSrv
}
