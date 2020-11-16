package websocket

import "github.com/rs/zerolog/log"

// sendQuestion 发送问题
func (w *srv) sendQuestion() {

	for {

		var (
			q     = <-w.qch
			conns = w.connPool[q.ID]
		)

		for _, v := range conns {

			err := v.WriteJSON(q)
			if err != nil {
				log.Error().Err(err).Msg("推送问题数据失败")
				continue
			}

			log.Info().Interface("客户端", v).Msg("推送数据中")
		}
	}
}
