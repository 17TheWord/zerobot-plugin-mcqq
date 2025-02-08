package mcqq

import (
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var PluginConfig = Config{}

func init() {
	if PluginConfig.CommandPriority < 1 {
		PluginConfig.CommandPriority = 1
		log.Warning("CommandPriority is too low, set to 1")
	} else if PluginConfig.CommandPriority > 98 {
		PluginConfig.CommandPriority = 98
		log.Warning("CommandPriority is too high, set to 98")
	} else {
		log.Infof("CommandPriority is %d, the MessagePriority is %d", PluginConfig.CommandPriority, PluginConfig.CommandPriority+1)
	}

	if PluginConfig.AccessToken == "" {
		log.Warning("AccessToken is empty!!!")
	}

	go startWebsocketServer()
	startWebsocketClient()

	zero.OnMessage(GroupRule).SetBlock(false).SetPriority(PluginConfig.CommandPriority + 1).Handle(handleQQMessage)
}
