package client

import (
	"sync"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
)

type ClientManager struct {
	ClientMap  map[int64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

type IClientManager interface {
	Start()
	SetClient(client *Client)
	GetClient(id int64) (*Client, bool)
	RemoveClient(id int64)
	LaunchMessage(message []byte)
}

var (
	Manager = NewClientManager()
)

func NewClientManager() *ClientManager {
	return &ClientManager{
		ClientMap:  make(map[int64]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		mu:         sync.RWMutex{},
	}
}

func (cm *ClientManager) SetClient(client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.ClientMap[client.ID] = client
}

func (cm *ClientManager) GetClient(id int64) (*Client, bool) {
	cm.Register <- &Client{ID: id}
	client, exists := cm.ClientMap[id]
	return client, exists
}

func (cm *ClientManager) RemoveClient(id int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.Unregister <- &Client{ID: id}
	delete(cm.ClientMap, id)
}

func (cm *ClientManager) Start() {
	for {
		select {
		case client := <-cm.Register:
			cm.SetClient(client)
			log.C(nil).Infow("Client registered", "id", client.ID)

		case client := <-cm.Unregister:
			cm.RemoveClient(client.ID)
			log.C(nil).Infow("Client unregistered", "id", client.ID)

		case message := <-cm.Broadcast:
			cm.LaunchMessage(message)
		}
	}
}

// TODO: LaunchMessage sends a message to all registered clients.
func (cm *ClientManager) LaunchMessage(message []byte) {

}
