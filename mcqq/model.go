package mcqq

import "github.com/RomiChan/websocket"

type MinecraftBot struct {
	Websocket  *websocket.Conn
	RconClient *RCONClientConn
}

var McBots = make(map[string]*MinecraftBot)
