package tictokmsg

type Code2SessionReq struct {
	Appid         string `json:"appid"`
	Secret        string `json:"secret"`
	AnonymousCode string `json:"anonymous_code"`
	Code          string `json:"code"`
}

type Code2SessionResp struct {
	ErrNo   int              `json:"err_no"`
	ErrTips string           `json:"err_tips"`
	Data    Code2SessionData `json:"data"`
}

type Code2SessionData struct {
	SessionKey      string `json:"session_key"`      // 会话密钥，如果请求时有 code 参数才会返回
	Openid          string `json:"openid"`           // 用户在当前小程序的 gen，如果请求时有 code 参数才会返回
	AnonymousOpenid string `json:"anonymous_openid"` // 匿名用户在当前小程序的 gen，如果请求时有 anonymous_code 参数才会返回
	Unionid         string `json:"unionid"`          // 用户在小程序平台的唯一标识符，请求时有 code 参数才会返回。如果开发者拥有多个小程序，可通过 unionid 来区分用户的唯一性。
}

// GetAccessTokenReq 获取接口调用凭证
type GetAccessTokenReq struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	GrantType string `json:"grant_type"`
}

type GetAccessTokenResp struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

// GetLiveInfoReq 获取直播房间信息
type GetLiveInfoReq struct {
	Token string `json:"token"`
}

type GetLiveInfoResp struct {
	Data    GetLiveInfoData `json:"data"`
	Errcode int             `json:"errcode"`
	Errmsg  string          `json:"errmsg"`
	Extra   struct {
		Now int64 `json:"now"`
	} `json:"extra"`
	StatusCode int `json:"status_code"`
}

type GetLiveInfoData struct {
	AckCfg     []interface{}  `json:"ack_cfg"`
	LinkerInfo LinkerInfoData `json:"linker_info"`
	Info       InfoData       `json:"info"`
}

type LinkerInfoData struct {
	LinkerId     int `json:"linker_id"`
	LinkerStatus int `json:"linker_status"`
	MasterStatus int `json:"master_status"`
}

type InfoData struct {
	RoomId       int    `json:"room_id"`
	AnchorOpenId string `json:"anchor_open_id"`
	AvatarUrl    string `json:"avatar_url"`
	NickName     string `json:"nick_name"`
}

// TaskStartReq 数据开放 任务推送
type TaskStartReq struct {
	Roomid  string `json:"roomid"`
	Appid   string `json:"appid"`
	MsgType string `json:"msg_type"`
	/*
		直播间消息类型，需要前置申请开通了对应类型直播间数据能力才可以调用。
		1. 评论：live_comment
		2. 礼物：live_gift
		3. 点赞：live_like
		4. 粉丝团：live_fansclub
	*/
}

type TaskStartResp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Logid  string `json:"logid"`
	Data   struct {
		TaskId string `json:"task_id"`
	} `json:"data"`
}

// TaskStopReq 数据开放 任务推送 停止
type TaskStopReq struct {
	Roomid  string `json:"roomid"`
	Appid   string `json:"appid"`
	MsgType string `json:"msg_type"`
	/*
		直播间消息类型，需要前置申请开通了对应类型直播间数据能力才可以调用。
		1. 评论：live_comment
		2. 礼物：live_gift
		3. 点赞：live_like
		4. 粉丝团：live_fansclub
	*/
}

type TaskStopResp struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Logid  string      `json:"logid"`
	Data   interface{} `json:"data"`
}

// GiftTopReq 礼物置顶
type GiftTopReq struct {
	RoomId        string   `json:"room_id"`
	AppId         string   `json:"app_id"`
	SecGiftIdList []string `json:"sec_gift_id_list"`
}

type GiftTopResp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Logid  string `json:"logid"`
	Data   struct {
		SuccessTopGiftIdList []string `json:"success_top_gift_id_list"`
	} `json:"data"`
}
