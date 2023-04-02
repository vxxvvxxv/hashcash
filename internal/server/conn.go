package server

import (
	"net"
	"sync"
)

type connManager struct {
	connections map[net.Conn]struct{}
	mu          sync.Mutex
}

func newConnManager() *connManager {
	return &connManager{
		connections: make(map[net.Conn]struct{}),
	}
}

func (cm *connManager) add(conn net.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[conn] = struct{}{}
}

func (cm *connManager) remove(conn net.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, conn)
}

func (cm *connManager) closeAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for conn := range cm.connections {
		_ = conn.Close()
	}
}
