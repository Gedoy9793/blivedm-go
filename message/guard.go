package message

import (
	"context"
	"github.com/gedoy9793/blivedm-go/utils"
	"github.com/tidwall/gjson"
)

type GuardBuy struct {
	Uid        int    `json:"uid"`
	Username   string `json:"username"`
	GuardLevel int    `json:"guard_level"`
	Num        int    `json:"num"`
	Price      int    `json:"price"`
	GiftId     int    `json:"gift_id"`
	GiftName   string `json:"gift_name"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
}

func (g *GuardBuy) Parse(ctx context.Context, data []byte) {
	logger := utils.GetLoggerFromContext(ctx)
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, g)
	if err != nil {
		logger.Error("parse GuardBuy failed")
	}
}
