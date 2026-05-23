package mcqq

import (
	"testing"

	"github.com/RomiChan/websocket"
)

func TestConnectionStoreSetGetExists(t *testing.T) {
	store := newConnectionStore()
	conn := &websocket.Conn{}

	if store.exists("Server") {
		t.Fatal("connection should not exist before set")
	}

	store.set("Server", conn)

	got, ok := store.get("Server")
	if !ok {
		t.Fatal("expected stored connection")
	}
	if got != conn {
		t.Fatal("stored connection mismatch")
	}
	if !store.exists("Server") {
		t.Fatal("connection should exist after set")
	}
}

func TestConnectionStoreDeleteIfSame(t *testing.T) {
	store := newConnectionStore()
	oldConn := &websocket.Conn{}
	newConn := &websocket.Conn{}

	store.set("Server", oldConn)
	store.set("Server", newConn)

	if store.deleteIfSame("Server", oldConn) {
		t.Fatal("old connection must not delete newer connection")
	}

	got, ok := store.get("Server")
	if !ok {
		t.Fatal("new connection should still exist")
	}
	if got != newConn {
		t.Fatal("new connection was replaced unexpectedly")
	}

	if !store.deleteIfSame("Server", newConn) {
		t.Fatal("expected new connection to be deleted")
	}
	if store.exists("Server") {
		t.Fatal("connection should not exist after delete")
	}
}
