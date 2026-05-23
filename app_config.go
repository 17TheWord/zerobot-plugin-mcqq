package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/17TheWord/zerobot-plugin-mcqq/mcqq"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Log  LogConfig   `yaml:"log"`
	Zero ZeroConfig  `yaml:"zero"`
	Mcqq mcqq.Config `yaml:"mcqq"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type ZeroConfig struct {
	NickName      []string       `yaml:"nick_name"`
	CommandPrefix string         `yaml:"command_prefix"`
	SuperUsers    []int64        `yaml:"super_users"`
	Drivers       []DriverConfig `yaml:"drivers"`
}

type DriverConfig struct {
	Type        string `yaml:"type"`
	ID          int    `yaml:"id"`
	URL         string `yaml:"url"`
	AccessToken string `yaml:"access_token"`
}

func loadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 %q: %w", path, err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 YAML 配置文件失败 %q: %w", path, err)
	}

	applyConfigDefaults(&config)
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func applyConfigDefaults(config *AppConfig) {
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if len(config.Zero.NickName) == 0 {
		config.Zero.NickName = []string{"bot"}
	}
	if config.Zero.CommandPrefix == "" {
		config.Zero.CommandPrefix = "/"
	}
	if config.Mcqq.CommandPriority == 0 {
		config.Mcqq.CommandPriority = 2
	}
	if config.Mcqq.WebsocketServer.Host == "" {
		config.Mcqq.WebsocketServer.Host = "127.0.0.1"
	}
	if config.Mcqq.WebsocketServer.Port == 0 {
		config.Mcqq.WebsocketServer.Port = 8080
	}
}

func validateConfig(config *AppConfig) error {
	var errs []string

	switch strings.ToLower(config.Log.Level) {
	case "debug", "info", "warn", "warning", "error":
	default:
		errs = append(errs, "log.level 仅支持 debug、info、warn、error")
	}

	if len(config.Zero.Drivers) == 0 {
		errs = append(errs, "zero.drivers 不能为空")
	}
	for i, driverConfig := range config.Zero.Drivers {
		typePath := fmt.Sprintf("zero.drivers[%d].type", i)
		urlPath := fmt.Sprintf("zero.drivers[%d].url", i)

		switch driverConfig.Type {
		case "websocket_server":
			if driverConfig.ID == 0 {
				errs = append(errs, fmt.Sprintf("zero.drivers[%d].id 不能为空或 0", i))
			}
		case "websocket_client":
		default:
			errs = append(errs, fmt.Sprintf("%s 仅支持 websocket_server 或 websocket_client", typePath))
		}

		if driverConfig.URL == "" {
			errs = append(errs, fmt.Sprintf("%s 不能为空", urlPath))
		}
	}

	if len(config.Mcqq.ServerMap) == 0 {
		errs = append(errs, "mcqq.server_map 不能为空")
	}
	for serverName, server := range config.Mcqq.ServerMap {
		if strings.TrimSpace(serverName) == "" {
			errs = append(errs, "mcqq.server_map 不能包含空服务器名")
		}
		if len(server.GroupList) == 0 {
			errs = append(errs, fmt.Sprintf("mcqq.server_map.%s.group_list 不能为空", serverName))
		}
		for i, group := range server.GroupList {
			if group.BotId == 0 {
				errs = append(errs, fmt.Sprintf("mcqq.server_map.%s.group_list[%d].bot_id 不能为空或 0", serverName, i))
			}
			if group.GroupId == 0 {
				errs = append(errs, fmt.Sprintf("mcqq.server_map.%s.group_list[%d].group_id 不能为空或 0", serverName, i))
			}
		}
	}

	for i, websocketClient := range config.Mcqq.WebsocketClient {
		if websocketClient.ServerName == "" {
			errs = append(errs, fmt.Sprintf("mcqq.websocket_client[%d].server_name 不能为空", i))
		} else if _, ok := config.Mcqq.ServerMap[websocketClient.ServerName]; !ok {
			errs = append(errs, fmt.Sprintf("mcqq.websocket_client[%d].server_name %q 未在 mcqq.server_map 中配置", i, websocketClient.ServerName))
		}
		if websocketClient.Url == "" {
			errs = append(errs, fmt.Sprintf("mcqq.websocket_client[%d].url 不能为空", i))
		}
	}

	if config.Mcqq.WebsocketServer.Enable {
		if config.Mcqq.WebsocketServer.Host == "" {
			errs = append(errs, "mcqq.websocket_server.host 不能为空")
		}
		if config.Mcqq.WebsocketServer.Port <= 0 {
			errs = append(errs, "mcqq.websocket_server.port 必须大于 0")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("配置校验失败:\n- %s", strings.Join(errs, "\n- "))
	}
	return nil
}

func buildDrivers(configs []DriverConfig) ([]zero.Driver, error) {
	drivers := make([]zero.Driver, 0, len(configs))
	for i, config := range configs {
		switch config.Type {
		case "websocket_server":
			drivers = append(drivers, driver.NewWebSocketServer(config.ID, config.URL, config.AccessToken))
		case "websocket_client":
			drivers = append(drivers, driver.NewWebSocketClient(config.URL, config.AccessToken))
		default:
			return nil, fmt.Errorf("zero.drivers[%d].type 仅支持 websocket_server 或 websocket_client", i)
		}
	}
	return drivers, nil
}
