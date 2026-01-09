package kafkas

import (
	"encoding/json"
)

type MsgKey struct {
	RoomId  string
	MsgType string
}

func GetMsgKey(roomId, msgType string) []byte {
	k := MsgKey{
		RoomId:  roomId,
		MsgType: msgType,
	}
	marshal, _ := json.Marshal(k)
	return marshal
}

func UnmarshalMsgKey(key []byte) (MsgKey, error) {
	var res MsgKey
	err := json.Unmarshal(key, &res)
	return res, err
}
