package internal

import (
	"github.com/ELQASASystem/server/internal/apis/http"
	"github.com/ELQASASystem/server/internal/apis/websocket"
	"github.com/ELQASASystem/server/internal/app"
)

func Main() {

	w := websocket.New()

	app.New(w.GetChannel().Question, w.GetChannel().WordStat)
	go http.New()
}
