package app

import (
	"github.com/ELQASASystem/server/internal/qq"

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

	if strings.HasPrefix(m.Chain[0].Text, ".tts ") {
		a.Cli.SendGroupMsg(a.Cli.NewTTSAudio(m.Chain[0].Text[5:]).To(m.Group.ID))
		return
	}

	a.handleAnswer(m) // 处理答案
	a.wch <- m        // 处理词云
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
