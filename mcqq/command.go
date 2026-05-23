package mcqq

import (
	"fmt"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func handleMCAPICommand(ctx *zero.Ctx) {
	args, _ := ctx.State["args"].(string)
	command, rest := cutField(strings.TrimSpace(args))
	if command == "" {
		ctx.Send(mcCommandHelp())
		return
	}

	switch command {
	case "status":
		handleStatusCommand(ctx, rest)
	case "rcon":
		handleRconCommand(ctx, rest)
	case "title":
		handleTitleCommand(ctx, rest)
	case "actionbar":
		handleActionBarCommand(ctx, rest)
	case "private":
		handlePrivateMessageCommand(ctx, rest)
	case "help":
		ctx.Send(mcCommandHelp())
	default:
		ctx.Send("未知 mc 命令: " + command + "\n" + mcCommandHelp())
	}
}

func handleStatusCommand(ctx *zero.Ctx, args string) {
	serverName, rest := cutField(args)
	if serverName == "" || rest != "" {
		ctx.Send("用法: /mc status <server>")
		return
	}
	sendMCAPIRequest(ctx, serverName, APIGetStatus, nil)
}

func handleRconCommand(ctx *zero.Ctx, args string) {
	serverName, command := cutField(args)
	if serverName == "" || command == "" {
		ctx.Send("用法: /mc rcon <server> <command>")
		return
	}
	sendMCAPIRequest(ctx, serverName, APISendRconCommand, map[string]interface{}{"command": command})
}

func handleTitleCommand(ctx *zero.Ctx, args string) {
	serverName, content := cutField(args)
	if serverName == "" || content == "" {
		ctx.Send("用法: /mc title <server> <title> [| <subtitle>]")
		return
	}

	title, subtitle, _ := strings.Cut(content, "|")
	title = strings.TrimSpace(title)
	subtitle = strings.TrimSpace(subtitle)
	if title == "" && subtitle == "" {
		ctx.Send("用法: /mc title <server> <title> [| <subtitle>]")
		return
	}

	data := map[string]interface{}{
		"fade_in":  20,
		"stay":     70,
		"fade_out": 20,
	}
	if title != "" {
		data["title"] = plainComponent(title)
	}
	if subtitle != "" {
		data["subtitle"] = plainComponent(subtitle)
	}
	sendMCAPIRequest(ctx, serverName, APISendTitle, data)
}

func handleActionBarCommand(ctx *zero.Ctx, args string) {
	serverName, content := cutField(args)
	if serverName == "" || content == "" {
		ctx.Send("用法: /mc actionbar <server> <message>")
		return
	}
	sendMCAPIRequest(ctx, serverName, APISendActionBar, map[string]interface{}{"message": []*Component{plainComponent(content)}})
}

func handlePrivateMessageCommand(ctx *zero.Ctx, args string) {
	serverName, rest := cutField(args)
	uuid, rest := cutField(rest)
	nickname, message := cutField(rest)
	if serverName == "" || uuid == "" || nickname == "" || message == "" {
		ctx.Send("用法: /mc private <server> <uuid|-> <nickname|-> <message>")
		return
	}
	if uuid == "-" {
		uuid = ""
	}
	if nickname == "-" {
		nickname = ""
	}
	if uuid == "" && nickname == "" {
		ctx.Send("uuid 和 nickname 至少需要填写一个")
		return
	}

	sendMCAPIRequest(ctx, serverName, APISendPrivateMsg, map[string]interface{}{
		"uuid":     uuid,
		"nickname": nickname,
		"message":  []*Component{plainComponent(message)},
	})
}

func sendMCAPIRequest(ctx *zero.Ctx, serverName string, api string, data interface{}) {
	conn, ok := mcConnections.get(serverName)
	if !ok {
		ctx.Send(fmt.Sprintf("服务器 %s 没有可用的 Minecraft WebSocket 连接", serverName))
		return
	}

	request := newWebsocketData(api, data)
	pendingRequests.add(pendingRequest{
		BotID:      ctx.Event.SelfID,
		GroupID:    ctx.Event.GroupID,
		ServerName: serverName,
		API:        api,
		Echo:       request.Echo,
	})

	if err := conn.WriteJSON(request); err != nil {
		pendingRequests.pop(request.Echo)
		ctx.Send(fmt.Sprintf("发送鹊桥 API 请求失败: %v", err))
	}
}

func plainComponent(text string) *Component {
	return &Component{Text: &text, Color: colorPtr(White)}
}

func cutField(text string) (string, string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ""
	}
	for i, r := range text {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return text[:i], strings.TrimSpace(text[i:])
		}
	}
	return text, ""
}

func mcCommandHelp() string {
	return strings.Join([]string{
		"用法:",
		"/mc status <server>",
		"/mc rcon <server> <command>",
		"/mc title <server> <title> [| <subtitle>]",
		"/mc actionbar <server> <message>",
		"/mc private <server> <uuid|-> <nickname|-> <message>",
	}, "\n")
}
