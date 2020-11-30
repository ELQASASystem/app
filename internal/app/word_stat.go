package app

import (
	"strings"

	"github.com/ELQASASystem/server/internal/qq"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var ls = map[uint64][]*websocket.Conn{} // ls 监听词汇的客户端

// handleWordStat 处理词云
func (a *App) handleWordStat(m *qq.Msg) {

	v, ok := ls[m.Group.ID]
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

	for _, conn := range v {
		if err := conn.WriteJSON(words); err != nil {
			log.Error().Err(err).Str("客户端", conn.RemoteAddr().String()).Msg("推送词云数据失败")
		}
	}
}

// AddConn 使用 gid：群ID 新增一个连入的客户端
func AddConn(gid uint64, c *websocket.Conn) {
	ls[gid] = append(ls[gid], c)
}

// RmConn 使用 gid：群ID 移出一个连入的客户端
func RmConn(gid uint64, conn *websocket.Conn) {
	if pool, ok := ls[gid]; ok {
		for k, wsconn := range pool {
			if wsconn == conn {
				ls[gid] = append(ls[gid][:k], ls[gid][k+1:]...)
			}
		}
	}
}
