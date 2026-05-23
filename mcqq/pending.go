package mcqq

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
)

const pendingResponseTimeout = 60 * time.Second

type pendingRequest struct {
	BotID      int64
	GroupID    int64
	ServerName string
	API        string
	Echo       string
}

type pendingRequestStore struct {
	mu       sync.Mutex
	requests map[string]pendingRequest
}

func newPendingRequestStore() *pendingRequestStore {
	return &pendingRequestStore{requests: make(map[string]pendingRequest)}
}

func (s *pendingRequestStore) add(request pendingRequest) {
	s.mu.Lock()
	s.requests[request.Echo] = request
	s.mu.Unlock()

	time.AfterFunc(pendingResponseTimeout, func() {
		if expired, ok := s.pop(request.Echo); ok {
			sendPendingResponse(expired, fmt.Sprintf("鹊桥 API 请求超时\n服务器: %s\n接口: %s", expired.ServerName, expired.API))
		}
	})
}

func (s *pendingRequestStore) pop(echo string) (pendingRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	request, ok := s.requests[echo]
	if ok {
		delete(s.requests, echo)
	}
	return request, ok
}

var pendingRequests = newPendingRequestStore()

func handleAPIResponse(response APIResponse) {
	request, ok := pendingRequests.pop(response.Echo)
	if !ok {
		return
	}
	sendPendingResponse(request, formatAPIResponse(response, request))
}

func sendPendingResponse(request pendingRequest, text string) {
	bot := zero.GetBot(request.BotID)
	if bot == nil {
		return
	}
	bot.SendGroupMessage(request.GroupID, text)
}

func formatAPIResponse(response APIResponse, request pendingRequest) string {
	var builder strings.Builder
	builder.WriteString("鹊桥 API 响应")
	builder.WriteString("\n服务器: ")
	builder.WriteString(request.ServerName)
	builder.WriteString("\n接口: ")
	if response.API != "" {
		builder.WriteString(response.API)
	} else {
		builder.WriteString(request.API)
	}
	builder.WriteString("\n状态: ")
	if response.Status != "" {
		builder.WriteString(response.Status)
	} else {
		builder.WriteString("UNKNOWN")
	}
	if response.Code != 0 {
		builder.WriteString(fmt.Sprintf("\n代码: %d", response.Code))
	}
	if response.Message != "" {
		builder.WriteString("\n消息: ")
		builder.WriteString(response.Message)
	}

	data := formatAPIResponseData(response.Data)
	if data != "" {
		builder.WriteString("\n数据: ")
		builder.WriteString(data)
	}

	return builder.String()
}

func formatAPIResponseData(data interface{}) string {
	if data == nil {
		return ""
	}
	if text, ok := data.(string); ok {
		return truncateResponseText(text)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return truncateResponseText(fmt.Sprint(data))
	}
	return truncateResponseText(string(jsonData))
}

func truncateResponseText(text string) string {
	const maxLen = 1500
	if len([]rune(text)) <= maxLen {
		return text
	}
	runes := []rune(text)
	return string(runes[:maxLen]) + "\n..."
}
