package app

import (
	ws "github.com/ELQASASystem/server/internal/apis/websocket"
	"github.com/ELQASASystem/server/internal/qq"
	"github.com/rs/zerolog/log"
	"strings"
)

// handleWordStat 处理词云
func (a *App) handleWordStat(m *qq.Msg) {

	ok := ws.HasConn(m.Group.ID)
	if !ok {
		return
	}

	// 不处理命令和空消息
	if len(m.Chain[0].Text) == 0 || strings.HasPrefix(m.Chain[0].Text, ".") {
		return
	}

	words, err := DoWordSplit(m.Chain[0].Text)
	if err != nil {
		log.Error().Err(err).Msg("分词时出错")
		return
	}

	ws.PushWordStat(ws.WordStat{
		Gid:     m.Group.ID,
		Context: words,
	})
}
