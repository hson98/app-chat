package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 10),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, oke := h.Rooms[cl.RoomId]; oke {
				r := h.Rooms[cl.RoomId]
				if _, oke := r.Clients[cl.ID]; !oke {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomId]; ok {
				if _, ok := h.Rooms[cl.RoomId].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomId].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "auth left the chat",
							RoomId:   cl.RoomId,
							Username: cl.Username,
						}
					}

					delete(h.Rooms[cl.RoomId].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomId]; ok {

				for _, cl := range h.Rooms[m.RoomId].Clients {
					cl.Message <- m
				}
			}
		}

	}
}
