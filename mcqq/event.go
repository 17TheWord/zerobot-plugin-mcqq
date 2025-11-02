package mcqq

// BasePlayer ...
type BasePlayer struct {
	IsOp     *bool   `json:"is_op,omitempty"`
	Nickname string  `json:"nickname"`
	UUID     *string `json:"uuid,omitempty"`
}

// BaseMessageEvent ...
type BaseMessageEvent struct {
	EventName     string     `json:"event_name"`
	ServerType    string     `json:"server_type"`
	ServerVersion string     `json:"server_version"`
	Message       string     `json:"message"`
	MessageID     string     `json:"message_id"`
	Player        BasePlayer `json:"player"`
	PostType      string     `json:"post_type"`
	ServerName    string     `json:"server_name"`
	SubType       string     `json:"sub_type"`
	Timestamp     int64      `json:"timestamp"`
}

// BaseNoticeEvent ...
type BaseNoticeEvent struct {
	EventName  string     `json:"event_name"`
	Player     BasePlayer `json:"player"`
	PostType   string     `json:"post_type"`
	ServerName string     `json:"server_name"`
	SubType    string     `json:"sub_type"`
	Timestamp  int64      `json:"timestamp"`
}
