package packet

import "github.com/Akegarasu/blivedm-go/client"

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket(c *client.Client) []byte {
	pkt := NewPacket(c, 1, HeartBeat, nil)
	return pkt.Build()
}
