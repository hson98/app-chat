package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WsHandler struct {
	hub *Hub
}

func NewWsHandler(h *Hub) *WsHandler {
	return &WsHandler{hub: h}
}

type CreateRoomReq struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *WsHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	h.hub.Rooms[req.Id] = &Room{
		ID:      req.Id,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")
	cl := &Client{
		ID:      clientID,
		RoomId:  roomID,
		Conn:    conn,
		Message: make(chan *Message, 10),
	}
	m := &Message{
		Content:  fmt.Sprintln(username, " đã tham gia vào phòng"),
		RoomId:   roomID,
		Username: username,
	}
	h.hub.Register <- cl
	h.hub.Broadcast <- m
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *WsHandler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *WsHandler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
