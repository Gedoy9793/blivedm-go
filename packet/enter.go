package packet

import (
	"encoding/json"
	"github.com/gedoy9793/blivedm-go/client"
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
func NewEnterPacket(c *client.Client, uid int, buvid string, roomID int, key string) []byte {
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
		c.Config.Logger.Error("NewEnterPacket JsonMarshal failed", err)
	}
	pkt := NewPlainPacket(c, RoomEnter, m)
	return pkt.Build()
}
