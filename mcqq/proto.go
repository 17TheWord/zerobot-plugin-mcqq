package mcqq

// WebsocketData  ...
type WebsocketData struct {
	API  string `json:"api"`
	Data any    `json:"data"`
}

type MessageData struct {
	Message []MessageSegment `json:"message"`
}

// MessageSegment ...
type MessageSegment struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

type Color string

const (
	Black       Color = "black"
	DarkBlue    Color = "dark_blue"
	DarkGreen   Color = "dark_green"
	DarkAqua    Color = "dark_aqua"
	DarkRed     Color = "dark_red"
	DarkPurple  Color = "dark_purple"
	Gold        Color = "gold"
	Gray        Color = "gray"
	DarkGray    Color = "dark_gray"
	Blue        Color = "blue"
	Green       Color = "green"
	Aqua        Color = "aqua"
	Red         Color = "red"
	LightPurple Color = "light_purple"
	Yellow      Color = "yellow"
	White       Color = "white"
)

// BaseComponent ...
type BaseComponent struct {
	// 加粗
	Bold *bool `json:"bold,omitempty"`
	// 颜色
	Color Color `json:"color,omitempty"`
	// 字体
	Font *string `json:"font,omitempty"`
	// 插入其他内容，未测试
	Insertion *string `json:"insertion,omitempty"`
	// 斜体
	Italic *bool `json:"italic,omitempty"`
	// 模糊的，未测试
	Obfuscated *bool `json:"obfuscated,omitempty"`
	// 删除线
	Strikethrough *bool `json:"strikethrough,omitempty"`
	// 文本
	Text string `json:"text"`
	// 下划线
	Underlined *bool `json:"underlined,omitempty"`
}

// ClickEventAction ...
type ClickEventAction string

const (
	ChangePage      ClickEventAction = "change_page"
	CopyToClipboard ClickEventAction = "copy_to_clipboard"
	OpenFile        ClickEventAction = "open_file"
	OpenURL         ClickEventAction = "open_url"
	RunCommand      ClickEventAction = "run_command"
	SuggestCommand  ClickEventAction = "suggest_command"
)

type HoverEventAction string

const (
	ShowEntity HoverEventAction = "show_entity"
	ShowItem   HoverEventAction = "show_item"
	ShowText   HoverEventAction = "show_text"
)

// ClickEvent ...
type ClickEvent struct {
	// 行为
	Action ClickEventAction `json:"action"`
	// 值
	Value string `json:"value,omitempty"`
}

// HoverItem ...
type HoverItem struct {
	Count int64  `json:"count,omitempty"`
	ID    string `json:"id,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

// HoverEntity ...
type HoverEntity struct {
	ID   string        `json:"id,omitempty"`
	Name BaseComponent `json:"name,omitempty"`
	Type string        `json:"type,omitempty"`
}

// HoverEvent ...
type HoverEvent struct {
	Action HoverEventAction `json:"action"`
	Entity *HoverEntity     `json:"entity,omitempty"`
	Item   *HoverItem       `json:"item,omitempty"`
	Text   []BaseComponent  `json:"text,omitempty"`
}

// TextComponent ...
type TextComponent struct {
	// 加粗
	Bold *bool `json:"bold,omitempty"`
	// 点击事件
	ClickEvent *ClickEvent `json:"click_event,omitempty"`
	// 颜色
	Color Color `json:"color,omitempty"`
	// 字体
	Font *string `json:"font,omitempty"`
	// 悬停事件
	HoverEvent *HoverEvent `json:"hover_event,omitempty"`
	// 插入其他内容，未测试
	Insertion *string `json:"insertion,omitempty"`
	// 斜体
	Italic *bool `json:"italic,omitempty"`
	// 模糊的，未测试
	Obfuscated *bool `json:"obfuscated,omitempty"`
	// 删除线
	Strikethrough *bool `json:"strikethrough,omitempty"`
	// 文本
	Text string `json:"text"`
	// 下划线
	Underlined *bool `json:"underlined,omitempty"`
}

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
