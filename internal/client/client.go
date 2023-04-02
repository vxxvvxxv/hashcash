package client

import (
	"context"
	"net"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/client/requests"
	"github.com/vxxvvxxv/hashcash/internal/logger"
	"github.com/vxxvvxxv/hashcash/pkg/hashcash"
)

type clientService struct {
	serverAddr   string
	conn         net.Conn
	isConnClosed bool
	connTimeout  time.Duration
	ctx          context.Context
	ctxCancel    context.CancelFunc
	log          logger.Logger

	// Header
	poolOpts []hashcash.PoolOption
}

func NewClient(
	serverAddr string,
	connTimeout time.Duration,
	log logger.Logger,
	headerMaxIterations int,
	headerTTL time.Duration,
) (Client, error) {
	if len(serverAddr) == 0 {
		return nil, ErrServerAddressRequired
	}

	poolOpts := make([]hashcash.PoolOption, 0)
	if headerTTL >= 0 {
		poolOpts = append(poolOpts, hashcash.WithPoolDuration(headerTTL))
	}
	if headerMaxIterations >= 0 {
		poolOpts = append(poolOpts, hashcash.WithPoolMaxIterations(headerMaxIterations))
	}

	return &clientService{
		serverAddr:  serverAddr,
		connTimeout: connTimeout,
		log:         log,
		poolOpts:    poolOpts,
	}, nil
}

func (c *clientService) Connect(ctx context.Context) (err error) {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	defer c.ctxCancel()

	if c.conn == nil || c.isConnClosed {
		c.log.Debug("Connecting to server with address: " + c.serverAddr + " with timeout: " + c.connTimeout.String())

		// Create connection
		var errDial error
		c.conn, errDial = net.DialTimeout("tcp", c.serverAddr, c.connTimeout)
		if errDial != nil {
			c.log.Error("can't connect to server with credentials: %s, err: %w", c.serverAddr, errDial)
			return errDial
		}

		c.log.Info("Connected to server")
	}
	return nil
}

func (c *clientService) Start(ctx context.Context) (err error) {
	c.ctx, c.ctxCancel = context.WithCancel(ctx)
	defer c.ctxCancel()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			if c.conn == nil || c.isConnClosed {
				err = c.Connect(ctx)
				if err != nil {
					return err
				}
			}
			if c.conn == nil {
				return ErrClientNotConnected
			}

			// Get data from server
			text, err := c.getDataFromServer(ctx, c.conn)
			if err != nil {
				return err
			}

			// Print data
			c.log.Info("SUCCESS! Data from server: \n", string(text))

			return nil
		}
	}
}

func (c *clientService) Stop() error {
	if c.ctxCancel == nil {
		return ErrClientNotConnected
	}
	c.ctxCancel()

	if c.conn != nil {
		c.isConnClosed = true
		return c.conn.Close()
	}

	return nil
}

func (c *clientService) getDataFromServer(ctx context.Context, conn net.Conn) ([]byte, error) {
	resTokenRaw, errToken := requests.RequestGetToken(c.log)(ctx, conn)
	if errToken != nil {
		return nil, errToken
	}
	resToken, ok := resTokenRaw.(requests.ResponseGetToken)
	if !ok {
		return nil, ErrWrongResponse
	}

	// Compute header
	c.log.Debug("RECEIVED: ", resToken.Header.String())

	info, errCompute := hashcash.ComputeWithPool(ctx, resToken.Header, c.poolOpts...)
	if errCompute != nil {
		return nil, errCompute
	}

	c.log.Debug("COMPUTED: ", info.Header.String())

	c.log.Info("Time to compute", info.Time.String(), " work number", info.WorkerNum, "hash", info.Header.Hash())

	// Get data from server
	resDataRaw, errData := requests.RequestGetData(c.log, resToken.Token, info.Header)(ctx, conn)
	if errData != nil {
		return nil, errData
	}
	resData, ok := resDataRaw.(requests.ResponseGetData)
	if !ok {
		return nil, ErrWrongResponse
	}

	return []byte(resData.Text), nil

}
