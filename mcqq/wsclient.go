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

func forwardWebsocket(config WebsocketClientConfig) {
	serverName := config.ServerName
	reconnectInterval := time.Duration(config.ReconnectInterval) * time.Second
	if reconnectInterval <= 0 {
		reconnectInterval = 5 * time.Second
	}

	for reconnectTimes := 0; ; reconnectTimes++ {
		if config.ReconnectMaxTimes > 0 && reconnectTimes > config.ReconnectMaxTimes {
			log.Warningf("Reconnect websocket [%s] stopped after %d attempts", serverName, config.ReconnectMaxTimes)
			return
		}

		if reconnectTimes > 0 {
			log.Infof("Reconnect websocket [%s] after %s, attempt %d", serverName, reconnectInterval, reconnectTimes)
			time.Sleep(reconnectInterval)
		}

		readWebsocket(config)
	}
}

func readWebsocket(config WebsocketClientConfig) {
	requestHeader := make(http.Header)
	requestHeader.Set("x-client-origin", "zerobot")
	requestHeader.Set("x-self-name", config.ServerName)
	requestHeader.Set("Authorization", "Bearer "+PluginConfig.AccessToken)

	conn, _, err := dialer.Dial(config.Url, requestHeader)
	if err != nil {
		log.Errorf("Failed to connect to websocket [%s]: %v", config.ServerName, err)
		return
	}

	mcConnections.set(config.ServerName, conn)
	log.Infof("Connected to websocket [%s]", config.ServerName)

	defer cleanupWebSocketConnection(conn, config.ServerName)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warningf("Read websocket message from [%s] failed: %v", config.ServerName, err)
			break
		}
		handleMinecraftMessage(message)
	}
}

func startWebsocketClient() {
	for _, server := range PluginConfig.WebsocketClient {
		go forwardWebsocket(server)
	}
}
