package client

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
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
	cm.mu.RLock()
	defer cm.mu.RUnlock()
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
			log.C(context.TODO()).Infow("Client registered", "id", client.ID)

		case client := <-cm.Unregister:
			cm.RemoveClient(client.ID)
			log.C(context.TODO()).Infow("Client unregistered", "id", client.ID)

		case msg := <-cm.Broadcast:
			cm.LaunchMessage(msg)
		}
	}
}

// TODO: LaunchMessage sends a message to all registered clients.
func (cm *ClientManager) LaunchMessage(msg []byte) {
	var wsMsg message.WsMessage

	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		return
	}

	channelType := wsMsg.ChannelType
	ReceiveId := wsMsg.ToID

	if channelType == message.ChannelTypePrivate || channelType == message.ChannelTypeGroup {

		if client, exists := cm.GetClient(ReceiveId); exists {
			// TODO: Perisist the message to the database or handle it accordingly
			log.C(context.TODO()).Infow("Sending message to client", "id", ReceiveId, "message", wsMsg.Message)
			client.Send <- msg
		} else {
			log.C(context.TODO()).Infow("Client not found", "id", ReceiveId)
			return
		}

	} else {
		// TODO: Handle broadcast messages
	}

}
