package ws

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID       string `json:"id"`
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
	Conn     *websocket.Conn
	Message  chan *Message
}

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		m, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(m)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomId:   c.RoomId,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
