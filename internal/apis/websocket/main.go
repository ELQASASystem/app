package websocket

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 使用 gorilla/websocket, 文档见 https://godoc.org/github.com/gorilla/websocket

var ConnPool map[string][]*websocket.Conn // 已连接的客户端

// StartWebsocketAPI 启动 Websocket 服务器
func StartWebsocketAPI() error {

	http.HandleFunc("/q", connHandler)
	err := http.ListenAndServe(":4041", nil)

	if err != nil {
		return err
	}

	return nil

}

// addConn 新增一个连入的客户端
func addConn(qid string, conn *websocket.Conn) {

	if _, ok := ConnPool[qid]; !ok {
		ConnPool[qid] = []*websocket.Conn{conn}
		return
	}

	ConnPool[qid] = append(ConnPool[qid], conn)
}

// rmConn 移出一个连入的客户端
func rmConn(qid string, conn *websocket.Conn) bool {
	if conns, ok := ConnPool[qid]; ok {

		if len(conns) == 0 {
			delete(ConnPool, qid)
		}

		for i, wsconn := range conns {
			if wsconn == conn {
				ConnPool[qid] = append(ConnPool[qid][:i], ConnPool[qid][i+1:]...)
			}
		}
	}

	return false
}

// connHandler Websocket 连接处理器
func connHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级至 Websocket
	wsconn, err := new(websocket.Upgrader).Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}

	listenQid := 0

	defer wsconn.Close()

	// 处理客户端需要监听的问题
	go questionHandler(wsconn, listenQid)
}

func questionHandler(wsconn *websocket.Conn, listenQid int) {
	isRegistered := false

	for {
		mt, msg, err := wsconn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}

		// 获取从 Websocket 传入的字段
		action := string(msg)

		log.Debug().Msg("收到消息: " + action)

		result := "添加问题ID成功"

		// 获取传入字段是否为合法的问题ID
		// 目前仅做监听/取消监听操作
		if !isRegistered {
			if qid, err := strconv.Atoi(action); qid != 0 {
				listenQid = qid
				result = "成功添加对问题 " + strconv.Itoa(qid) + "的监听"
				isRegistered = true
			} else if err != nil {
				result = "无法解析需监听的问题ID"
			}

			// 向客户端发送操作结果
			err = wsconn.WriteMessage(mt, []byte(result))
			if err != nil {
				log.Error().Err(err).Msg("写入消息时出现问题")
				break
			}
		}
	}

	rmConn(strconv.Itoa(listenQid), wsconn)
}
