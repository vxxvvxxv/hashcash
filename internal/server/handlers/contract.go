package handlers

import (
	"context"
	"net"
)

type HandlerManager interface {
	GetHandler(name string) (Handler, error)
}

type Handler func(ctx context.Context, conn net.Conn, req []byte) error

const (
	GetTokenHandlerName = "get_token"
	GetDataHandlerName  = "get_data"
)
