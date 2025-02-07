package mcqq

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func processQQMessage2MinecraftProtocol(ctx *zero.Ctx) WebsocketData {
	messageTextComponentList := make([]MessageSegment, len(ctx.Event.Message)+3)
	groupNameMap := map[int64]string{}
	groupName := "未知群聊"

	if value, exist := groupNameMap[ctx.Event.GroupID]; exist {
		groupName = value
	} else {
		groupInfo := ctx.GetGroupInfo(ctx.Event.GroupID, false)
		groupNameMap[ctx.Event.GroupID] = groupInfo.Name
		groupName = groupInfo.Name
	}
	messageTextComponentList[0] = MessageSegment{Type: "text", Data: TextComponent{Text: groupName}}

	nickname := ctx.Event.Sender.NickName
	if nickname == "" {
		nickname = ctx.Event.Sender.Card
	}

	messageTextComponentList[1] = MessageSegment{Type: "text", Data: TextComponent{Text: nickname}}

	messageTextComponentList[2] = MessageSegment{Type: "text", Data: TextComponent{Text: "说："}}

	for i := 0; i < len(ctx.Event.Message); i++ {
		var tempTextComponent TextComponent
		if ctx.Event.Message[i].Type == "text" {
			tempTextComponent = TextComponent{Text: ctx.Event.Message[i].Data["text"]}
		} else if ctx.Event.Message[i].Type == "image" {
			tempTextComponent = TextComponent{
				Text:  "[图片]",
				Color: Aqua,
				HoverEvent: &HoverEvent{
					Action: ShowText,
					Text:   []BaseComponent{{Text: fmt.Sprintf("[[CICode,url=%s,name=图片]]", ctx.Event.Message[i].Data["url"])}},
				},
			}
		} else {
			tempTextComponent = TextComponent{Text: "未知消息类型"}
		}
		messageTextComponentList[i+3] = MessageSegment{Type: "text", Data: tempTextComponent}
	}

	return WebsocketData{API: "send_msg", Data: messageTextComponentList}
}

func handleQQMessage(ctx *zero.Ctx) {
	protoMessage := processQQMessage2MinecraftProtocol(ctx)

	jsonData, jsonErr := json.Marshal(protoMessage)
	jsonString := string(jsonData)
	if jsonErr != nil {
		log.Errorln("Failed to convert protoMessage to json:", jsonErr)
		return
	}

	targetServername := getTargetServerName(ctx.Event.GroupID)
	if targetServername == "" {
		log.Errorf("Failed to get target server name with group by: %d", ctx.Event.GroupID)
		return
	}

	targetServer := getTargetServer(targetServername)
	if targetServer == nil {
		log.Errorf("Failed to get target server [%s] with group by: %d", targetServername, ctx.Event.GroupID)
		return
	}

	websocketErr := targetServer.Websocket.WriteJSON(jsonString)
	if websocketErr != nil {
		return
	}
	log.Infoln("Send message to Minecraft:", jsonString)
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
		break

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
		break

	default:
		log.Error("Unknown sub_type event from" + serverName + ": " + string(messageBytes))
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
