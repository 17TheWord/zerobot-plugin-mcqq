package mcqq

import "encoding/json"

type BaseEvent struct {
	PostType   string `json:"post_type"`
	ServerName string `json:"server_name"`
	SubType    string `json:"sub_type"`
}

type APIResponse struct {
	Code     int         `json:"code"`
	API      string      `json:"api"`
	PostType string      `json:"post_type"`
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	Echo     string      `json:"echo,omitempty"`
}

type TranslateModel struct {
	Key  string        `json:"key"`
	Args []interface{} `json:"args"`
	Text string        `json:"text"`
}

type FlexibleTranslateText struct {
	Text      string
	Translate TranslateModel
}

func (f *FlexibleTranslateText) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		f.Text = text
		return nil
	}

	var translate TranslateModel
	if err := json.Unmarshal(data, &translate); err != nil {
		return err
	}
	f.Translate = translate
	return nil
}

func (f FlexibleTranslateText) String() string {
	if f.Translate.Text != "" {
		return f.Translate.Text
	}
	if f.Text != "" {
		return f.Text
	}
	return f.Translate.Key
}

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
	Key  string        `json:"key"`
	Args []interface{} `json:"args"`
	Text string        `json:"text"`
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
	Key       string         `json:"key"`
	Display   DisplayModel   `json:"display"`
	Text      string         `json:"text"`
	Translate TranslateModel `json:"translate"`
}

type DisplayModel struct {
	Title       FlexibleTranslateText `json:"title"`
	Description FlexibleTranslateText `json:"description"`
	Frame       string                `json:"frame"`
}
