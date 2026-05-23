package mcqq

import (
	"strings"
	"testing"
)

func TestPendingRequestStorePop(t *testing.T) {
	store := newPendingRequestStore()
	request := pendingRequest{Echo: "echo-1", API: APIGetStatus, ServerName: "Server"}

	store.add(request)

	got, ok := store.pop("echo-1")
	if !ok {
		t.Fatal("expected pending request")
	}
	if got.API != APIGetStatus || got.ServerName != "Server" {
		t.Fatalf("unexpected pending request: %#v", got)
	}

	if _, ok := store.pop("echo-1"); ok {
		t.Fatal("pending request should be removed after pop")
	}
}

func TestFormatAPIResponse(t *testing.T) {
	response := APIResponse{
		Code:     200,
		API:      APIGetStatus,
		PostType: "response",
		Status:   "SUCCESS",
		Message:  "success",
		Data: map[string]interface{}{
			"server_type": "paper",
		},
		Echo: "echo-1",
	}
	request := pendingRequest{ServerName: "Server", API: APIGetStatus, Echo: "echo-1"}

	text := formatAPIResponse(response, request)
	for _, expected := range []string{"鹊桥 API 响应", "服务器: Server", "接口: get_status", "状态: SUCCESS", "server_type"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("response text should contain %q: %s", expected, text)
		}
	}
}

func TestTruncateResponseText(t *testing.T) {
	longText := strings.Repeat("一", 1600)
	truncated := truncateResponseText(longText)
	if !strings.HasSuffix(truncated, "\n...") {
		t.Fatalf("expected truncated suffix, got %q", truncated[len(truncated)-4:])
	}
}
