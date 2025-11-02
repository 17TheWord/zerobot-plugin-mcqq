package mcqq

import (
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var PluginConfig = Config{}
var McBots = make(map[string]*websocket.Conn)

func InitPlugin() {
	groupIdSet = make(map[int64]struct{})
	for _, server := range PluginConfig.ServerMap {
		for _, group := range server.GroupList {
			groupIdSet[group.GroupId] = struct{}{}
		}
	}

	if PluginConfig.CommandPriority < 2 {
		PluginConfig.CommandPriority = 2
		log.Warning("命令优先级过低，已设置为 2")
	} else if PluginConfig.CommandPriority > 98 {
		PluginConfig.CommandPriority = 98
		log.Warning("命令优先级过高，已设置为 98")
	} else {
		log.Infof("命令优先级设置为 %d", PluginConfig.CommandPriority)
	}

	if PluginConfig.AccessToken == "" {
		log.Info("未设置访问令牌 AccessToken")
	} else {
		log.Info("已设置访问令牌 AccessToken")
	}

	if PluginConfig.WebsocketServer.Enable {
		log.Info("正在启动 Websocket 服务器...")
		go startWebsocketServer()
		log.Info("Websocket 服务器已启动。")
	} else {
		log.Info("Websocket 服务器未启用。")
	}

	if len(PluginConfig.WebsocketClient) > 0 {
		log.Info("正在启动 Websocket 客户端...")
		startWebsocketClient()
		log.Info("Websocket 客户端已启动。")
	} else {
		log.Info("未配置任何 Websocket 客户端。")
	}

	zero.OnMessage(GroupRule).SetBlock(false).SetPriority(PluginConfig.CommandPriority + 1).Handle(handleQQMessage)
}
