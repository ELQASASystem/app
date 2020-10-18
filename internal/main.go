package internal

import (
	"github.com/ELQASASystem/app/internal/apis/http"
	"github.com/ELQASASystem/app/internal/apis/websocket"
	"github.com/ELQASASystem/app/internal/app"
)

func Main() {
	class.New()
	go http.StartAPI()
	go websocket.New()
}
