package tictokapi

import (
	"encoding/json"
	"fmt"
	"github.com/caoyuewen/components/third/tictok/tictokmsg"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
)

const (
	ApiGetAccessToken = "https://developer.toutiao.com/api/apps/v2/token"          // 过去接口调用凭证
	ApiCode2Session   = "https://developer.toutiao.com/api/apps/v2/jscode2session" // 验证主播登录
	ApiGetLiveInfo    = "https://webcast.bytedance.com/api/webcastmate/info"       // 获取直播信息
	/*
		启动任务 数据开放 任务推送
		注意：不同类型的数据需要启动不同的任务单独监听，比如礼物数据单独启动一个，评论数据单独启动一个
		启动直播间数据推送，启动成功后，直播间数据才会同步推送给开发者服务器
		频率限制：单个 app_id 调用上限为 10 次/秒。
	*/
	ApiTaskStart = "https://webcast.bytedance.com/api/live_data/task/start" // 数据开放 任务推送 开始
	ApiTaskStop  = "https://webcast.bytedance.com/api/live_data/task/stop"  // 数据开放 任务推送 停止
	ApiGiftTop   = "https://webcast.bytedance.com/api/gift/top_gift"        // 礼物置顶

)

func Code2Session(appid, secret, code, anonymousCode string) (*tictokmsg.Code2SessionResp, error) {
	req := tictokmsg.Code2SessionReq{
		Appid:         appid,
		Secret:        secret,
		Code:          code,
		AnonymousCode: anonymousCode,
	}
	header := make(map[string]string)
	header["content-type"] = "application/json"
	body, err := Post(ApiCode2Session, req, header)
	if err != nil {
		log.Errorf("Code2Session err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.Code2SessionResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("Code2Session err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	return &resp, nil
}

func GetAccessToken(appid, secret string) (*tictokmsg.GetAccessTokenResp, error) {
	req := tictokmsg.GetAccessTokenReq{
		Appid:     appid,
		Secret:    secret,
		GrantType: "client_credential",
	}
	header := make(map[string]string)
	header["content-type"] = "application/json"
	body, err := Post(ApiGetAccessToken, req, header)
	if err != nil {
		log.Errorf("GetAccessToken err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.GetAccessTokenResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("GetAccessToken err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	return &resp, nil
}

func GetLiveInfo(accessToken, token string) (*tictokmsg.GetLiveInfoResp, error) {
	req := tictokmsg.GetLiveInfoReq{
		Token: token,
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["X-Token"] = accessToken
	body, err := Post(ApiGetLiveInfo, req, header)
	if err != nil {
		log.Errorf("GetLiveInfo err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.GetLiveInfoResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("GetLiveInfo err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	if resp.StatusCode != 0 || resp.Errmsg != "" {
		log.Errorf("GetLiveInfo err resp body:%s", string(body))
		return &resp, fmt.Errorf("GetLiveInfo err:%s", resp.Errmsg)
	}
	return &resp, nil
}

func TaskStart(accessToken, appid, roomid, msgType string) (*tictokmsg.TaskStartResp, error) {
	req := tictokmsg.TaskStartReq{
		Roomid:  roomid,
		Appid:   appid,
		MsgType: msgType,
	}
	header := make(map[string]string)
	header["content-type"] = "application/json"
	header["access-token"] = accessToken
	body, err := Post(ApiTaskStart, req, header)
	if err != nil {
		log.Errorf("TaskStart err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.TaskStartResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("TaskStart err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	return &resp, nil
}

func TaskStop(accessToken, appid, roomid, msgType string) (*tictokmsg.TaskStopResp, error) {
	req := tictokmsg.TaskStopReq{
		Roomid:  roomid,
		Appid:   appid,
		MsgType: msgType,
	}
	header := make(map[string]string)
	header["content-type"] = "application/json"
	header["access-token"] = accessToken
	body, err := Post(ApiTaskStop, req, header)
	if err != nil {
		log.Errorf("TaskStop err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.TaskStopResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("TaskStop err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	return &resp, nil
}

func GiftTop(accessToken, appid, roomid string, giftIdList []string) (*tictokmsg.GiftTopResp, error) {
	req := tictokmsg.GiftTopReq{
		RoomId:        roomid,
		AppId:         appid,
		SecGiftIdList: giftIdList,
	}
	header := make(map[string]string)
	header["content-type"] = "application/json"
	header["x-token"] = accessToken
	body, err := Post(ApiGiftTop, req, header)
	if err != nil {
		log.Errorf("GiftTop err:%s req :%s", err.Error(), util.ToJson(req))
		return nil, err
	}
	var resp tictokmsg.GiftTopResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf("GiftTop err:%s resp body:%s", err.Error(), string(body))
		return nil, err
	}
	return &resp, nil
}
