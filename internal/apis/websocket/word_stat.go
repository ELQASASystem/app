package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	wls = map[uint64][]*websocket.Conn{} // wls 监听词汇的客户端
	wch = make(chan WordStat, 10)
)

type WordStat struct {
	Gid     uint64
	Context []string
}

// sendQuestion 发送问题
func (w *srv) sendWordStat() {

	for {

		var (
			w     = <-wch
			conns = wls[w.Gid]
		)

		for _, v := range conns {

			err := v.WriteJSON(w)
			if err != nil {
				log.Error().Err(err).Msg("推送词云词汇失败")
				continue
			}

			log.Info().Interface("客户端", v).Msg("推送词云数据中")
		}
	}
}

// AddConn 使用 gid：群ID 新增一个连入的客户端
func AddConn(gid uint64, c *websocket.Conn) {
	wls[gid] = append(wls[gid], c)
}

// RmConn 使用 gid：群ID 移出一个连入的客户端
func RmConn(gid uint64, conn *websocket.Conn) {
	if pool, ok := wls[gid]; ok {
		for k, wsc := range pool {
			if wsc == conn {
				wls[gid] = append(wls[gid][:k], wls[gid][k+1:]...)
			}
		}
	}
}

func HasConn(gid uint64) (ok bool) {
	_, ok = wls[gid]
	return
}

func PushWordStat(w WordStat) { wch <- w }
