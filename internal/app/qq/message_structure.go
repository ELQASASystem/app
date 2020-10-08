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
func NewMsg() *Message { return &Message{chain: m2.SendingMessage{}} }

// NewText 新建文本消息结构体
func NewText(t string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddText(t) }

// NewImage 新建图片消息结构体
func NewImage(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddImage(p) }

// NewAudio 新建音频消息结构体
func NewAudio(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddAudio(p) }

// NewTTSAudio 新建文字转语音消息结构体
func NewTTSAudio(t string) *Message {
	m := &Message{chain: m2.SendingMessage{}}
	return m.AddTTSAudio(t)
}

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
func (m *Message) AddTTSAudio(t string) *Message {

	v, err := m.super.c.GetTts(t)
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
