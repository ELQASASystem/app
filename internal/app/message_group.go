package app

import (
	"github.com/ELQASASystem/app/internal/qq"

	"github.com/rs/zerolog/log"
	"strings"
)

// monitorGroup 监听群消息
func (a *App) monitorGroup() {
	for {
		go a.processGroup(<-a.mch)
	}
}

// processGroup 处理群消息
func (a *App) processGroup(m *qq.Msg) {

	if a.block(m) {
		return
	}

	if m.Chain[0].Text == ".hello" {
		a.Cli.SendGroupMsg(a.Cli.NewText("Hello, Client!").To(m.Group.ID))
		return
	}

	if strings.HasPrefix(m.Chain[0].Text, ".fenci ") {

		res, err := a.Cli.C.GetWordSegmentation(m.Chain[0].Text[7:])
		if err != nil {
			log.Error().Err(err).Msg("分词时出错")
			return
		}

		for k, v := range res {
			res[k] = strings.ReplaceAll(v, "\u0000", "")
		}

		a.Cli.SendGroupMsg(a.Cli.NewText(strings.Join(res, " | ")).To(m.Group.ID))
		return

	}

	if strings.HasPrefix(m.Chain[0].Text, ".tts ") {

		a.Cli.SendGroupMsg(a.Cli.NewTTSAudio(m.Chain[0].Text[5:]).To(m.Group.ID))
		return

	}

	a.handleAnswer(m) // 处理答案
}

// block 阻止可能的意外
func (a *App) block(m *qq.Msg) bool {

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
