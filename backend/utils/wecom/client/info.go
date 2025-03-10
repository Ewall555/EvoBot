package client

type DeptUser struct {
	UserID     string `json:"userid"`
	Department int    `json:"department"`
}

type WecomConfig struct {
	Type           string `json:"type"`
	CorpID         string `json:"corp_id"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Secret         string `json:"secret"`
	AgentID        string `json:"agent_id"`
}

type SignatureOptions struct {
	Signature string `form:"msg_signature"`
	TimeStamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
	EchoStr   string `form:"echostr"`
}

type MessageInfo struct {
	MessageID   string `json:"messageid"`
	SendTime    uint64 `json:"send_time"`
	Origin      uint32 `json:"origin"`
	KFID        string `json:"kfid"`
	KHID        string `json:"khid"`
	StaffID     string `json:"staffid"`
	Message     string `json:"message"`
	MessageType string `json:"messagetype"`
	Credential  string `json:"credential"`
	ChatState   int    `json:"chatstate"`
}

type SendTextMsgOptions struct {
	KFID    string `json:"kfid"`
	KHID    string `json:"khid"`
	Message string `json:"message"`
}

type SendMenuMsgOptions struct {
	KFID string `json:"kfid"`
	KHID string `json:"khid"`
	MenuMsgOptions
}

type MenuMsgOptions struct {
	HeadContent string     `json:"headContent,omitempty"`
	MenuList    []MenuItem `json:"menuList,omitempty"`
	TailContent string     `json:"tailContent,omitempty"`
}

type SendTextMsgOnEventOptions struct {
	Message    string `json:"message"`
	Credential string `json:"credential"`
}

type SendMenuMsgOnEventOptions struct {
	MenuMsgOptions
	Credential string `json:"credential"`
}

type MenuItem struct {
	Type  string `json:"type"` // click, view, miniprogram, text 等类型
	Click *struct {
		ID      string `json:"id,omitempty"`
		Content string `json:"content"`
	} `json:"click,omitempty"`
	View *struct {
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"view,omitempty"`
	MiniProgram *struct {
		AppID    string `json:"appid"`
		PagePath string `json:"pagepath"`
		Content  string `json:"content"`
	} `json:"miniprogram,omitempty"`
	Text *struct {
		Content   string `json:"content"`
		NoNewLine int    `json:"no_newline,omitempty"`
	} `json:"text,omitempty"`
}

type ServiceStateTransOptions struct {
	OpenKFID       string `json:"open_kfid"`       // 客服帐号ID
	ExternalUserID string `json:"external_userid"` // 微信客户的external_userid
	ServicerUserID string `json:"servicer_userid"` // 接待人员的userid，当state=3时要求必填，接待人员须处于“正在接待”中
}

type ServiceStateGetOptions struct {
	OpenKFID       string `json:"open_kfid"`       // 客服帐号ID
	ExternalUserID string `json:"external_userid"` // 微信客户的external_userid
}

type ReceptionistOptions struct {
	OpenKFID   string   `json:"open_kfid"`   // 客服帐号ID
	UserIDList []string `json:"userid_list"` // 接待人员userid列表。第三方应用填密文userid，即open_userid 可填充个数：1 ~ 100。超过100个需分批调用。
}

type ReceptionistList struct {
	UserID string `json:"userid"` // 接待人员的userid。第三方应用获取到的为密文userid，即open_userid
	Status int    `json:"status"` // 接待人员的接待状态。0:接待中,1:停止接待。第三方应用需具有“管理帐号、分配会话和收发消息”权限才可获取
}

type AccountAddOptions struct {
	Name    string `json:"name"`     // 客服帐号名称, 不多于16个字符
	MediaID string `json:"media_id"` // 客服头像临时素材。可以调用上传临时素材接口获取, 不多于128个字节
}

type AccountUpdateOptions struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号ID, 不多于64字节
	Name     string `json:"name"`      // 客服帐号名称, 不多于16个字符
	MediaID  string `json:"media_id"`  // 客服头像临时素材。可以调用上传临时素材接口获取, 不多于128个字节
}

type AccountInfoSchema struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号ID
	Name     string `json:"name"`      // 客服帐号名称
	Avatar   string `json:"avatar"`    // 客服头像URL
}

// AddMsgTemplateRequest 创建企业群发请求
type AddMsgTemplateRequest struct {
	Sender      string   `json:"sender,omitempty"`
	Text        MsgText  `json:"text"`
	AllowSelect bool     `json:"allow_select,omitempty"`
	ChatIDList  []string `json:"chat_id_list,omitempty"`
}

// MsgText 文本消息
type MsgText struct {
	Content string `json:"content"`
}
