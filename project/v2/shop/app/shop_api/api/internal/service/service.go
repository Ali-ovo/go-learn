package service

import (
	"shop/app/shop_api/api/internal/data/v1"
	srvGoods "shop/app/shop_api/api/internal/service/goods"
	serviceGoods "shop/app/shop_api/api/internal/service/goods/v1"
	srvSms "shop/app/shop_api/api/internal/service/sms"
	serviceSms "shop/app/shop_api/api/internal/service/sms/v1"
	srvUser "shop/app/shop_api/api/internal/service/user"
	serviceUser "shop/app/shop_api/api/internal/service/user/v1"
	"shop/pkg/options"
)

type ServiceFactory interface {
	Goods() srvGoods.GoodsSrv
	User() srvUser.UserSrv
	Sms() srvSms.SmsSrv
}

type service struct {
	data    data.DataFactory
	smsOpts *options.SmsOptions
	jwtOpts *options.JwtOptions
}

func (s *service) Sms() srvSms.SmsSrv {
	return serviceSms.NewSmsService(s.smsOpts)
}

func (s *service) Goods() srvGoods.GoodsSrv {
	return serviceGoods.NewGoodsService(s.data)
}

func (s *service) User() srvUser.UserSrv {
	return serviceUser.NewUserService(s.data, s.jwtOpts)
}

func NewService(data data.DataFactory, smsOpts *options.SmsOptions, jwtOpts *options.JwtOptions) *service {
	return &service{data: data,
		smsOpts: smsOpts,
		jwtOpts: jwtOpts,
	}
}

var _ ServiceFactory = &service{}
