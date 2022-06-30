package ws

type server struct {
	Clients map[string]Client
}

func NewServer() Server {
	return &server{
		Clients: map[string]Client{},
	}
}

func (s *server) AddClient(client Client) {
	s.Clients[client.Id] = client
}

func (s *server) RemoveClient(client Client) {
	delete(s.Clients, client.Id)
}

func (s *server) Publish(data interface{}) {
	for _, c := range s.Clients {
		c.Conn.WriteJSON(data)
	}
}
