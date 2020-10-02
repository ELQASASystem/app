package class

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

// 使用 gorilla/websocket, 文档见 https://godoc.org/github.com/gorilla/websocket

// 设置服务器地址
var addr = flag.String("addr", "localhost:8081", "Websocket service address")

// Websocket 升级器, 用于将 Http 连接升级为 Websocket 处理
var upgrader = websocket.Upgrader{}

// echo 官方 Demo 用例, 简单的回声处理器
func echo(w http.ResponseWriter, r *http.Request) {

	// 将 HTTP 连接升级至 Websocket
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}
		log.Printf("收到消息: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error().Err(err).Msg("写入消息时出现问题")
			break
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/echo", echo)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Error().Err(err).Msg("Websocket 出现异常")
	}
}
