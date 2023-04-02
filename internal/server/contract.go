package server

import (
	"context"
)

type Server interface {
	ListenAndServe(ctx context.Context) error
	Shutdown() error
}
