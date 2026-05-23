package mcqq

type GroupConfig struct {
	GroupId int64 `yaml:"group_id"`
	BotId   int64 `yaml:"bot_id"`
}

type ServerConfig struct {
	GroupList []GroupConfig `yaml:"group_list"`
	RconMsg   bool          `yaml:"rcon_msg"`
}

type WebsocketClientConfig struct {
	ServerName        string `yaml:"server_name"`
	Url               string `yaml:"url"`
	ReconnectInterval int    `yaml:"reconnect_interval"`
	ReconnectMaxTimes int    `yaml:"reconnect_max_times"`
}

type WebsocketServerConfig struct {
	Enable bool   `yaml:"enable"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type Config struct {
	WebsocketServer WebsocketServerConfig   `yaml:"websocket_server"`
	WebsocketClient []WebsocketClientConfig `yaml:"websocket_client"`
	ServerMap       map[string]ServerConfig `yaml:"server_map"`
	CommandPriority int                     `yaml:"command_priority"`
	ChatImage       bool                    `yaml:"chat_image"`
	AccessToken     string                  `yaml:"access_token"`
}
