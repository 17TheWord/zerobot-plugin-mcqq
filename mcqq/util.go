package mcqq

import (
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

func getTargetServer(serverName string) *MinecraftBot {
	if server, exist := McBots[serverName]; exist {
		return server
	}
	return nil
}

func getTargetServerName(groupId int64) string {
	for serverName, server := range PluginConfig.ServerMap {
		for _, group := range server.GroupList {
			if group.GroupId == groupId {
				return serverName
			}
		}
	}
	return ""
}

func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func cleanupWebSocketConnection(conn *websocket.Conn, serverName string) {
	err := conn.Close()
	if err != nil {
		log.Infof("Close websocket connection from [%s] failed: %v", serverName, err)
		return
	}
	delete(McBots, serverName)
	log.Infoln("Disconnected from websocket [%s]\n", serverName)
}
