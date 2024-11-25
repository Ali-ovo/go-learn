package server

import (
	"context"
	"net/url"
)

type Server interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
