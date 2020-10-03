package class

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	m2 "github.com/Mrs4s/MiraiGo/message"
	"github.com/rs/zerolog/log"
)

type (
	// Rina Rina QQ 客户端
	Rina struct {
		c       *client.QQClient // 客户端
		msgChan *chan *QQMsg     // 消息管道
	}

	// QQMsg 接收的 QQ 消息
	QQMsg struct {
		Chain []Chain           // 消息链
		Call  map[string]string // 参数
		Group struct {
			ID   uint64 // 群号
			Name string // 群名
		} // 群相关
		User struct {
			ID   uint64 // QQ号
			Name string // QQ名
		} // 用户相关
	}

	// Chain 消息链
	Chain struct {
		Type string // 类型：text、image、at
		Text string // text
		URL  string // image
		QQ   uint64 // at
	}

	// Message 返回的 QQ 消息
	Message struct {
		target uint64            // 目标
		chain  m2.SendingMessage // 消息链
	}
)

// newRina 新增 Rina
func newRina(i uint64, p string, ch *chan *QQMsg) (r *Rina) {

	err := client.SystemDeviceInfo.ReadJson([]byte("{\"display\":\"MIRAI.373480.001\",\"product\":\"mirai\",\"device\":\"mirai\",\"board\":\"mirai\",\"model\":\"mirai\",\"finger_print\":\"mamoe/mirai/mirai:10/MIRAI.200122.001/6671789:user/release-keys\",\"boot_id\":\"7794a02c-d854-18ac-649e-35fedfd0b37a\",\"proc_version\":\"Linux version 3.0.31-47Fxpwhn (android-build@xxx.xxx.xxx.xxx.com)\",\"protocol\":0,\"imei\":\"678319144775066\"}"))
	if err != nil {
		log.Panic().Err(err).Msg("设置设备信息失败")
	}
	client.SystemDeviceInfo.Protocol = client.AndroidPhone

	c := client.NewClient(int64(i), p)
	c.OnLog(func(q *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info().Str("信息", e.Message).Msg("协议")

		case "ERROR":
			log.Error().Str("信息", e.Message).Msg("协议")
		}
	})

	r = &Rina{c: c, msgChan: ch}
	if err := r.login(); err != nil {
		log.Panic().Msg("登录失败")
	}

	return

}

// login 登录
func (r Rina) login() (err error) {

	res, err := r.c.Login()
	if err != nil {
		return
	}
	for !res.Success {

		switch res.Error {
		case client.NeedCaptcha:
			_, err := r.c.SubmitCaptcha(r.needCap(res), res.CaptchaSign)
			if err != nil {
				log.Error().Err(err).Msg("提交验证码错误")
				continue
			}

		default:
			log.Panic().Str("原因", res.ErrorMessage).Msg("无法登录")
		}

	}

	log.Info().Msg("登录成功：" + r.c.Nickname)

	err = r.c.ReloadGroupList()
	if err != nil {
		log.Error().Err(err).Msg("加载群列表失败")
		return
	}

	err = r.c.ReloadFriendList()
	if err != nil {
		log.Error().Err(err).Msg("加载好友列表失败")
		return
	}

	log.Info().Int("个数", len(r.c.GroupList)).Msg("加载群列表成功")
	log.Info().Int("个数", len(r.c.FriendList)).Msg("加载好友列表成功")

	return

}

// needCap 触发输入验证码
func (r Rina) needCap(res *client.LoginResponse) string {

	file, err := os.Create("ca.jpg")
	if err != nil {
		log.Error().Err(err).Msg("创建验证码图片失败")
	}

	_, err = io.Copy(file, bytes.NewReader(res.CaptchaImage))
	if err != nil {
		log.Error().Err(err).Msg("写入验证码图片失败")
	}

	log.Info().Msg("请打开图片（ca.jpg）填写验证码")

	var c string
	if _, err := fmt.Scanln(&c); err != nil {
		log.Error().Err(err).Msg("读取错误，写的什么东西，爬")
	}

	return c

}

