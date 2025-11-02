package mcqq

import (
	"net/http"
	"time"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var dialer = websocket.Dialer{
	HandshakeTimeout: 5 * time.Second,
}

func forwardWebsocket(serverName string, url string) {
	requestHeader := make(http.Header)
	requestHeader.Set("x-client-origin", "zerobot")
	requestHeader.Set("x-self-name", serverName)
	requestHeader.Set("Authorization", "Bearer "+PluginConfig.AccessToken)

	conn, _, err := dialer.Dial(url, requestHeader)
	if err != nil {
		log.Errorf("Failed to connect to websocket [%s]: %v", serverName, err)
		return
	}

	McBots[serverName] = conn
	log.Infof("Connected to websocket [%s]", serverName)

	defer cleanupWebSocketConnection(conn, serverName)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warningf("Read websocket message from [%s] failed: %v", serverName, err)
			break
		}
		handleMinecraftMessage(message)
	}
}

func startWebsocketClient() {
	for _, server := range PluginConfig.WebsocketClient {
		go forwardWebsocket(server.ServerName, server.Url)
	}
}
