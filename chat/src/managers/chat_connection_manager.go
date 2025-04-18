package managers

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ChatConnManager struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func NewChatConnManager() *ChatConnManager {
	return &ChatConnManager{
		conns: make(map[string]*websocket.Conn),
	}
}

func (m *ChatConnManager) Add(ip string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conns[ip] = conn
}

func (m *ChatConnManager) Remove(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, ip)
}

func (m *ChatConnManager) Get(ip string) (*websocket.Conn, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.conns[ip]
	return conn, ok
}
