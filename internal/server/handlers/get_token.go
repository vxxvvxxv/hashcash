package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

func (h *handlerManager) GetTokenHandler() Handler {
	return func(ctx context.Context, conn net.Conn, req []byte) error {
		// handler\ntoken\nheader
		h.log.Debug("<-", string(req))

		// operation_name = contents[0]
		// token = contents[1]
		// header = contents[2]
		contents := bytes.Split(req, []byte("\n"))
		if len(contents) < 2 {
			return fmt.Errorf("wrong request")
		}

		if string(contents[0]) != GetTokenHandlerName {
			return fmt.Errorf("wrong handler name")
		}

		header, err := createHeader(h.serverAddr, h.headerDifficult, h.headerTTL)
		if err != nil {
			h.log.Error("error on creating header", err)
			return err
		}

		h.log.Debug("CREATED: ", header.String())

		token := generateTokenFromHeader(header)

		res := fmt.Sprintf("%s\n%s\n%s\n", contents[0], token, header.String())

		h.log.Debug("->", res)

		if _, err = fmt.Fprint(conn, res); err != nil {
			h.log.Error("error on sending response", err)
			return err
		}

		return nil
	}
}

func generateTokenFromHeader(header hashcash.Header) string {
	h := md5.New()
	h.Write([]byte(header.Nonce))
	mid := h.Sum(nil)
	h.Reset()

	h.Write(mid)
	h.Write([]byte(header.Subject))
	h.Write([]byte(header.Alg))
	h.Write([]byte(strconv.Itoa(int(header.ExpiredAt))))

	return hex.EncodeToString(h.Sum(nil))
}

func createHeader(subject string, difficult int, ttl time.Duration) (hashcash.Header, error) {
	expired := time.Now().Add(ttl).Unix()
	header, err := hashcash.Default(subject, difficult, expired)
	if err != nil {
		return hashcash.Header{}, err
	}

	return header, nil
}
