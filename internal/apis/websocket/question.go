package websocket

import "github.com/rs/zerolog/log"

// pushRemoteQA 向远程客户端推送问答数据
func (w *srv) pushRemoteQA() {

	for {
		var (
			q     = <-w.qch
			conns = w.connPool[q.ID]
		)

		for _, v := range conns {
			if err := v.WriteJSON(q); err != nil {
				log.Error().Err(err).Str("客户端", v.RemoteAddr().String()).Msg("推送问题数据失败")
				continue
			}
		}
	}
}
