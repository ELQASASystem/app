package websocket

import (
	"github.com/ELQASASystem/server/internal/qq"
	"net/http"
	"strconv"

	"github.com/ELQASASystem/server/internal/app"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 文档： https://godoc.org/github.com/gorilla/websocket

type (
	// wsSRV Websocket 服务。
	wsSRV struct {
		connPool map[uint32][]*websocket.Conn // connPool 已连接的客户端
		ch       ch                           // ch 管道
	}

	// ch 管道
	ch struct {
		Question chan *app.Question // question 问题管道
		WordStat chan *qq.Msg       // wordStat 词云管道
	}
)

// New 新建 Websocket 服务
func New() (w *wsSRV) {

	var (
		qch = make(chan *app.Question, 10)
		wch = make(chan *qq.Msg, 50)
	)

	w = &wsSRV{
		map[uint32][]*websocket.Conn{},
		ch{qch, wch},
	}

	go w.start()
	return
}

// GetChannel 获取管道
func (w *wsSRV) GetChannel() *ch { return &w.ch }

// start 启动 Websocket 服务
func (w *wsSRV) start() {

	go w.pushRemoteQA()
	go w.sendWordStat()

	http.HandleFunc("/question", w.handleQuestion)
	http.HandleFunc("/stream/word_statistics", w.handleWordStat)

	if err := http.ListenAndServe(":4041", nil); err != nil {
		log.Panic().Err(err).Msg("启动 Websocket API 服务 失败")
	}
}

// handleWordStat 处理词云请求
func (w *wsSRV) handleWordStat(writer http.ResponseWriter, r *http.Request) {

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

	AddConn(i, wsconn)
	defer RmConn(i, wsconn)

	log.Info().Uint64("问题ID", i).Msg("成功添加词云词汇监听")

	for {
		if _, _, err := wsconn.ReadMessage(); err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
	}
}

// rmConn 使用 i：问题ID(ID) 移出一个连入的客户端
func (w *wsSRV) rmConn(i uint32, conn *websocket.Conn) {

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
