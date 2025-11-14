package mcqq

// Player 玩家信息
type Player struct {
	Nickname           string  `json:"nickname"`
	UUID               string  `json:"uuid"`
	IsOP               bool    `json:"is_op"`
	Address            string  `json:"address"`
	Health             float64 `json:"health"`
	MaxHealth          float64 `json:"max_health"`
	ExperienceLevel    int     `json:"experience_level"`
	ExperienceProgress float64 `json:"experience_progress"`
	TotalExperience    int     `json:"total_experience"`
	WalkSpeed          float64 `json:"walk_speed"`
	X                  float64 `json:"x"`
	Y                  float64 `json:"y"`
	Z                  float64 `json:"z"`
}

// PlayerChatEvent 玩家聊天事件
type PlayerChatEvent struct {
	Timestamp     int    `json:"timestamp"`
	PostType      string `json:"post_type"`
	EventName     string `json:"event_name"`
	ServerName    string `json:"server_name"`
	ServerVersion string `json:"server_version"`
	ServerType    string `json:"server_type"`
	SubType       string `json:"sub_type"`
	MessageID     string `json:"message_id"`
	RawMessage    string `json:"raw_message"`
	Player        Player `json:"player"`
	Message       string `json:"message"`
}

// PlayerCommandEvent 玩家命令事件
type PlayerCommandEvent struct {
	Timestamp     int    `json:"timestamp"`
	PostType      string `json:"post_type"`
	EventName     string `json:"event_name"`
	ServerName    string `json:"server_name"`
	ServerVersion string `json:"server_version"`
	ServerType    string `json:"server_type"`
	SubType       string `json:"sub_type"`
	MessageID     string `json:"message_id"`
	RawMessage    string `json:"raw_message"`
	Player        Player `json:"player"`
	Command       string `json:"command"`
}

type PlayerNoticeEvent struct {
	Timestamp     int    `json:"timestamp"`
	PostType      string `json:"post_type"`
	EventName     string `json:"event_name"`
	ServerName    string `json:"server_name"`
	ServerVersion string `json:"server_version"`
	ServerType    string `json:"server_type"`
	SubType       string `json:"sub_type"`
	Player        Player `json:"player"`
}

type PlayerDeathEvent struct {
	Timestamp     int        `json:"timestamp"`
	PostType      string     `json:"post_type"`
	EventName     string     `json:"event_name"`
	ServerName    string     `json:"server_name"`
	ServerVersion string     `json:"server_version"`
	ServerType    string     `json:"server_type"`
	SubType       string     `json:"sub_type"`
	Player        Player     `json:"player"`
	Death         DeathModel `json:"death"`
}

type DeathModel struct {
	Key  string   `json:"key"`
	Args []string `json:"args"`
	Text string   `json:"text"`
}

type PlayerAchievementEvent struct {
	Timestamp     int              `json:"timestamp"`
	PostType      string           `json:"post_type"`
	EventName     string           `json:"event_name"`
	ServerName    string           `json:"server_name"`
	ServerVersion string           `json:"server_version"`
	ServerType    string           `json:"server_type"`
	SubType       string           `json:"sub_type"`
	Player        Player           `json:"player"`
	Achievement   AchievementModel `json:"achievement"`
}

type AchievementModel struct {
	Key     string       `json:"key"`
	Display DisplayModel `json:"display"`
	Text    string       `json:"text"`
}

type DisplayModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Frame       string `json:"frame"`
}
