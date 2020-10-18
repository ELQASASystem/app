package class

import (
	"github.com/ELQASASystem/app/internal/app/qq"

	"github.com/rs/zerolog/log"
	"strings"
)

// monitorGroup 监听群消息
func monitorGroup() {
	for {
		go processGroup(<-*Bot.MsgChan)
	}
}

// processGroup 处理群消息
func processGroup(m *qq.Msg) {

	if block(m) {
		return
	}

	if m.Chain[0].Text == ".hello" {
		Bot.SendGroupMsg(Bot.NewText("Hello, Client!").To(m.Group.ID))
		return
	}

	if strings.HasPrefix(m.Chain[0].Text, ".fenci ") {

		res, err := Bot.C.GetWordSegmentation(m.Chain[0].Text[7:])
		if err != nil {
			log.Error().Err(err).Msg("分词时出错")
			return
		}

		for k, v := range res {
			res[k] = strings.ReplaceAll(v, "\u0000", "")
		}

		Bot.SendGroupMsg(Bot.NewText(strings.Join(res, " | ")).To(m.Group.ID))
		return

	}

	if strings.HasPrefix(m.Chain[0].Text, ".tts ") {

		Bot.SendGroupMsg(Bot.NewTTSAudio(m.Chain[0].Text[5:]).To(m.Group.ID))
		return

	}

	handleAnswer(m) // 处理答案
}

// block 阻止可能的意外
func block(m *qq.Msg) bool {

	// 当长度小于1时消息无法获取
	if len(m.Chain) < 1 {
		return true
	}

	// 匿名用户禁止
	if m.User.ID == 80000000 {
		return true
	}

	return false

}
