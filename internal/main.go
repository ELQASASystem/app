package internal

import (
	"github.com/ELQASASystem/app/internal/apis/http"
	"github.com/ELQASASystem/app/internal/apis/websocket"
	"github.com/ELQASASystem/app/internal/app"
)

func Main() {

	class.New(websocket.New())
	go http.StartAPI()

}
