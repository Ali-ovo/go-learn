package srvSms

import "context"

type SmsSrv interface {
	SendSms(ctx context.Context, mobile string, tp string) error
}
