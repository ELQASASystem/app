package internal

import (
	"github.com/ELQASASystem/server/internal/apis/http"
	"github.com/ELQASASystem/server/internal/apis/websocket"
	"github.com/ELQASASystem/server/internal/app"
)

func Main() {

	app.New(websocket.New())
	go http.New()

}
