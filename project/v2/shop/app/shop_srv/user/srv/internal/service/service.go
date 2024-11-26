package service

type ServiceFactory interface {
	User() UserSrv
}
