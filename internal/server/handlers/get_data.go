package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

func (h *handlerManager) GetDataHandler() Handler {
	return func(ctx context.Context, conn net.Conn, req []byte) error {
		h.log.Debug("<-", string(req))

		// operation_name = contents[0]
		// token = contents[1]
		// header = contents[2]
		contents := bytes.Split(req, []byte("\n"))
		if len(contents) < 2 {
			return fmt.Errorf("wrong request")
		}

		if string(contents[0]) != GetDataHandlerName {
			return fmt.Errorf("wrong handler name")
		}

		if len(contents[2]) == 0 {
			return fmt.Errorf("wrong header")
		}

		header, err := hashcash.Parse(string(contents[2]))
		if err != nil {
			return fmt.Errorf("wrong header: %w", err)
		}

		token := string(contents[1])

		h.log.Debug("RECEIVED_TOKEN", token)
		generatedToken := generateTokenFromHeader(header)
		h.log.Debug("GENERATED_TOKEN", generatedToken)

		if token != generatedToken {
			return fmt.Errorf("wrong token: generated: %s, received: %s", generatedToken, token)
		}

		t := time.Unix(header.ExpiredAt, 0)
		if !time.Now().Before(t) {
			return fmt.Errorf("header expired")
		}

		if !header.IsValid() {
			return fmt.Errorf("header is wrong value")
		}

		// Return always the quote with the same ID
		text := h.dbInstance.GetRandomDataFromDB()

		h.log.Debug("->", text)

		if _, err = fmt.Fprintf(conn, "%s\n%s\n%s\n", GetDataHandlerName, contents[1], text); err != nil {
			if !errors.Is(err, net.ErrClosed) || !errors.Is(err, io.EOF) || !errors.Is(err, io.ErrClosedPipe) {
				h.log.Error("error on reading info from client", err)
				return err
			}
		}

		h.log.Info("Sent text:", text)

		return nil
	}
}
