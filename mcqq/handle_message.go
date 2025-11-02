package mcqq

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func colorPtr(c Color) *Color {
	return &c
}

func processQQMessage2MinecraftProtocol(ctx *zero.Ctx) []Component {
	messageComponentList := make([]Component, len(ctx.Event.Message)+3)
	groupNameMap := map[int64]string{}

	var groupName string

	if value, exist := groupNameMap[ctx.Event.GroupID]; exist {
		groupName = value
	} else {
		groupInfo := ctx.GetGroupInfo(ctx.Event.GroupID, false)
		groupNameMap[ctx.Event.GroupID] = groupInfo.Name
		groupName = groupInfo.Name
	}

	groupName = " [" + groupName + "] "
	messageComponentList[0] = Component{Text: &groupName, Color: colorPtr(Aqua)}

	nickname := ctx.Event.Sender.NickName
	if nickname == "" {
		nickname = ctx.Event.Sender.Card
	}

	messageComponentList[1] = Component{Text: &nickname, Color: colorPtr(Green)}

	var sayText = "说: "
	messageComponentList[2] = Component{Text: &sayText, Color: colorPtr(White)}

	for i := 0; i < len(ctx.Event.Message); i++ {
		var text string
		var color Color
		if ctx.Event.Message[i].Type == "text" {
			text = ctx.Event.Message[i].Data["text"]
			color = White
		} else if ctx.Event.Message[i].Type == "image" {
			if PluginConfig.ChatImage {
				url := ctx.Event.Message[i].Data["url"]
				text = "[[CICode,url=" + url + ",name=图片]]"
			} else {
				text = "[图片]"
			}
			color = LightPurple
		} else {
			text = ctx.Event.Message[i].Type
			color = Gray
		}
		messageComponentList[i+3] = Component{Text: &text, Color: &color}
	}

	return messageComponentList
}

func handleQQMessage(ctx *zero.Ctx) {
	log.Info("正在处理来自QQ群 ", ctx.Event.GroupID, " 的消息...")
	protoMessage := processQQMessage2MinecraftProtocol(ctx)
	log.Info("处理消息完成，准备发送到Minecraft服务器...")

	messageData := map[string]interface{}{"message": protoMessage}

	timestamp := time.Now().UnixMilli()
	echoId := strconv.FormatInt(timestamp, 10)
	websocketData := WebsocketData{"send_msg", messageData, echoId}

	targetServerNameList := getTargetServerNameList(ctx.Event.GroupID)
	if len(targetServerNameList) == 0 {
		log.Errorf("No target server found for group: %d", ctx.Event.GroupID)
		return
	}

	targetServerList := getTargetServerWebsocketList(targetServerNameList)
	if len(targetServerList) == 0 {
		log.Errorf("No active websocket connection for group: %d", ctx.Event.GroupID)
		return
	}

	for _, targetServer := range targetServerList {
		websocketErr := targetServer.WriteJSON(websocketData)
		if websocketErr != nil {
			log.Errorln("Failed to send message to Minecraft server:", websocketErr)
		}
	}

}

func handleMinecraftMessage(messageBytes []byte) {
	var base map[string]interface{}
	err := json.Unmarshal(messageBytes, &base)
	if err != nil {
		log.Errorln("Error unmarshalling Minecraft message: ")
		log.Errorln(string(messageBytes))
		log.Errorln(err)
		return
	}

	postType := base["post_type"].(string)

	if postType == "response" {
		log.Info("接收到响应消息: " + string(messageBytes))
		return
	}

	serverName := base["server_name"].(string)
	subType := base["sub_type"].(string)

	var message = "[" + serverName + "] "

	switch subType {
	case "chat", "death":
		var messageEvent BaseMessageEvent

		err := json.Unmarshal(messageBytes, &messageEvent)
		if err != nil {
			log.Error("Error unmarshalling MessageEvent: ", err)
			return
		}

		if messageEvent.SubType == "chat" {
			message += messageEvent.Player.Nickname + " 说：" + messageEvent.Message
		} else {
			message += messageEvent.Message
		}

	case "join", "quit":
		var noticeEvent BaseNoticeEvent
		err := json.Unmarshal(messageBytes, &noticeEvent)
		if err != nil {
			log.Error("Error unmarshalling NoticeEvent: ", err)
			return
		}

		if noticeEvent.SubType == "join" {
			message += noticeEvent.Player.Nickname + " 加入了服务器"
		} else {
			message += noticeEvent.Player.Nickname + " 退出了服务器"
		}

	default:
		log.Error("Unsupported sub_type event from" + serverName + ": " + string(messageBytes))
		return
	}

	log.Infof("Received message from [%s]: %s", serverName, message)
	sendMcMsg2QQGroup(serverName, message)
}

func sendMcMsg2QQGroup(serverName string, message string) {
	if server, exists := PluginConfig.ServerMap[serverName]; exists {
		for _, group := range server.GroupList {
			bot := zero.GetBot(group.BotId)
			if bot == nil {
				log.Warningln("Failed to get bot with id: ", group.BotId)
				return
			}
			bot.SendGroupMessage(group.GroupId, message)
		}
	} else {
		log.Warningln("Failed to get server config with name: ", serverName)
	}
}
