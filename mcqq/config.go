package mcqq

type Group struct {
	GroupId int64
	BotId   int64
}

type Server struct {
	GroupList []Group
	RconMsg   bool
	RconCmd   bool
}

type Config struct {
	Host            string
	Port            int
	ServerMap       map[string]Server
	CommandPriority int
	AccessToken     string
}
