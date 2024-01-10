package packet

import (
	"context"
	"encoding/json"
	"github.com/gedoy9793/blivedm-go/utils"
)

type Enter struct {
	UID      int    `json:"uid"`
	Buvid    string `json:"buvid"`
	RoomID   int    `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

// NewEnterPacket 构造进入房间的包
// uid 可以为 0, key 在使用 broadcastlv 服务器的时候不需要
func NewEnterPacket(ctx context.Context, uid int, buvid string, roomID int, key string) []byte {
	logger := utils.GetLoggerFromContext(ctx)
	ent := &Enter{
		UID:      uid,
		Buvid:    buvid,
		RoomID:   roomID,
		ProtoVer: 3,
		Platform: "danmuji",
		Type:     2,
		Key:      key,
	}
	m, err := json.Marshal(ent)
	if err != nil {
		logger.Error("NewEnterPacket JsonMarshal failed", err)
	}
	pkt := NewPlainPacket(ctx, RoomEnter, m)
	return pkt.Build()
}
