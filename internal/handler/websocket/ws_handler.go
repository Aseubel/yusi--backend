package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"yusi-backend/internal/svc"
	ws "yusi-backend/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境需要更严格的验证
	},
}

// WebSocketHandler 处理 WebSocket 连接
func WebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 路径获取房间代码
		roomCode := r.URL.Query().Get("roomCode")
		if roomCode == "" {
			// 尝试从路径参数获取
			roomCode = r.URL.Path[len("/api/ws/"):]
		}

		if roomCode == "" {
			http.Error(w, "房间代码不能为空", http.StatusBadRequest)
			return
		}

		// 从 Query 参数获取用户 ID
		userId := r.URL.Query().Get("userId")
		if userId == "" {
			http.Error(w, "用户ID不能为空", http.StatusUnauthorized)
			return
		}

		// 升级 HTTP 连接为 WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket 升级失败: %v", err)
			return
		}

		// 创建客户端
		client := &ws.Client{
			ID:     userId + "_" + roomCode,
			RoomID: roomCode,
			UserID: userId,
			Conn:   conn,
			Send:   make(chan []byte, 256),
		}

		// 注册客户端到 Hub
		svcCtx.WsHub.Register <- client

		// 启动读写协程
		go client.WritePump()
		go client.ReadPump(svcCtx.WsHub)
	}
}
