package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second

	// 读取超时时间
	pongWait = 60 * time.Second

	// ping 发送周期，必须小于 pongWait
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

// Client WebSocket 客户端
type Client struct {
	ID     string
	RoomID string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Message WebSocket 消息结构
type Message struct {
	Type    string      `json:"type"`    // join, leave, message, narrative_submit, room_start, room_finish
	RoomID  string      `json:"roomId"`
	UserID  string      `json:"userId"`
	Content interface{} `json:"content"`
}

// Hub 管理所有的 WebSocket 连接
type Hub struct {
	// 房间ID -> 客户端列表
	Rooms map[string]map[*Client]bool

	// 注册客户端
	Register chan *Client

	// 注销客户端
	Unregister chan *Client

	// 广播消息到指定房间
	Broadcast chan *Message

	mu sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			h.mu.Unlock()

			log.Printf("客户端 %s 加入房间 %s", client.UserID, client.RoomID)

			// 广播用户加入消息
			h.BroadcastToRoom(client.RoomID, &Message{
				Type:   "user_joined",
				RoomID: client.RoomID,
				UserID: client.UserID,
				Content: map[string]interface{}{
					"userId": client.UserID,
				},
			})

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					// 如果房间没有客户端了，删除房间
					if len(clients) == 0 {
						delete(h.Rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()

			log.Printf("客户端 %s 离开房间 %s", client.UserID, client.RoomID)

			// 广播用户离开消息
			h.BroadcastToRoom(client.RoomID, &Message{
				Type:   "user_left",
				RoomID: client.RoomID,
				UserID: client.UserID,
				Content: map[string]interface{}{
					"userId": client.UserID,
				},
			})

		case message := <-h.Broadcast:
			h.BroadcastToRoom(message.RoomID, message)
		}
	}
}

// BroadcastToRoom 向指定房间广播消息
func (h *Hub) BroadcastToRoom(roomID string, message *Message) {
	h.mu.RLock()
	clients := h.Rooms[roomID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.Send <- messageBytes:
		default:
			// 发送失败，关闭连接
			close(client.Send)
			h.mu.Lock()
			delete(clients, client)
			h.mu.Unlock()
		}
	}
}

// GetRoomMemberCount 获取房间在线人数
func (h *Hub) GetRoomMemberCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.Rooms[roomID]; ok {
		return len(clients)
	}
	return 0
}

// ReadPump 从 WebSocket 连接读取消息
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 读取错误: %v", err)
			}
			break
		}

		// 设置消息的房间ID和用户ID
		message.RoomID = c.RoomID
		message.UserID = c.UserID

		// 广播消息到房间
		hub.Broadcast <- &message
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的其他消息也写入
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
