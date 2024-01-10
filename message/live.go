package message

import (
	"context"
	"encoding/json"
	"github.com/gedoy9793/blivedm-go/utils"
)

type StopLiveRoomList struct {
	RoomIdList []int `json:"room_id_list"`
}

type Live struct {
	Cmd             string `json:"cmd"`
	LiveKey         string `json:"live_key"`
	VoiceBackground string `json:"voice_background"`
	SubSessionKey   string `json:"sub_session_key"`
	LivePlatform    string `json:"live_platform"`
	LiveModel       int    `json:"live_model"`
	LiveTime        int    `json:"live_time"`
	Roomid          int    `json:"roomid"`
}

type Preparing struct {
	Cmd    string `json:"cmd"`
	Roomid string `json:"roomid"`
}

func (l *Live) Parse(ctx context.Context, data []byte) {
	logger := utils.GetLoggerFromContext(ctx)
	err := json.Unmarshal(data, l)
	if err != nil {
		logger.Error("parse live failed")
	}
}
