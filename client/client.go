package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/gedoy9793/blivedm-go/api"
	"github.com/gedoy9793/blivedm-go/packet"
	"github.com/gedoy9793/blivedm-go/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	conn                *websocket.Conn
	RoomID              int
	Uid                 int
	Buvid               string
	Cookie              string
	token               string
	host                string
	hostList            []string
	retryCount          int
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	cancel              context.CancelFunc
	done                <-chan struct{}
	Config              *Config
	ctx                 context.Context
}

type Config struct {
	UserAgent string
	Logger    logrus.FieldLogger
}

// NewClient 创建一个新的弹幕 client
func NewClient(roomID int, config *Config) *Client {
	if config.UserAgent == "" {
		config.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"
	}
	if config.Logger == nil {
		config.Logger = logrus.StandardLogger()
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = utils.SetLoggerToContext(ctx, config.Logger)
	return &Client{
		RoomID:              roomID,
		retryCount:          0,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
		Config:              config,
		ctx:                 ctx,
	}
}

func (c *Client) SetCookie(cookie string) {
	c.Cookie = cookie
}

// init 初始化 获取真实 RoomID 和 弹幕服务器 host
func (c *Client) init() error {
	if c.Cookie != "" {
		if !strings.Contains(c.Cookie, "bili_jct") || !strings.Contains(c.Cookie, "SESSDATA") {
			c.Config.Logger.Errorf("cannot found account token")
			return errors.New("账号未登录")
		}
		uid, err := api.GetUid(c.Cookie)
		if err != nil {
			c.Config.Logger.Error(err)
		}
		c.Uid = uid
		re := regexp.MustCompile("_uuid=(.+?);")
		result := re.FindAllStringSubmatch(c.Cookie, -1)
		if len(result) > 0 {
			c.Buvid = result[0][1]
		}
	}
	roomInfo, err := api.GetRoomInfo(c.RoomID)
	// 失败降级
	if err != nil || roomInfo.Code != 0 {
		c.Config.Logger.Errorf("room=%s init GetRoomInfo fialed, %s", c.RoomID, err)
	}
	c.RoomID = roomInfo.Data.RoomId
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.RoomID, c.Cookie)
		if err != nil {
			c.hostList = []string{"broadcastlv.chat.bilibili.com"}
		} else {
			for _, h := range info.Data.HostList {
				c.hostList = append(c.hostList, h.Host)
			}
		}
		c.token = info.Data.Token
	}
	return nil
}

func (c *Client) connect() error {
	reqHeader := &http.Header{}
	reqHeader.Set("User-Agent", c.Config.UserAgent)
retry:
	c.host = c.hostList[c.retryCount%len(c.hostList)]
	c.retryCount++
	conn, res, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/sub", c.host), *reqHeader)
	if err != nil {
		c.Config.Logger.Errorf("connect dial failed, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	c.conn = conn
	_ = res.Body.Close()
	if err = c.sendEnterPacket(); err != nil {
		c.Config.Logger.Errorf("failed to send enter packet, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	return nil
}

func (c *Client) wsLoop() {
	for {
		select {
		case <-c.done:
			c.Config.Logger.Debug("current client closed")
			return
		default:
			msgType, data, err := c.conn.ReadMessage()
			if err != nil {
				c.Config.Logger.Error("ws message read failed, reconnecting")
				time.Sleep(time.Duration(3) * time.Millisecond)
				_ = c.connect()
				continue
			}
			if msgType != websocket.BinaryMessage {
				c.Config.Logger.Error("packet not binary")
				continue
			}
			for _, fn := range c.eventHandlers.rawDataHandlers {
				go cover(c, func() { fn(&data) })
			}
			for _, pkt := range packet.DecodePacket(c.ctx, data).Parse() {
				go c.Handle(pkt)
			}
		}
	}
}

func (c *Client) heartBeatLoop() {
	pkt := packet.NewHeartBeatPacket(c.ctx)
	for {
		select {
		case <-c.done:
			return
		case <-time.After(30 * time.Second):
			if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
				c.Config.Logger.Error(err)
			}
			c.Config.Logger.Debug("send: HeartBeat")
		}
	}
}

// Start 启动弹幕 Client 初始化并连接 ws、发送心跳包
func (c *Client) Start() error {
	if err := c.init(); err != nil {
		return err
	}
	if err := c.connect(); err != nil {
		return err
	}
	go c.wsLoop()
	go c.heartBeatLoop()
	return nil
}

// Stop 停止弹幕 Client
func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) SetHost(host string) {
	c.host = host
}

// UseDefaultHost 使用默认 host broadcastlv.chat.bilibili.com
func (c *Client) UseDefaultHost() {
	c.hostList = []string{"broadcastlv.chat.bilibili.com"}
}

func (c *Client) sendEnterPacket() error {
	pkt := packet.NewEnterPacket(c.ctx, c.Uid, c.Buvid, c.RoomID, c.token)
	if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	c.Config.Logger.Debugf("send: EnterPacket")
	return nil
}
