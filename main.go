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
		Host: "127.0.0.1",
		Port: 8085,
		ForwardUrlList: []mcqq.ForwardServer{
			{
				ServerName: "Server",
				Url:        "ws://127.0.0.1:8080",
			},
		},
		AccessToken:     "",
		CommandPriority: 1,
		ServerMap: map[string]mcqq.Server{
			"Server": {
				GroupList: []mcqq.Group{
					{
						BotId:   0,
						GroupId: 0,
					},
				},
				RconCmd: false,
				RconMsg: false,
			},
		},
	}
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{0},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
			// 反向 WS
			driver.NewWebSocketServer(16, "ws://127.0.0.1:8081", ""),
		},
	}, nil)
}
