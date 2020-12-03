package websocket

import (
	"net/http"

	"github.com/ELQASASystem/server/internal/app"
	"github.com/ELQASASystem/server/internal/qq"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 文档： https://godoc.org/github.com/gorilla/websocket

type (
	// wsSRV Websocket 服务。
	wsSRV struct {
		ch *ch // ch 管道
		qa *qa // qa 问答
		ws *ws // ws 词云
	}

	// ch 管道
	ch struct {
		Question chan *app.Question // Question 问题管道
		WordStat chan *qq.Msg       // WordStat 词云管道
	}

	// qa 问答
	qa struct {
		pool map[uint32][]*websocket.Conn // pool 客户端池
		ch   *ch                          // ch 管道
	}

	// ws 词云
	ws struct {
		pool map[uint64][]*websocket.Conn // pool 客户端池
		ch   *ch                          // ch 管道
	}
)

// New 新建 Websocket 服务
func New() (w *wsSRV) {

	var (
		qch = make(chan *app.Question, 10)
		wch = make(chan *qq.Msg, 50)
		ch  = &ch{qch, wch}
	)

	w = &wsSRV{
		ch,
		&qa{map[uint32][]*websocket.Conn{}, ch},
		&ws{map[uint64][]*websocket.Conn{}, ch},
	}

	go w.start()
	return
}

// GetChannel 获取管道
func (w *wsSRV) GetChannel() *ch { return w.ch }

// start 启动 Websocket 服务
func (w *wsSRV) start() {

	go w.qa.pushQA()
	go w.ws.pushWordStat()

	http.HandleFunc("/question", w.qa.handleQuestion)
	http.HandleFunc("/stream/word_statistics", w.ws.handleWordStat)

	if err := http.ListenAndServe(":4041", nil); err != nil {
		log.Panic().Err(err).Msg("启动 Websocket API 服务 失败")
	}
}
