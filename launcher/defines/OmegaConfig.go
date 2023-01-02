package defines

// Omega 群服互通配置文件结构
type QGroupLink struct {
	Address                   string                        `json:"CQHTTP正向Websocket代理地址"`
	GameMessageFormat         string                        `json:"游戏消息格式化模版"`
	QQMessageFormat           string                        `json:"Q群消息格式化模版"`
	Groups                    map[string]int64              `json:"链接的Q群"`
	Selector                  string                        `json:"游戏内可以听到QQ消息的玩家的选择器"`
	NoBotMsg                  bool                          `json:"不要转发机器人的消息"`
	ChatOnly                  bool                          `json:"只转发聊天消息"`
	MuteIgnored               bool                          `json:"屏蔽其他群的消息"`
	FilterQQToServerMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到服务器"`
	FilterServerToQQMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到QQ"`
	AllowedCmdExecutor        map[int64]bool                `json:"允许这些人透过QQ执行命令"`
	AllowdFakeCmdExecutor     map[int64]map[string][]string `json:"允许这些人透过QQ执行伪命令"`
	DenyCmds                  map[string]string             `json:"屏蔽这些指令"`
	AllowCmds                 []string                      `json:"允许所有人使用这些指令"`
	SendJoinAndLeaveMsg       bool                          `json:"向Q群发送玩家进出消息"`
	ShowExchangeDetail        bool                          `json:"在控制台显示消息转发详情"`
}

type ComponentConfig struct {
	Name        string      `json:"名称"`
	Description string      `json:"描述"`
	Disabled    bool        `json:"是否禁用"`
	Version     string      `json:"版本"`
	Source      string      `json:"来源"`
	Configs     *QGroupLink `json:"配置"`
}
