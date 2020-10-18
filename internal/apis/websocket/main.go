package websocket

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 文档： https://godoc.org/github.com/gorilla/websocket

type websocketSrv struct {
	connPool map[uint64][]*websocket.Conn // 已连接的客户端
}

// New 新建
func New() *websocketSrv {
	return &websocketSrv{}
}

// start 启动 Websocket 服务
func (w *websocketSrv) start() (err error) {

	http.HandleFunc("/question", func(writer http.ResponseWriter, r *http.Request) {

		wsconn, err := new(websocket.Upgrader).Upgrade(writer, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
			return
		}
		defer wsconn.Close()

		go w.questionHandler(wsconn)

	})
	err = http.ListenAndServe(":4041", nil)

	if err != nil {
		return
	}

	return
}

// addConn 使用 i：问题ID(ID) 新增一个连入的客户端
func (w *websocketSrv) addConn(i uint64, c *websocket.Conn) {

	if _, ok := w.connPool[i]; !ok {
		w.connPool[i] = []*websocket.Conn{c}
		return
	}

	w.connPool[i] = append(w.connPool[i], c)
}

// rmConn 移出一个连入的客户端
func (w *websocketSrv) rmConn(qid uint64, conn *websocket.Conn) {

	conns := w.connPool[qid]

	for k, wsconn := range conns {
		if wsconn == conn {
			w.connPool[qid] = append(w.connPool[qid][:k], w.connPool[qid][k+1:]...)
		}
	}

	if len(conns) == 0 {
		delete(w.connPool, qid)
	}

}

// questionHandler 问题处理器
func (w *websocketSrv) questionHandler(wsconn *websocket.Conn) {

	for {
		_, msg, err := wsconn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}

		action := string(msg)
		log.Info().Str("消息", action).Msg("收到 Websocket 消息")

		// 获取传入字段是否为合法的问题ID
		// 目前仅做监听/取消监听操作

		qid, err := strconv.ParseUint(action, 10, 64)
		if err != nil {
			break
		}

		log.Info().Uint64("问题ID", qid).Msg("成功添加WS问题监听")

		// 向客户端发送操作结果
		/*
			err = wsconn.WriteMessage(mt, []byte(result))
			if err != nil {
				log.Error().Err(err).Msg("写入消息时出现问题")
				break
			}
		*/

	}

	w.rmConn(uint64(0), wsconn)
}
