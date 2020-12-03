package websocket

import (
	"net/http"
	"strconv"

	"github.com/ELQASASystem/server/internal/app"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// handleWordStat 处理词云请求
func (w *ws) handleWordStat(writer http.ResponseWriter, r *http.Request) {

	up := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	wsconn, err := up.Upgrade(writer, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}

	_, msg, err := wsconn.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("读取消息失败")
		return
	}

	i, err := strconv.ParseUint(string(msg), 10, 64) // 群号
	if err != nil {
		log.Error().Err(err).Msg("解析群号失败")
		return
	}

	w.pool[i] = append(w.pool[i], wsconn)
	defer w.RmConn(i, wsconn)

	log.Info().Uint64("问题ID", i).Msg("成功添加词云监听")

	for {
		if _, _, err := wsconn.ReadMessage(); err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
	}
}

// pushWordStat 推送词云数据
func (w *ws) pushWordStat() {

	for {
		m := <-w.ch.WordStat

		conns, ok := w.pool[m.Group.ID]
		if !ok {
			continue
		}

		// 不处理空消息
		if len(m.Chain[0].Text) == 0 {
			continue
		}

		words, err := app.DoWordSplit(m.Chain[0].Text)
		if err != nil {
			log.Error().Err(err).Msg("分词时出错")
			continue
		}

		for _, v := range conns {
			if err := v.WriteJSON(words); err != nil {
				log.Error().Err(err).Str("客户端", v.RemoteAddr().String()).Msg("推送词云数据失败")
				continue
			}
		}
		log.Info().Msg("推送词云数据")
	}
}

// RmConn 使用 gid：群ID 移出一个连入的客户端
func (w *ws) RmConn(gid uint64, conn *websocket.Conn) {
	if pool, ok := w.pool[gid]; ok {
		for k, wsc := range pool {
			if wsc == conn {
				w.pool[gid] = append(w.pool[gid][:k], w.pool[gid][k+1:]...)
			}
		}
	}
}
