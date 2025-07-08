package client

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/cbhcbhcbh/Quantum/internal/apiserver/v1/store"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/model"
	"github.com/cbhcbhcbh/Quantum/pkg/kafka"
	"github.com/gorilla/websocket"
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
	ctx := context.Background()

	for {
		select {
		case client := <-cm.Register:
			cm.SetClient(client)
			cm.PullPrivateOfflineMessage(ctx, client)
			cm.ConsumingGroupOfflineMessages(ctx, client)
			log.C(ctx).Infow("Client registered", "id", client.ID)

		case client := <-cm.Unregister:
			cm.RemoveClient(client.ID)
			log.C(ctx).Infow("Client unregistered", "id", client.ID)

		case msg := <-cm.PrivateChannel:
			cm.LaunchPrivateMessage(ctx, msg)

		case msg := <-cm.BroadcastChannel:
			cm.LaunchBroadcastMessage(ctx, msg)

		case msg := <-cm.GroupChannel:
			cm.LaunchGroupMessage(ctx, msg)
		}
	}
}

func (cm *ClientManager) LaunchPrivateMessage(ctx context.Context, msg []byte) {
	log.C(ctx).Infow("Launching private message", "message", string(msg))
	var wsMsg message.WsMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.C(ctx).Errorw("Failed to unmarshal private message", "error", err)
		return
	}

	receiveId := wsMsg.ToID

	if client, ok := cm.ClientMap[receiveId]; ok {
		client.Send <- msg
	} else {
		_, _, err := kafka.P.Push(msg, strconv.FormatInt(receiveId, 10), known.OfflinePrivateTopic)
		if err != nil {
			log.C(ctx).Errorw("Failed to push private message to Kafka", "error", err)
		} else {
			log.C(ctx).Infow("private message pushed to Kafka", "user_id", receiveId)
		}
	}

}

func (cm *ClientManager) LaunchBroadcastMessage(ctx context.Context, msg []byte) {
	// TODO: Implement broadcast message handling
}

func (cm *ClientManager) LaunchGroupMessage(ctx context.Context, msg []byte) {
	log.C(ctx).Infow("Launching group message", "message", string(msg))
	var wsMsg message.WsMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.C(ctx).Errorw("Failed to unmarshal group message", "error", err)
		return
	}

	groupId := wsMsg.ToID
	channelType := wsMsg.ChannelType

	groupUser, err := store.S.GroupUser().List(ctx, groupId)
	if err != nil {
		log.C(ctx).Errorw("Failed to get group users", "error", err)
		return
	}

	for _, user := range groupUser {
		receiveId := user.UserID
		if client, ok := cm.ClientMap[receiveId]; ok {
			client.Send <- msg
		} else {
			_, _, err := kafka.P.Push(msg, strconv.Itoa(channelType), known.OfflineGroupTopic)
			if err != nil {
				log.C(ctx).Errorw("Failed to push group message to Kafka", "error", err)
			} else {
				log.C(ctx).Infow("Group message pushed to Kafka", "group_id", groupId, "user_id", receiveId)
			}
		}
	}
}

func (cm *ClientManager) PullPrivateOfflineMessage(ctx context.Context, client *Client) {
	pullAndPushOfflineMessages(
		ctx,
		client,
		func(ctx context.Context, start, end int64, status int16) ([]*model.OfflineMessageM, error) {
			return store.S.OfflineMessage().ListByTimeRangeAndStatus(ctx, start, end, status)
		},
		store.S.OfflineMessage().UpdateStatuByID,
		func(m *model.OfflineMessageM) ([]byte, int64) {
			return []byte(m.Message), m.ID
		},
	)
}

func (cm *ClientManager) ConsumingGroupOfflineMessages(ctx context.Context, client *Client) {
	pullAndPushOfflineMessages(
		ctx,
		client,
		func(ctx context.Context, start, end int64, status int16) ([]*model.GroupOfflineMessageM, error) {
			return store.S.GroupOfflineMessage().ListByTimeRangeAndStatus(ctx, start, end, status)
		},
		store.S.GroupOfflineMessage().UpdateStatuByID,
		func(m *model.GroupOfflineMessageM) ([]byte, int64) {
			return []byte(m.Message), m.ID
		},
	)
}

// TODO: 把这个函数重新放个地方
func pullAndPushOfflineMessages[T any](
	ctx context.Context,
	client *Client,
	listFunc func(context.Context, int64, int64, int16) ([]*T, error),
	updateFunc func(context.Context, []int64, int16) error,
	getMsg func(*T) ([]byte, int64),
) {
	nowTime := time.Now()
	lastTime := nowTime.Add(-15 * 24 * time.Hour) // 15 days ago

	var ids []int64
	var messages []*T
	var err error

	if messages, err = listFunc(ctx, lastTime.Unix(), nowTime.Unix(), 0); err != nil {
		log.C(ctx).Errorw("Pull offline message pushed to user", "user_id", client.ID, "error", err)
		return
	}

	for _, message := range messages {
		msgBytes, id := getMsg(message)
		_ = client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
		ids = append(ids, id)
	}

	if len(ids) > 0 {
		_ = updateFunc(ctx, ids, 1)
	}
}

func (cm *ClientManager) SendFriendActionMessage(msg message.CreateFriendMessage) bool {
	toId := msg.ToID
	message, _ := json.Marshal(msg)

	client, ok := cm.ClientMap[toId]
	if ok {
		client.Send <- message
		return true
	}
	return false
}
