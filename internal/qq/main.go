package qq

import (
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	m2 "github.com/Mrs4s/MiraiGo/message"
	"github.com/rs/zerolog/log"
)

type (
	// Rina Rina QQ 客户端
	Rina struct {
		Message Message          // 消息
		C       *client.QQClient // 客户端
		MsgChan *chan *Msg       // 消息管道
	}

	// Msg 接收的 QQ 消息
	Msg struct {
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
)

// NewRina 新增 Rina
func NewRina(i uint64, p string, ch *chan *Msg) (r *Rina) {

	err := client.SystemDeviceInfo.ReadJson([]byte("{\"display\":\"MIRAI.373480.001\",\"product\":\"mirai\",\"device\":\"mirai\",\"board\":\"mirai\",\"model\":\"mirai\",\"finger_print\":\"mamoe/mirai/mirai:10/MIRAI.200122.001/6671789:user/release-keys\",\"boot_id\":\"7794a02c-d854-18ac-649e-35fedfd0b37a\",\"proc_version\":\"Linux version 3.0.31-47Fxpwhn (android-build@xxx.xxx.xxx.xxx.com)\",\"protocol\":0,\"imei\":\"678319144775066\"}"))
	if err != nil {
		log.Panic().Err(err).Msg("设置设备信息失败")
	}
	client.SystemDeviceInfo.Protocol = client.AndroidPhone

	c := client.NewClient(int64(i), p)
	c.OnLog(func(q *client.QQClient, e *client.LogEvent) {
		switch e.Type {
		case "INFO":
			log.Info().Str("信息", e.Message).Msg("协议信息")

		case "ERROR":
			log.Error().Str("信息", e.Message).Msg("协议错误")
		}
	})

	self := &Rina{C: c, MsgChan: ch}
	self.Message.super = self

	r = self
	if err := r.login(); err != nil {
		log.Panic().Msg("登录失败")
	}

	return

}

// login 登录
func (r Rina) login() (err error) {

	for res, err := r.C.Login(); err != nil || !res.Success; res, err = r.C.Login() {

		if err != nil {
			if err == client.ErrAlreadyOnline {
				return nil
			}

			log.Error().Err(err).Msg("登录失败")
			return err
		}

		switch res.Error {
		default:
			log.Panic().Str("原因", res.ErrorMessage).Msg("无法登录")
		}

	}

	log.Info().Msg("登录成功：" + r.C.Nickname)

	err = r.C.ReloadGroupList()
	if err != nil {
		log.Error().Err(err).Msg("加载群列表失败")
		return
	}

	err = r.C.ReloadFriendList()
	if err != nil {
		log.Error().Err(err).Msg("加载好友列表失败")
		return
	}

	log.Info().Int("个数", len(r.C.GroupList)).Msg("加载群列表成功")
	log.Info().Int("个数", len(r.C.FriendList)).Msg("加载好友列表成功")

	return

}

// RegEventHandle 注册基本监听事件
func (r Rina) RegEventHandle() {

	r.C.OnGroupMessage(r.onGroupMsg)
	r.C.OnPrivateMessage(r.onFriendMsg)

	// 断线重连
	r.C.OnDisconnected(func(q *client.QQClient, e *client.ClientDisconnectedEvent) {
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
	r.C.OnServerUpdated(func(q *client.QQClient, e *client.ServerUpdatedEvent) bool {
		log.Warn().Interface("数据", e.Servers).Msg("更新服务器")
		return true
	})

}

// onGroupMsg 触发群消息
func (r Rina) onGroupMsg(_ *client.QQClient, m *m2.GroupMessage) {

	msg := &Msg{
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

	log.Info().Str("群名", msg.Group.Name).Str("昵称", msg.User.Name).Interface("原文", msg.Chain).
		Msg("收到群消息")

	*r.MsgChan <- msg

}

// onFriendMsg 触发好友消息
func (r Rina) onFriendMsg(_ *client.QQClient, _ *m2.PrivateMessage) {

	// TODO 好友消息

}

// SendGroupMsg 发送群消息
func (r Rina) SendGroupMsg(m *Message) {

	for k, v := range m.chain.Elements {
		if nm, ok := v.(*m2.ImageElement); ok {
			am, err := r.C.UploadGroupImage(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传图片失败")
			} else {
				m.chain.Elements[k] = am
			}

		}

		if nm, ok := v.(*m2.VoiceElement); ok {
			am, err := r.C.UploadGroupPtt(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传语音失败")
			} else {
				m.chain.Elements[k] = am
			}
		}

		if nm, ok := v.(*m2.AtElement); ok {
			mem := r.C.FindGroupByUin(int64(m.target)).FindMember(nm.Target)
			if c := mem.CardName; c != "" {
				nm.Display = "@" + c
			} else {
				nm.Display = "@" + mem.Nickname
			}
			m.chain.Elements[k] = nm
		}

	}

	log.Info().Msg("发送群消息")
	r.C.SendGroupMessage(int64(m.target), &m.chain)
}
