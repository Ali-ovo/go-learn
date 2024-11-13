package server

import "context"

type Server interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
}
