package app

import (
	"strings"

	"github.com/ELQASASystem/server/internal/qq"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var ls map[uint64][]*websocket.Conn // ls 监听词汇的客户端

func (a *App) handleWordCloud(m *qq.Msg) {

	// 不处理命令
	if strings.HasPrefix(m.Chain[0].Text, ".") {
		return
	}

	v, ok := ls[m.Group.ID]
	if !ok {
		return
	}

	res, err := a.Cli.C.GetWordSegmentation(m.Chain[0].Text)
	if err != nil {
		log.Error().Err(err).Msg("分词时出错")
		return
	}

	for _, conn := range v {

		err := conn.WriteJSON(res)
		if err != nil {
			log.Error().Err(err).Msg("推送词云数据失败")
		}

		log.Info().Interface("客户端", v).Msg("推送词云数据中")
	}
}

// AddConn 使用 gid：群ID 新增一个连入的客户端
func AddConn(gid uint64, c *websocket.Conn) {

	if _, ok := ls[gid]; !ok {
		ls[gid] = append(ls[gid], c)
	}

	ls[gid] = append(ls[gid], c)
}

// rmConn 使用 gid：群ID 移出一个连入的客户端
func RmConn(gid uint64, conn *websocket.Conn) {
	if pool, ok := ls[gid]; ok {
		for k, wsconn := range pool {
			if wsconn == conn {
				ls[gid] = append(ls[gid][:k], ls[gid][k+1:]...)
			}
		}
	}
}
