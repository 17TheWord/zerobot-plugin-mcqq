package mcqq

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var upGrader = websocket.Upgrader{}

func handleWebsocket(writer http.ResponseWriter, request *http.Request) {
	// 获取请求头中的 x-self-name 原始名称
	oriSelfName := request.Header.Get("x-self-name")
	// 如果请求头中没有 x-self-name 或为空，则拒绝连接
	if oriSelfName == "" {
		log.Warningf("Missing X-Self-Name Header from %s", request.RemoteAddr)
		_, _ = writer.Write([]byte("Missing X-Self-Name Header"))
		return
	}

	clientOrigin := request.Header.Get("x-client-origin")

	if clientOrigin == "zerobot" {
		log.Warningf("X-Client-Origin Header cannot be zerobot from %s", request.RemoteAddr)
		_, _ = writer.Write([]byte("X-Client-Origin Header cannot be zerobot"))
		return
	}

	// 解码
	selfName, err := url.QueryUnescape(oriSelfName)
	if err != nil {
		log.Warningf("X-Self-Name Header is not valid from %s with [%s]", request.RemoteAddr, selfName)
		_, _ = writer.Write([]byte("X-Self-Name Header is not valid"))
		return
	}

	if PluginConfig.AccessToken != "" {
		accessToken := request.Header.Get("Authorization")
		if accessToken == "" {
			log.Warningf("Missing Authorization Header from %s with [%s]", request.RemoteAddr, selfName)
			_, _ = writer.Write([]byte("Missing Authorization Header"))
			return
		}

		if accessToken != "Bearer "+PluginConfig.AccessToken {
			log.Warningf("Authorization Header is not valid from %s with [%s]", request.RemoteAddr, selfName)
			_, _ = writer.Write([]byte("Authorization Header is not valid"))
			return
		}
	}

	if _, exists := McBots[selfName]; exists {
		log.Warningf("X-Self-Name Header is already in use from %s with [%s]", request.RemoteAddr, selfName)
		_, _ = writer.Write([]byte("X-Self-Name Header is already in use"))
		return
	}

	conn, err := upGrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Warningf("Websocket from %s with [%s] connection failed: %v", request.RemoteAddr, selfName, err)
		return
	}

	McBots[selfName] = conn
	log.Infof("Websocket from [%s] connected", selfName)

	defer cleanupWebSocketConnection(conn, selfName)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warningf("Read websocket message from [%s] failed: %v", selfName, err)
			break
		}
		handleMinecraftMessage(message)
	}
}

func startWebsocketServer() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/minecraft", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/ws", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/ws/", handleWebsocket)

	host := fmt.Sprintf("%s:%d", PluginConfig.WebsocketServer.Host, PluginConfig.WebsocketServer.Port)
	log.Infof("Starting server at %s...", host)
	if err := http.ListenAndServe(host, httpHandler); err != nil {
		log.Fatalf("Failed to start server on %s: %v", host, err)
	}
}
