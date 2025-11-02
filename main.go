package main

import (
	"github.com/17TheWord/zerobot-plugin-mcqq/mcqq"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/wdvxdr1123/ZeroBot/driver"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	mcqq.PluginConfig = mcqq.Config{
		WebsocketServer: mcqq.WebsocketServerConfig{
			Enable: false,
			Host:   "127.0.0.1",
			Port:   8080,
		},
		WebsocketClient: []mcqq.WebsocketClientConfig{
			{
				ServerName: "Server",
				Url:        "ws://127.0.0.1:8080",
			},
		},
		AccessToken:     "",
		CommandPriority: 2,
		ChatImage:       true,
		ServerMap: map[string]mcqq.ServerConfig{
			"Server": {
				GroupList: []mcqq.GroupConfig{
					{
						BotId:   0,
						GroupId: 0,
					},
				},
				RconMsg: false,
			},
		},
	}
	mcqq.InitPlugin()
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{0},
		Driver: []zero.Driver{
			// 正向 WS
			//driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
			// 反向 WS
			driver.NewWebSocketServer(16, "ws://127.0.0.1:9090", ""),
		},
	}, nil)
}
