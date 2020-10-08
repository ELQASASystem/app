package websocket

import (
	"flag"
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

var (
	// 设置服务器地址
	addr = flag.String("addr", "localhost:8081", "Websocket service address")
	// Websocket 升级器, 用于将 Http 连接升级为 Websocket 处理
	upgrader = websocket.Upgrader{}
	// 连接用户储存池, key 为 UUID, value 为请求监听的答题ID
	connPool []NewConn
)

func addConn(wsconn *websocket.Conn, uuid string) (*NewConn, bool) {
	for _, conn := range connPool {
		if conn.uuid == uuid {
			return nil, false
		}
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

func GetConnByQID(qid uint32) (*NewConn, bool) {
	for _, ele := range connPool {
		if ele.listenQID == qid {
			return &ele, true
		}
	}

	return nil, false
}

func connHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级至 Websocket
	wsconn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msg("处理 WebSocket 连接时出现异常")
		return
	}

	defer wsconn.Close()

	// 为本次连接生成一个新的 UUID
	u := uuid.New()
	conn, _ := addConn(wsconn, u.String())
	log.Log().Msg("客户端 " + u.String() + " 已连接")

	for {
		mt, msg, err := wsconn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("读取消息失败")
			break
		}

		action := string(msg)

		log.Debug().Msg("收到消息: " + action)

		result := "添加问题ID成功"

		if conn.listenQID == 0 {
			if qid, err := strconv.Atoi(action); qid != 0 {
				conn.listenQID = uint32(qid)
			} else if err != nil {
				result = "无法解析需监听的问题ID"
			}

			err = wsconn.WriteMessage(mt, []byte(result))
			if err != nil {
				log.Error().Err(err).Msg("写入消息时出现问题")
				break
			}
		}
	}

	log.Log().Msg("客户端 " + u.String() + " 已断开")
}

func main() {
	flag.Parse()
	http.HandleFunc("/q", connHandler)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Error().Err(err).Msg("Websocket 出现异常")
	}
}
