package message

import (
	"github.com/gedoy9793/blivedm-go/client"
	"github.com/gedoy9793/blivedm-go/utils"
	"github.com/tidwall/gjson"
)

type UserToast struct {
	AnchorShow       bool   `json:"anchor_show"`
	Color            string `json:"color"`
	Dmscore          int    `json:"dmscore"`
	EffectId         int    `json:"effect_id"`
	EndTime          int    `json:"end_time"`
	FaceEffectId     int    `json:"face_effect_id"`
	GiftId           int    `json:"gift_id"`
	GuardLevel       int    `json:"guard_level"`
	IsShow           int    `json:"is_show"`
	Num              int    `json:"num"`
	OpType           int    `json:"op_type"`
	PayflowId        string `json:"payflow_id"`
	Price            int    `json:"price"`
	RoleName         string `json:"role_name"`
	RoomEffectId     int    `json:"room_effect_id"`
	StartTime        int    `json:"start_time"`
	SvgaBlock        int    `json:"svga_block"`
	TargetGuardCount int    `json:"target_guard_count"`
	ToastMsg         string `json:"toast_msg"`
	Uid              int    `json:"uid"`
	Unit             string `json:"unit"`
	UserShow         bool   `json:"user_show"`
	Username         string `json:"username"`
}

func (u *UserToast) Parse(c *client.Client, data []byte) {
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, u)
	if err != nil {
		c.Config.Logger.Error("parse UserToast failed")
	}
}
