package mcqq

type GroupConfig struct {
	GroupId int64
	BotId   int64
}

type ServerConfig struct {
	GroupList []GroupConfig
	RconMsg   bool
}

type WebsocketClientConfig struct {
	ServerName string
	Url        string
}

type WebsocketServerConfig struct {
	Enable bool
	Host   string
	Port   int
}

type Config struct {
	WebsocketServer WebsocketServerConfig
	WebsocketClient []WebsocketClientConfig
	ServerMap       map[string]ServerConfig
	CommandPriority int
	ChatImage       bool
	AccessToken     string
}
