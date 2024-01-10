package packet

import (
	"context"
)

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket(ctx context.Context) []byte {
	pkt := NewPacket(ctx, 1, HeartBeat, nil)
	return pkt.Build()
}
