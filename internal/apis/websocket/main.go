package websocket

import (
	"net/http"
	"strconv"

	class "github.com/ELQASASystem/app/internal/app"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// srv Websocket 服务。文档： https://godoc.org/github.com/gorilla/websocket
type srv struct {
	connPool map[uint64][]*websocket.Conn // connPool 已连接的客户端
	qch      chan *class.Question         // qch 问题管道
}

// New 新建 Websocket 服务
func New() chan *class.Question {

	var (
		qch = make(chan *class.Question, 10)
		s   = &srv{connPool: map[uint64][]*websocket.Conn{}, qch: qch}
	)

	go s.start()

	return qch
}

// start 启动 Websocket 服务
func (w *srv) start() {

	http.HandleFunc("/question", func(writer http.ResponseWriter, r *http.Request) {

		upgrader := websocket.Upgrader{
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		wsconn, err := upgrader.Upgrade(writer, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
			return
		}
		defer wsconn.Close()

		w.questionHandler(wsconn)

	})

	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	go w.sendQuestion()

}

// addConn 使用 i：问题ID(ID) 新增一个连入的客户端
func (w *srv) addConn(i uint64, c *websocket.Conn) {

	if _, ok := w.connPool[i]; !ok {
		w.connPool[i] = []*websocket.Conn{c}
		return
	}

	w.connPool[i] = append(w.connPool[i], c)
}

// rmConn 移出一个连入的客户端
func (w *srv) rmConn(qid uint64, conn *websocket.Conn) {

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
func (w *srv) questionHandler(wsconn *websocket.Conn) {

	_, msg, err := wsconn.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("读取消息失败")
		return
	}

	action := string(msg)
	log.Info().Str("消息", action).Msg("收到 Websocket 消息")

	i, _ := strconv.ParseUint(action, 10, 64)

	log.Info().Uint64("问题ID", i).Msg("成功添加WS问题监听")

	w.addConn(i, wsconn)

}
