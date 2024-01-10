package main

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/Akegarasu/blivedm-go/client"
	_ "github.com/Akegarasu/blivedm-go/utils"
	"github.com/tidwall/gjson"
)

const roomId = 8792912

var dumps = []string{"GUARD_BUY", "USER_TOAST_MSG"}

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	c := client.NewClient(roomId, &client.Config{Logger: log})
	c.SetCookie("")
	for _, v := range dumps {
		vv := v
		c.RegisterCustomEventHandler(vv, func(s string) {
			data := gjson.Get(s, "data").String()
			fmt.Printf("[%s] %s\n", vv, data)
		})
	}
	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started bili dumper")
	select {}
}
