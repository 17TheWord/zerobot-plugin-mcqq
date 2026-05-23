package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestLoadConfigWithDefaults(t *testing.T) {
	path := writeTempConfig(t, `
zero:
  drivers:
    - type: websocket_server
      id: 16
      url: ws://127.0.0.1:9090
mcqq:
  websocket_client:
    - server_name: Server
      url: ws://127.0.0.1:8080/minecraft/ws
  server_map:
    Server:
      group_list:
        - bot_id: 123456789
          group_id: 987654321
`)

	config, err := loadConfig(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Log.Level != "info" {
		t.Fatalf("expected default log level info, got %q", config.Log.Level)
	}
	if len(config.Zero.NickName) != 1 || config.Zero.NickName[0] != "bot" {
		t.Fatalf("unexpected default nick names: %#v", config.Zero.NickName)
	}
	if config.Zero.CommandPrefix != "/" {
		t.Fatalf("expected default command prefix /, got %q", config.Zero.CommandPrefix)
	}
	if config.Mcqq.CommandPriority != 2 {
		t.Fatalf("expected default command priority 2, got %d", config.Mcqq.CommandPriority)
	}
	if config.Mcqq.WebsocketServer.Host != "127.0.0.1" {
		t.Fatalf("expected default websocket server host, got %q", config.Mcqq.WebsocketServer.Host)
	}
	if config.Mcqq.WebsocketServer.Port != 8080 {
		t.Fatalf("expected default websocket server port 8080, got %d", config.Mcqq.WebsocketServer.Port)
	}
}

func TestValidateConfigRejectsMissingMinecraftConnection(t *testing.T) {
	path := writeTempConfig(t, `
zero:
  drivers:
    - type: websocket_server
      id: 16
      url: ws://127.0.0.1:9090
mcqq:
  server_map:
    Server:
      group_list:
        - bot_id: 123456789
          group_id: 987654321
`)

	_, err := loadConfig(path)
	if err == nil {
		t.Fatal("expected missing Minecraft connection error")
	}
	if !strings.Contains(err.Error(), "mcqq.websocket_client 至少需要配置一个连接") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConfigRejectsDuplicateGroupMapping(t *testing.T) {
	path := writeTempConfig(t, `
zero:
  drivers:
    - type: websocket_server
      id: 16
      url: ws://127.0.0.1:9090
mcqq:
  websocket_client:
    - server_name: ServerA
      url: ws://127.0.0.1:8080/minecraft/ws
  server_map:
    ServerA:
      group_list:
        - bot_id: 123456789
          group_id: 987654321
    ServerB:
      group_list:
        - bot_id: 123456789
          group_id: 987654321
`)

	_, err := loadConfig(path)
	if err == nil {
		t.Fatal("expected duplicate group mapping error")
	}
	if !strings.Contains(err.Error(), "同时映射到多个服务器") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConfigRejectsNegativeReconnect(t *testing.T) {
	path := writeTempConfig(t, `
zero:
  drivers:
    - type: websocket_server
      id: 16
      url: ws://127.0.0.1:9090
mcqq:
  websocket_client:
    - server_name: Server
      url: ws://127.0.0.1:8080/minecraft/ws
      reconnect_interval: -1
      reconnect_max_times: -1
  server_map:
    Server:
      group_list:
        - bot_id: 123456789
          group_id: 987654321
`)

	_, err := loadConfig(path)
	if err == nil {
		t.Fatal("expected negative reconnect config error")
	}
	if !strings.Contains(err.Error(), "reconnect_interval 不能小于 0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "reconnect_max_times 不能小于 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}