// regEventHandle 注册基本监听事件
func (r Rina) regEventHandle() {

	r.c.OnGroupMessage(r.onGroupMsg)
	r.c.OnPrivateMessage(r.onFriendMsg)

	// 断线重连
	r.c.OnDisconnected(func(q *client.QQClient, e *client.ClientDisconnectedEvent) {
		for {

			log.Warn().Msg("啊哦连接丢失了，准备重连中...1s")
			time.Sleep(time.Second)
			if err := r.login(); err != nil {
				log.Warn().Msg("重登录失败，再次尝试中...")
				continue
			}

			return

		}
	})

	// 更新服务器
	r.c.OnServerUpdated(func(q *client.QQClient, e *client.ServerUpdatedEvent) {
		log.Warn().Interface("数据", e.Servers).Msg("更新服务器")

		if len(e.Servers) < 1 {
			log.Error().Str("原因", "服务器地址长度为 0").Msg("更新服务器失败")
			return
		}

		var a []*net.TCPAddr
		for _, v := range e.Servers {
			a = append(a, &net.TCPAddr{
				IP:   net.ParseIP(v.Server),
				Port: int(v.Port),
			})
		}

		r.c.SetCustomServer(a)
	})

}

// onGroupMsg 触发群消息
func (r Rina) onGroupMsg(_ *client.QQClient, m *m2.GroupMessage) {

	msg := &QQMsg{
		Chain: []Chain{},
		Group: struct {
			ID   uint64
			Name string
		}{
			uint64(m.GroupCode),
			m.GroupName,
		},
		User: struct {
			ID   uint64
			Name string
		}{
			uint64(m.Sender.Uin),
			m.Sender.Nickname,
		},
	}

	for _, v := range m.Elements {
		switch e := v.(type) {
		case *m2.TextElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "text",
				Text: e.Content,
			})

		case *m2.AtElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "at",
				QQ:   uint64(e.Target),
			})

		case *m2.ImageElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "image",
				URL:  e.Url,
			})

		}
	}

	log.Info().
		Str("群名", msg.Group.Name).
		Str("昵称", msg.User.Name).
		Interface("原文", msg.Chain).
		Msg("收到群消息")

	*r.msgChan <- msg

}

// onFriendMsg 触发好友消息
func (r Rina) onFriendMsg(_ *client.QQClient, m *m2.PrivateMessage) {

	// TODO 好友消息

}

// NewMsg 新建消息结构体
func NewMsg() *Message { return &Message{chain: m2.SendingMessage{}} }

// NewText 新建文本消息结构体
func NewText(t string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddText(t) }

// NewImage 新建图片消息结构体
func NewImage(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddImage(p) }

// NewAudio 新建音频消息结构体
func NewAudio(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddAudio(p) }

// NewJSON 新建 JSON 卡片消息结构体
func NewJSON(s string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddJSON(s) }

// AddText 添加文本
func (m *Message) AddText(t string) *Message { m.chain.Append(m2.NewText(t)); return m }

// AddImage 添加图片
func (m *Message) AddImage(p string) *Message {

	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error().Err(err).Msg("读取图片失败")
		return m
	}

	m.chain.Append(m2.NewImage(b))

	return m

}

// AddAudio 添加音频
func (m *Message) AddAudio(p string) *Message {

	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error().Err(err).Msg("读取语音失败")
		return m
	}

	m.chain.Append(&m2.VoiceElement{Data: b})

	return m

}

// AddTTSAudio 添加 TTS 音频
func (m *Message) AddTTSAudio(text string) *Message {
	v, err := classBot.c.GetTts(text)

	if err != nil {
		log.Error().Err(err).Msg("转换文本为语音失败")
		return m
	}

	m.chain.Append(&m2.VoiceElement{Data: v})
	return m
}

// AddJSON 添加 JSON 卡片
func (m *Message) AddJSON(s string) *Message { m.chain.Append(m2.NewLightApp(s)); return m }

// To 发送的目标
func (m *Message) To(i uint64) *Message { m.target = i; return m }

// SendGroupMsg 发送群消息
func (r Rina) SendGroupMsg(m *Message) {

	for k, v := range m.chain.Elements {
		if nm, ok := v.(*m2.ImageElement); ok {
			am, err := r.c.UploadGroupImage(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传图片失败")
			} else {
				m.chain.Elements[k] = am
			}

		}

		if nm, ok := v.(*m2.VoiceElement); ok {
			am, err := r.c.UploadGroupPtt(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传语音失败")
			} else {
				m.chain.Elements[k] = am
			}
		}
	}

	log.Info().Msg("发送群消息")

	r.c.SendGroupMessage(int64(m.target), &m.chain)

}
