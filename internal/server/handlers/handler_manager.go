package handlers

import (
	"crypto/md5"
	"fmt"
	"hash"
	"strings"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/db"
	"github.com/vxxvvxxv/hashcash/internal/logger"
)

type handlerManager struct {
	dbInstance      db.DB
	log             logger.Logger
	serverAddr      string
	hashServer      hash.Hash
	headerTTL       time.Duration
	headerDifficult int
}

func NewHandlerManager(
	dbInstance db.DB,
	log logger.Logger,
	serverAddr string,
	headerTTL time.Duration,
	headerDifficult int,
) HandlerManager {
	return &handlerManager{
		dbInstance:      dbInstance,
		log:             log,
		serverAddr:      strings.ReplaceAll(serverAddr, ":", "_"),
		hashServer:      md5.New(),
		headerTTL:       headerTTL,
		headerDifficult: headerDifficult,
	}
}

func (h *handlerManager) GetHandler(name string) (Handler, error) {
	switch name {
	case GetTokenHandlerName:
		return h.GetTokenHandler(), nil
	case GetDataHandlerName:
		return h.GetDataHandler(), nil
	default:
		return nil, fmt.Errorf("wrong handler name")
	}
}
