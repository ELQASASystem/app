package websocket

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// handleQuestion 处理请求
func (w *wsSRV) handleQuestion(writer http.ResponseWriter, r *http.Request) {

	up := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	wsconn, err := up.Upgrade(writer, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}
	defer wsconn.Close()

	_, msg, err := wsconn.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("读取消息失败")
		return
	}

	i, _ := strconv.ParseUint(string(msg), 10, 32)
	id := uint32(i)

	w.connPool[id] = append(w.connPool[id], wsconn)
	defer w.rmConn(uint32(i), wsconn)

	log.Info().Uint64("问题ID", i).Msg("成功添加 WS 问题监听")

	for {
		if _, _, err := wsconn.ReadMessage(); err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
	}
}

// pushRemoteQA 向远程客户端推送问答数据
func (w *wsSRV) pushRemoteQA() {

	for {
		var (
			q     = <-w.ch.Question
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
