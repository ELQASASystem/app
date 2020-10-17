package qq

import (
	"io/ioutil"

	m2 "github.com/Mrs4s/MiraiGo/message"
	"github.com/rs/zerolog/log"
)

// Message 构造消息
type Message struct {
	super  *Rina             // 父
	target uint64            // 目标
	chain  m2.SendingMessage // 消息链
}

// NewMsg 新建消息结构体
func (r *Rina) NewMsg() *Message { return &Message{super: r, chain: m2.SendingMessage{}} }

// NewText 新建文本消息结构体
func (r Rina) NewText(t string) *Message { return r.NewMsg().AddText(t) }

// NewImage 新建图片消息结构体
func (r Rina) NewImage(p string) *Message { return r.NewMsg().AddImage(p) }

// NewAudio 新建音频消息结构体
func (r Rina) NewAudio(p string) *Message { return r.NewMsg().AddAudio(p) }

// NewTTSAudio 新建文字转语音消息结构体
func (r Rina) NewTTSAudio(t string) *Message { return r.NewMsg().AddTTSAudio(t) }

// NewJSON 新建 JSON 卡片消息结构体
func (r Rina) NewJSON(s string) *Message { return r.NewMsg().AddJSON(s) }

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

// AddJSON 添加提醒
func (m *Message) AddAt(id uint64) *Message { m.chain.Append(m2.NewAt(int64(id))); return m }

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
func (m *Message) AddTTSAudio(t string) *Message {

	v, err := m.super.C.GetTts(t)
	if err != nil {
		log.Error().Err(err).Msg("文本转语音失败")
		return m
	}

	m.chain.Append(&m2.VoiceElement{Data: v})
	return m

}

// AddJSON 添加 JSON 卡片
func (m *Message) AddJSON(s string) *Message { m.chain.Append(m2.NewLightApp(s)); return m }

// To 发送的目标
func (m *Message) To(i uint64) *Message { m.target = i; return m }
