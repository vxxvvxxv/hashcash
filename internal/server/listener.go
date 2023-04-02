package server

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
	keepAlive bool
	ttl       time.Duration
}

func wrapTCPKeepAliveListener(l *net.TCPListener, ttl time.Duration) *tcpKeepAliveListener {
	return &tcpKeepAliveListener{
		TCPListener: l,
		keepAlive:   true,
		ttl:         ttl,
	}
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	if err = tc.SetKeepAlive(ln.keepAlive); err != nil {
		return tc, err
	}
	if err = tc.SetKeepAlivePeriod(ln.ttl); err != nil {
		return tc, err
	}
	return tc, nil
}

func (ln tcpKeepAliveListener) Close() error {
	return ln.TCPListener.Close()
}
