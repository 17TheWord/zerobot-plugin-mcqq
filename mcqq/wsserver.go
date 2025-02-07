package mcqq

import (
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var upGrader = websocket.Upgrader{}

func handleWebsocket(writer http.ResponseWriter, request *http.Request) {
	// 获取请求头中的 x-self-name 原始名称
	oriSelfName := request.Header.Get("x-self-name")
	// 如果请求头中没有 x-self-name 或为空，则拒绝连接
	if oriSelfName == "" {
		log.Warningln("Missing X-Self-Name Header")
		_, _ = writer.Write([]byte("Missing X-Self-Name Header"))
		return
	}

	clientOrigin := request.Header.Get("x-client-origin")
	if clientOrigin == "zerobot" {
		log.Warningln("X-Client-Origin Header cannot be zerobot")
		_, _ = writer.Write([]byte("X-Client-Origin Header cannot be zerobot"))
		return
	}

	// 解码
	selfName, err := url.QueryUnescape(oriSelfName)
	if err != nil {
		log.Warningln("X-Self-Name Header is not valid")
		_, _ = writer.Write([]byte("X-Self-Name Header is not valid"))
		return
	}

	if PluginConfig.AccessToken != "" {
		accessToken := request.Header.Get("Authorization")
		if accessToken == "" {
			log.Warningln("Missing Authorization Header")
			_, _ = writer.Write([]byte("Missing Authorization Header"))
			return
		}

		if accessToken != "Bearer "+PluginConfig.AccessToken {
			log.Warningln("Authorization Header is not valid")
			_, _ = writer.Write([]byte("Authorization Header is not valid"))
			return
		}
	}

	if _, exists := McBots[selfName]; exists {
		log.Warningln("X-Self-Name Header is already in use")
		_, _ = writer.Write([]byte("X-Self-Name Header is already in use"))
		return
	}

	connect, err := upGrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Warningln("Websocket from ["+selfName+"] connection failed:", err)
		return
	}

	McBots[selfName] = &MinecraftBot{
		Websocket:  connect,
		RconClient: nil,
	}

	log.Infoln("Websocket from [" + selfName + "] connected")

	defer func() {
		err := connect.Close()
		if err != nil {
			log.Error("Close websocket connection from ["+selfName+"] failed:", err)
			return
		}
		delete(McBots, selfName)
	}()

	for {
		_, message, err := connect.ReadMessage()
		if err != nil {
			log.Warningln("Read websocket message from ["+selfName+"] failed:", err)
			break
		}
		handleMinecraftMessage(message)
	}
}

func StartWebsocket() {
	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/minecraft", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/ws", handleWebsocket)
	httpHandler.HandleFunc("/minecraft/ws/", handleWebsocket)
	go func() {
		log.Println("WebSocket server is running and integrated with ZeroBot HTTP service.")
		log.Fatal(http.ListenAndServe("localhost:8085", httpHandler))
	}()
}
