package websocket

import (
	"net/http"
	"strconv"

	"github.com/ELQASASystem/server/internal/app"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// srv Websocket 服务。文档： https://godoc.org/github.com/gorilla/websocket
type srv struct {
	connPool map[uint32][]*websocket.Conn // connPool 已连接的客户端
	qch      chan *app.Question           // qch 问题管道
}

// New 新建 Websocket 服务
func New() chan *app.Question {

	var (
		qch = make(chan *app.Question, 10)
		s   = &srv{connPool: map[uint32][]*websocket.Conn{}, qch: qch}
	)

	go s.start()

	return qch
}

// start 启动 Websocket 服务
func (w *srv) start() {

	http.HandleFunc("/question", w.handle)
	http.HandleFunc("/stream/word_statistics", w.handleWordStat)
	go w.pushRemoteQA()
	go w.sendWordStat()

	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Error().Err(err).Msg("失败")
	}

}

// handle 处理请求
func (w *srv) handle(writer http.ResponseWriter, r *http.Request) {

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

// handleWordStat 处理请求
func (w *srv) handleWordStat(writer http.ResponseWriter, r *http.Request) {

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

	i, _ := strconv.ParseUint(string(msg), 10, 64) // 群号

	app.AddConn(i, wsconn)
	defer app.RmConn(i, wsconn)

	log.Info().Uint64("问题ID", i).Msg("成功添加词云词汇监听")

	for {
		if _, _, err := wsconn.ReadMessage(); err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
	}
}

// rmConn 使用 i：问题ID(ID) 移出一个连入的客户端
func (w *srv) rmConn(i uint32, conn *websocket.Conn) {

	conns := w.connPool[i]

	for k, wsconn := range conns {
		if wsconn == conn {
			w.connPool[i] = append(w.connPool[i][:k], w.connPool[i][k+1:]...)
		}
	}

	if len(conns) == 0 {
		delete(w.connPool, i)
	}

	log.Info().Msg("WS客户端下线")
}
