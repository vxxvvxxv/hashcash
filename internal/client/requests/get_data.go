package requests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/vxxvvxxv/hashcash/internal/logger"
	"github.com/vxxvvxxv/hashcash/internal/server/handlers"
	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

func RequestGetData(log logger.Logger, token string, computedHeader hashcash.Header) Request {
	return func(ctx context.Context, conn net.Conn) (interface{}, error) {

		req := fmt.Sprintf("%s\n%s\n%s\n", handlers.GetDataHandlerName, token, computedHeader.String())

		log.Debug("->", req)

		if _, err := fmt.Fprint(conn, req); err != nil {
			if !errors.Is(err, net.ErrClosed) || !errors.Is(err, io.EOF) || !errors.Is(err, io.ErrClosedPipe) {
				log.Error("error on sending info to server", err)
				return nil, err
			}
		}

		// Response can be very big
		res := make([]byte, 1024)
		_, err := conn.Read(res)
		if err != nil {
			log.Error("error on read info from server", err)
			return nil, err
		}

		log.Debug("<-", string(res))

		contents := bytes.Split(res, []byte("\n"))
		if len(contents) < 2 {
			return nil, fmt.Errorf("data: wrong response, need at least 2 lines")
		}

		if string(contents[0]) != handlers.GetDataHandlerName {
			return nil, fmt.Errorf("wrong handler name in response")
		}

		return ResponseGetData{
			Text: string(contents[2]),
		}, nil
	}

}
