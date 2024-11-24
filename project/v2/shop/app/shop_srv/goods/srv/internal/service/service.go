package service

type ServiceFactory interface {
	Goods() GoodsSrv
}
