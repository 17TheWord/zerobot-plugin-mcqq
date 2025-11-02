package mcqq

// ========================
// 基础结构与枚举定义
// ========================

// ScoreComponent 表示计分板文本组件，显示玩家或实体的计分板数值。
type ScoreComponent struct {
	Name      string  `json:"name"`            // 计分板条目的名称（通常为玩家名或选择器）
	Objective string  `json:"objective"`       // 计分板目标名称
	Value     *string `json:"value,omitempty"` // 可选：直接指定显示值（通常仅用于客户端显示）
}

// ClickAction 点击事件类型
type ClickAction string

const (
	OpenURL         ClickAction = "open_url"
	OpenFile        ClickAction = "open_file"
	RunCommand      ClickAction = "run_command"
	SuggestCommand  ClickAction = "suggest_command"
	ChangePage      ClickAction = "change_page"
	CopyToClipboard ClickAction = "copy_to_clipboard"
)

// ClickEvent 点击事件定义，当玩家点击文本时触发的操作。
type ClickEvent struct {
	Action ClickAction `json:"action"` // 点击事件类型
	Value  string      `json:"value"`  // 参数值（命令、URL或页码）
}

// HoverAction 悬停事件类型
type HoverAction string

const (
	ShowText   HoverAction = "show_text"
	ShowItem   HoverAction = "show_item"
	ShowEntity HoverAction = "show_entity"
)

// HoverShowItem 鼠标悬停显示物品信息。
type HoverShowItem struct {
	ID    string      `json:"id"`              // 物品ID (例: "minecraft:diamond_sword")
	Count *int        `json:"count,omitempty"` // 物品数量
	Tag   interface{} `json:"tag,omitempty"`   // 物品NBT数据
}

// HoverShowEntity 鼠标悬停显示实体信息。
type HoverShowEntity struct {
	Name *Component `json:"name,omitempty"` // 实体显示名称，可为文本组件
	Type *string    `json:"type,omitempty"` // 实体类型ID (例: "minecraft:zombie")
	ID   *string    `json:"id,omitempty"`   // 实体UUID字符串
}

// HoverEvent 悬停事件定义，当鼠标悬停时显示额外信息。
type HoverEvent struct {
	Action   HoverAction `json:"action"`             // 悬停事件类型
	Contents interface{} `json:"contents,omitempty"` // 显示内容，可为文本、物品或实体信息
}

// Color 文本颜色枚举
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

// Component Minecraft 聊天文本组件 (Chat Component)
type Component struct {
	// === 内容字段 ===
	Text      *string       `json:"text,omitempty"`      // 纯文本内容
	Translate *string       `json:"translate,omitempty"` // 翻译键 (如 "chat.type.text")
	Fallback  *string       `json:"fallback,omitempty"`  // 翻译失败时的备用文本
	With      []interface{} `json:"with,omitempty"`      // 翻译参数，可为字符串或子组件

	Score     *ScoreComponent `json:"score,omitempty"`     // 计分板组件
	Selector  *string         `json:"selector,omitempty"`  // 实体选择器 (如 "@p")
	Separator *Component      `json:"separator,omitempty"` // 多实体分隔符文本
	Keybind   *string         `json:"keybind,omitempty"`   // 键位绑定名 (如 "key.jump")

	NBT     *string `json:"nbt,omitempty"`     // NBT路径表达式
	Block   *string `json:"block,omitempty"`   // NBT来源：方块坐标
	Entity  *string `json:"entity,omitempty"`  // NBT来源：实体选择器
	Storage *string `json:"storage,omitempty"` // NBT来源：存储命名空间

	// === 样式属性 ===
	Color         *Color  `json:"color,omitempty"`         // 文本颜色
	Font          *string `json:"font,omitempty"`          // 字体资源路径
	Bold          *bool   `json:"bold,omitempty"`          // 是否加粗
	Italic        *bool   `json:"italic,omitempty"`        // 是否斜体
	Underlined    *bool   `json:"underlined,omitempty"`    // 是否带下划线
	Strikethrough *bool   `json:"strikethrough,omitempty"` // 是否带删除线
	Obfuscated    *bool   `json:"obfuscated,omitempty"`    // 是否混淆显示

	// === 行为属性 ===
	Insertion  *string     `json:"insertion,omitempty"`  // Shift+点击插入到聊天栏
	ClickEvent *ClickEvent `json:"clickEvent,omitempty"` // 点击事件定义
	HoverEvent *HoverEvent `json:"hoverEvent,omitempty"` // 悬停事件定义

	// === 递归结构 ===
	Extra []*Component `json:"extra,omitempty"` // 附加子文本组件
}

type WebsocketData struct {
	API  string      `json:"api"`
	Data interface{} `json:"data"`
	Echo string      `json:"echo,omitempty"`
}
