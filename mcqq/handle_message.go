package mcqq

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func colorPtr(c Color) *Color {
	return &c
}

// processQQMessageList 处理QQ消息列表，转换为Minecraft协议的Component列表
func processQQMessageList(ctx *zero.Ctx, message message.Message, replyModel bool) []*Component {
	messageList := make([]*Component, 0)
	for i := 0; i < len(message); i++ {
		msgType := message[i].Type
		msgData := message[i].Data

		var text string
		var color Color
		var hoverEvent *HoverEvent = nil
		var clickEvent *ClickEvent = nil
		var ciCode string

		if msgType == "reply" {
			text = "回复内容:\n\n"
			color = Gray
		} else if msgType == "text" {
			text = msgData["text"]
			color = White
		} else if msgType == "face" {
			text = "[表情]"
			color = Gold
			faceId := "表情ID: " + msgData["id"]
			hoverEvent = &HoverEvent{
				Action: "show_text",
				Contents: Component{
					Text:  &faceId,
					Color: colorPtr(DarkPurple),
				},
			}
		} else if msgType == "file" {
			text = "[文件]"
			color = Gold
			fileName := msgData["name"]
			hoverEvent = &HoverEvent{
				Action: "show_text",
				Contents: Component{
					Text:  &fileName,
					Color: colorPtr(DarkPurple),
				},
			}
		} else if msgType == "image" {
			url := msgData["url"]
			ciCode = "[[CICode,url=" + url + ",name=图片]]"

			text = "[图片]"
			color = LightPurple
			hoverText := "点击前往浏览器查看图片"
			hoverEvent = &HoverEvent{
				Action: "show_text",
				Contents: Component{
					Text:  &hoverText,
					Color: colorPtr(DarkPurple),
				},
			}
			clickEvent = &ClickEvent{
				Action: "open_url",
				Value:  url,
			}
		} else if msgType == "record" {
			text = "[语音]"
			color = Gold
		} else if msgType == "video" {
			text = "[视频]"
			color = LightPurple
			url := msgData["url"]
			hoverText := "点击前往浏览器查看视频"
			hoverEvent = &HoverEvent{
				Action: "show_text",
				Contents: Component{
					Text:  &hoverText,
					Color: colorPtr(DarkPurple),
				},
			}
			clickEvent = &ClickEvent{
				Action: "open_url",
				Value:  url,
			}
		} else if msgType == "at" {
			var name string
			if msgData["qq"] == "all" {
				name = "@所有人"
			} else {
				qqStr := msgData["qq"]
				qqInt, err := strconv.ParseInt(qqStr, 10, 64)
				if err == nil {
					// 优先从缓存获取
					if groupMap, ok := GroupMemberNameMap[ctx.Event.GroupID]; ok {
						if cachedName, ok2 := groupMap[qqInt]; ok2 && cachedName != "" {
							name = cachedName
						} else {
							groupMemberInfo := ctx.GetThisGroupMemberInfo(qqInt, false)
							if !groupMemberInfo.Exists() {
								name = "@" + qqStr
							} else {
								name = groupMemberInfo.Get("card").String()
								if name == "" {
									name = groupMemberInfo.Get("nickname").String()
								}
								// 写入缓存
								GroupMemberNameMap[ctx.Event.GroupID][qqInt] = name
							}
						}
					} else {
						// 初始化群成员缓存map
						GroupMemberNameMap[ctx.Event.GroupID] = map[int64]string{}
						groupMemberInfo := ctx.GetThisGroupMemberInfo(qqInt, false)
						if !groupMemberInfo.Exists() {
							name = "@" + qqStr
						} else {
							name = groupMemberInfo.Get("card").String()
							if name == "" {
								name = groupMemberInfo.Get("nickname").String()
							}
							// 写入缓存
							GroupMemberNameMap[ctx.Event.GroupID][qqInt] = name
						}
					}
				} else {
					name = "@" + qqStr
				}
			}
			text = "@" + name
			color = Green
		} else {
			text = "[" + msgType + "]"
			color = Gray
		}

		var component Component
		if replyModel {
			component = Component{Text: &text, Color: &color}
		} else {
			if PluginConfig.ChatImage {
				text = ciCode
				component = Component{Text: &text}
			} else {
				component = Component{Text: &text, Color: &color, HoverEvent: hoverEvent, ClickEvent: clickEvent}
			}
		}

		messageList = append(messageList, &component)
	}
	return messageList
}

func processQQMessage2MinecraftProtocol(ctx *zero.Ctx) []*Component {
	messageComponentList := make([]*Component, 3)

	var groupName string
	if value, exist := GroupNameMap[ctx.Event.GroupID]; exist {
		groupName = value
	} else {
		groupInfo := ctx.GetGroupInfo(ctx.Event.GroupID, false)
		GroupNameMap[ctx.Event.GroupID] = groupInfo.Name
		groupName = groupInfo.Name
	}

	groupName = " [" + groupName + "] "
	messageComponentList[0] = &Component{Text: &groupName, Color: colorPtr(Aqua)}

	nickname := ctx.Event.Sender.NickName
	if nickname == "" {
		nickname = ctx.Event.Sender.Card
	}

	messageComponentList[1] = &Component{Text: &nickname, Color: colorPtr(Green)}

	newMsg := make(message.Message, len(ctx.Event.Message))
	copy(newMsg, ctx.Event.Message)

	if newMsg[0].Type == "reply" {
		replyMsgId := newMsg[0].Data["id"]
		replyMsg := ctx.GetMessage(replyMsgId)
		replyUserName := replyMsg.Sender.Card
		if replyUserName == "" {
			replyUserName = replyMsg.Sender.NickName
		}
		replyComponentList := processQQMessageList(ctx, replyMsg.Elements, true)

		var replyText = " 回复 @" + replyUserName + " 的消息: "

		newMsg = newMsg[2:]
		qqMessageComponentList := processQQMessageList(ctx, newMsg, false)
		messageComponentList[2] = &Component{Text: &replyText, Color: colorPtr(Gray), HoverEvent: &HoverEvent{
			Action:   "show_text",
			Contents: replyComponentList,
		}}
		messageComponentList = append(messageComponentList, qqMessageComponentList...)
	} else {
		qqMessageComponentList := processQQMessageList(ctx, newMsg, false)
		var sayText = "说: "
		messageComponentList[2] = &Component{Text: &sayText, Color: colorPtr(White), Extra: qqMessageComponentList}
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
