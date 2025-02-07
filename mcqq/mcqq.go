package mcqq

import (
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var PluginConfig = Config{}

func init() {
	if PluginConfig.CommandPriority < 1 {
		PluginConfig.CommandPriority = 1
		log.Warn("CommandPriority is too low, set to 1")
	} else if PluginConfig.CommandPriority > 98 {
		PluginConfig.CommandPriority = 98
		log.Warn("CommandPriority is too high, set to 98")
	}

	if PluginConfig.AccessToken == "" {
		log.Infoln("AccessToken is empty!!!")
	}

	StartWebsocket()
	zero.OnMessage(GroupRule).SetBlock(false).SetPriority(PluginConfig.CommandPriority + 1).Handle(handleQQMessage)
}
