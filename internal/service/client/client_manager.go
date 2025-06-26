package client

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/pkg/kafka"
)

type ClientManager struct {
	ClientMap        map[int64]*Client
	BroadcastChannel chan []byte
	PrivateChannel   chan []byte
	GroupChannel     chan []byte
	Register         chan *Client
	Unregister       chan *Client
	mu               sync.RWMutex
}

type IClientManager interface {
	Start()
	SetClient(client *Client)
	GetClient(id int64) (*Client, bool)
	RemoveClient(id int64)
	LaunchPrivateMessage(message []byte)
	LaunchBroadcastMessage(message []byte)
	LaunchGroupMessage(message []byte)
}

var (
	Manager = NewClientManager()
)

func NewClientManager() *ClientManager {
	return &ClientManager{
		ClientMap:        make(map[int64]*Client),
		BroadcastChannel: make(chan []byte),
		PrivateChannel:   make(chan []byte),
		GroupChannel:     make(chan []byte),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		mu:               sync.RWMutex{},
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

		case msg := <-cm.PrivateChannel:
			cm.LaunchPrivateMessage(msg)

		case msg := <-cm.BroadcastChannel:
			cm.LaunchBroadcastMessage(msg)

		case msg := <-cm.GroupChannel:
			cm.LaunchGroupMessage(msg)
		}
	}
}

func (cm *ClientManager) LaunchPrivateMessage(msg []byte) {
	// TODO: Implement private message handling
}

func (cm *ClientManager) LaunchBroadcastMessage(msg []byte) {
	// TODO: Implement broadcast message handling
}

func (cm *ClientManager) LaunchGroupMessage(msg []byte) {
	log.C(context.TODO()).Infow("Launching group message", "message", string(msg))
	var wsMsg message.WsMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.C(context.TODO()).Errorw("Failed to unmarshal group message", "error", err)
		return
	}

	group_Id := wsMsg.ToID
	channelType := wsMsg.ChannelType

	groupUser, err := store.S.GroupUser().List(context.TODO(), group_Id)
	if err != nil {
		log.C(context.TODO()).Errorw("Failed to get group users", "error", err)
		return
	}

	for _, user := range groupUser {
		receiveId := user.UserID
		if client, ok := cm.ClientMap[receiveId]; ok {
			client.Send <- msg
		} else {
			_, _, err := kafka.P.Push(msg, string(channelType), known.OfflineGroupTopic)
			if err != nil {
				log.C(context.TODO()).Errorw("Failed to push group message to Kafka", "error", err)
			} else {
				log.C(context.TODO()).Infow("Group message pushed to Kafka", "group_id", group_Id, "user_id", receiveId)
			}
		}
	}
}
