package requests

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/vxxvvxxv/hashcash/internal/logger"
	"github.com/vxxvvxxv/hashcash/internal/server/handlers"
	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

func RequestGetToken(
	log logger.Logger,
) Request {
	return func(ctx context.Context, conn net.Conn) (interface{}, error) {
		req := fmt.Sprintf("%s\n", handlers.GetTokenHandlerName)

		log.Debug("->", req)

		if _, err := fmt.Fprint(conn, req); err != nil {
			log.Error("error on sending info to server", err)
			return nil, err
		}

		// Read response
		res := make([]byte, 512)
		_, err := conn.Read(res)
		if err != nil {
			log.Error("error on read info from server", err)
			return nil, err
		}

		// _, err := conn.Read(res)
		// if err != nil {
		// 	log.Error("error on read info from server", err)
		// 	return nil, err
		// }

		log.Debug("<-", string(res))

		contents := bytes.Split(res, []byte("\n"))
		if len(contents) < 2 {
			return nil, fmt.Errorf("token: wrong response, need at least 2 lines")
		}

		if string(contents[0]) != handlers.GetTokenHandlerName {
			return nil, fmt.Errorf("wrong handler name in response")
		}
		if len(contents[1]) == 0 {
			return nil, fmt.Errorf("empty token")
		}
		if len(contents[2]) == 0 {
			return nil, fmt.Errorf("empty header")
		}

		log.Debug("CONTENTS:")
		for i, v := range contents {
			log.Debug(i, " ", string(v))
		}

		token := string(contents[1])
		headerStr := string(contents[2])

		log.Debug("SERVER TOKEN: ", token)
		log.Debug("SERVER HEADER: ", headerStr)

		header, err := hashcash.Parse(headerStr)
		if err != nil {
			return nil, err
		}

		log.Debug("PARSED HEADER: ", header.String())

		return ResponseGetToken{
			Token:  string(contents[1]),
			Header: header,
		}, nil
	}

}
