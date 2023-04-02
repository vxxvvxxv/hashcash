package requests

import (
	"context"
	"net"

	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

type Request func(ctx context.Context, conn net.Conn) (interface{}, error)

type ResponseGetToken struct {
	Token  string
	Header hashcash.Header
}

type ResponseGetData struct {
	Text string
}
