package websocket

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// handleQuestion 处理请求
func (q *qa) handleQuestion(writer http.ResponseWriter, r *http.Request) {

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

	q.pool[id] = append(q.pool[id], wsconn)
	defer q.rmConn(uint32(i), wsconn)

	log.Info().Uint64("问题ID", i).Msg("成功添加 WS 问题监听")

	for {
		if _, _, err := wsconn.ReadMessage(); err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
	}
}

// pushQA 向远程客户端推送问答数据
func (q *qa) pushQA() {

	for {
		var (
			que   = <-q.ch.Question
			conns = q.pool[que.ID]
		)

		for _, v := range conns {
			if err := v.WriteJSON(que); err != nil {
				log.Error().Err(err).Str("客户端", v.RemoteAddr().String()).Msg("推送问题数据失败")
				continue
			}
		}
		log.Info().Uint32("问题ID", que.ID).Msg("问答推送数据")
	}
}

// rmConn 使用 i：问题ID(ID) 移出一个连入的客户端
func (q *qa) rmConn(i uint32, conn *websocket.Conn) {

	conns := q.pool[i]

	for k, wsconn := range conns {
		if wsconn == conn {
			q.pool[i] = append(q.pool[i][:k], q.pool[i][k+1:]...)
		}
	}

	if len(conns) == 0 {
		delete(q.pool, i)
	}

	log.Info().Msg("WS客户端下线")
}
