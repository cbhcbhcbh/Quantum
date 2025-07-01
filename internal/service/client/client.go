package client

import (
	"context"
	"sync"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/log"
	"github.com/cbhcbhcbh/Quantum/internal/pkg/message"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int64
	Conn *websocket.Conn
	Send chan []byte
	Mux  sync.RWMutex
}

type IClient interface {
	Read()
	Write()
	Close()
}

func NewClient(id int64, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) Read() {
	ctx := context.Background()

	defer func() {
		Manager.Unregister <- c
		c.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Close()
			break
		}

		log.C(ctx).Infow("Received message", "id", c.ID, "message", string(msg))

		msgByte, ackMsg, channel, err := message.ValidationMsg(msg)

		if err != nil {
			log.C(ctx).Errorw("Message validation error", "id", c.ID, "error", err)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msgByte)
		} else {
			switch channel {
			case message.PRIVATE:
				_ = c.Conn.WriteMessage(websocket.TextMessage, ackMsg)
				Manager.PrivateChannel <- msgByte
			case message.GROUP:
				_ = c.Conn.WriteMessage(websocket.TextMessage, ackMsg)
				Manager.GroupChannel <- msgByte
			case message.BROADCAST:
				_ = c.Conn.WriteMessage(websocket.TextMessage, ackMsg)
				Manager.BroadcastChannel <- msgByte
			case message.PING:
				_ = c.Conn.WriteMessage(websocket.TextMessage, ackMsg)
			default:
				log.C(ctx).Infow("Unknown channel type", "id", c.ID, "channel", channel)
			}
		}
	}
}

func (c *Client) Write() {
	ctx := context.Background()

	defer c.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.C(ctx).Errorw("Error writing message", "id", c.ID, "error", err)
			return
		}
	}

	_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) Close() {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	select {
	case <-c.Send:
	default:
		close(c.Send)
	}
	_ = c.Conn.Close()
}
