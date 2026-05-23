package mcqq

import (
	"sync"

	"github.com/RomiChan/websocket"
)

type connectionStore struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func newConnectionStore() *connectionStore {
	return &connectionStore{conns: make(map[string]*websocket.Conn)}
}

func (s *connectionStore) get(serverName string) (*websocket.Conn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.conns[serverName]
	return conn, ok
}

func (s *connectionStore) set(serverName string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[serverName] = conn
}

func (s *connectionStore) exists(serverName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.conns[serverName]
	return ok
}

func (s *connectionStore) deleteIfSame(serverName string, conn *websocket.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, ok := s.conns[serverName]; ok && current == conn {
		delete(s.conns, serverName)
		return true
	}
	return false
}
