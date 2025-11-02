package mcqq

import (
	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

func getTargetServerWebsocketList(serverNameList []string) []*websocket.Conn {
	targetServerWebsocketList := make([]*websocket.Conn, 0)
	for _, serverName := range serverNameList {
		if websocketConn, exist := McBots[serverName]; exist {
			targetServerWebsocketList = append(targetServerWebsocketList, websocketConn)
		}
	}
	return targetServerWebsocketList
}

func getTargetServerNameList(groupId int64) []string {
	var groupList []string
	for serverName, server := range PluginConfig.ServerMap {
		for _, group := range server.GroupList {
			if group.GroupId == groupId {
				groupList = append(groupList, serverName)
			}
		}
	}
	return groupList
}

func cleanupWebSocketConnection(conn *websocket.Conn, serverName string) {
	err := conn.Close()
	if err != nil {
		log.Infof("Close websocket connection from [%s] failed: %v", serverName, err)
		return
	}
	delete(McBots, serverName)
	log.Infof("Disconnected from websocket [%s]", serverName)
}
