package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/db"
	"github.com/vxxvvxxv/hashcash/internal/logger"
	"github.com/vxxvvxxv/hashcash/internal/server/handlers"
)

type serverService struct {
	serverAddr string
	db         db.DB
	log        logger.Logger

	listener net.Listener

	ctx       context.Context
	ctxCancel context.CancelFunc

	wg             sync.WaitGroup
	conns          *connManager
	connTTL        time.Duration
	handlerManager handlers.HandlerManager
}

func NewServer(
	serverAddr string,
	db db.DB,
	log logger.Logger,
	connTTL time.Duration,
	handlerManager handlers.HandlerManager,
) (Server, error) {
	if len(serverAddr) == 0 {
		return nil, errors.New("addr is wrong")
	}

	s := &serverService{
		serverAddr:     serverAddr,
		db:             db,
		log:            log,
		connTTL:        connTTL,
		conns:          newConnManager(),
		handlerManager: handlerManager,
	}

	return s, nil
}

func (s *serverService) ListenAndServe(ctx context.Context) error {
	s.ctx, s.ctxCancel = context.WithCancel(ctx)
	defer s.ctxCancel()

	s.wg.Add(1)
	defer s.wg.Done()

	var errListen error

	s.listener, errListen = net.Listen("tcp", s.serverAddr)
	if errListen != nil {
		s.log.Error("can't create serverService: ", errListen)
		return errListen
	}
	defer func() {
		if errClose := s.listener.Close(); errClose != nil && !errors.Is(errClose, net.ErrClosed) {
			s.log.Error("can't close serverService: ", errClose)
		}
	}()

	s.log.Info("Listening on: " + s.serverAddr)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			// TODO: Check this

			conn, errAccept := s.listener.Accept()
			if errAccept != nil {
				select {
				case <-s.ctx.Done():
					return nil
				default:
					s.log.Error("can't accept connection: ", errAccept)
					return errAccept
				}
			}

			s.log.Debug("accepted connection")

			// Add connection to manager
			s.conns.add(conn)

			s.wg.Add(1)
			// Will be closed in handleConnection
			go s.handleConnection(s.ctx, conn)
		}
	}
}

func (s *serverService) Shutdown() error {
	s.log.Debug("Shutdown serverService")
	s.ctxCancel()

	s.log.Debug("Closing all connections")
	s.conns.closeAll()
	s.log.Debug("Closed all connections")

	s.log.Debug("Stopping listener")
	if s.listener != nil {
		if errClose := s.listener.Close(); errClose != nil && !errors.Is(errClose, net.ErrClosed) {
			s.log.Error("can't close serverService: ", errClose)
		}
	}
	s.log.Debug("Stopped listener")

	s.log.Debug("Waiting for compete all connections")
	s.wg.Wait()
	s.log.Debug("Completed all connections")

	return nil
}

func (s *serverService) handleConnection(ctx context.Context, conn net.Conn) {
	s.log.Debug("Start handleConnection")
	defer s.log.Debug("End handleConnection")

	defer s.wg.Done()
	defer s.conns.remove(conn)

	defer func() {
		if errClose := conn.Close(); errClose != nil && !errors.Is(errClose, net.ErrClosed) {
			s.log.Error("can't close connection: ", errClose)
		}
	}()

	if err := conn.SetDeadline(time.Now().Add(s.connTTL)); err != nil {
		s.log.Error("conn set deadline: %v", err)
		return
	}

	// r := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := s.handleOperation(ctx, conn); err != nil {
				if !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
					s.log.Error("can't handle operation: ", err)
				}
				return
			}
		}
	}
}

func (s *serverService) handleOperation(ctx context.Context, conn net.Conn) error {
	// req, err = r.ReadBytes('\n')
	req := make([]byte, 512)
	_, err := conn.Read(req)
	if err != nil {
		if !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
			s.log.Error(fmt.Sprintf("can't read from connection: %v", err))
		}
		return err
	}

	// handler = contents[0]
	contents := bytes.Split(req, []byte("\n"))
	if len(contents) < 1 {
		return fmt.Errorf("wrong request: %v", string(req))
	}

	handler, err := s.handlerManager.GetHandler(string(contents[0]))
	if err != nil {
		return err
	}

	if err = handler(ctx, conn, req); err != nil {
		return err
	}

	return nil
}
