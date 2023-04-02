package client

import (
	"context"
	"errors"
	"net"
)

type Client interface {
	Connect(ctx context.Context) error
	Start(ctx context.Context) error
	Stop() error
}

type Middleware func(ctx context.Context, conn net.Conn, buf []byte) ([]byte, error)

var (
	ErrServerAddressRequired = errors.New("server address is required")
	ErrWrongResponse         = errors.New("wrong response")
	ErrClientNotConnected    = errors.New("client not connected")
)
