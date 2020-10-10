package websocket

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 使用 gorilla/websocket, 文档见 https://godoc.org/github.com/gorilla/websocket

// 连接结构体
type NewConn struct {
	Conn      websocket.Conn // 连接
	uuid      string         // 客户端 UUID
	isActive  bool           // 是否活跃
	listenQID uint32         // 需监听问题的 ID
	Mt        int            // 消息类型
}

var connPool []NewConn // 连接用户储存池, key 为 UUID, value 为请求监听的答题ID

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
func addConn(wsconn *websocket.Conn, uuid string) (*NewConn, bool) {
	if _, _, ok := getConn(uuid); ok {
		return nil, false
	}

	newConn := NewConn{
		Conn:      *wsconn,
		uuid:      uuid,
		isActive:  true,
		listenQID: 0,
	}

	connPool = append(connPool, newConn)

	return &newConn, true
}

// remConn 移出一个连入的客户端
func remConn(uuid string) bool {
	if _, i, ok := getConn(uuid); !ok {
		return false
	} else {
		connPool = append(connPool[:i], connPool[i+1:]...)
		return true
	}
}

// getConn 通过 UUID 获取连接信息
func getConn(uuid string) (*NewConn, int, bool) {
	for i, ele := range connPool {
		if ele.uuid == uuid {
			return &ele, i, true
		}
	}

	return nil, -1, false
}

// getConnByQID 通过问题 ID 获取连接信息
func GetConnByQID(qid uint32) (*NewConn, int, bool) {
	for i, ele := range connPool {
		if ele.listenQID == qid {
			return &ele, i, true
		}
	}

	return nil, -1, false
}

// connHandler Websocket 连接处理器
func connHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级至 Websocket
	wsconn, err := new(websocket.Upgrader).Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}

	defer wsconn.Close()

	// 为本次连接生成一个新的 UUID
	u := uuid.New()
	conn, _ := addConn(wsconn, u.String())
	log.Log().Msg("客户端 " + u.String() + " 已连接")

	// 持续接收客户端传入的字段
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
		if qid, err := strconv.Atoi(action); qid != 0 {
			conn.listenQID = uint32(qid)
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

	remConn(u.String())
	log.Log().Msg("客户端 " + u.String() + " 已断开")
}
