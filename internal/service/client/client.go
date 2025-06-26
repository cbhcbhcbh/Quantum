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
	defer func() {
		Manager.Unregister <- c
		c.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		log.C(context.TODO()).Infow("Received message", "id", c.ID, "message", string(msg))

		msgString, ackMsg, err := message.ValidationMsg(msg)

		c.Conn.WriteMessage(websocket.TextMessage, []byte(ackMsg))
		if err == nil {
			Manager.BroadcastChannel <- []byte(msgString)
		}
	}
}

func (c *Client) Write() {
	defer c.Close()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.C(context.TODO()).Errorw("Error writing message", "id", c.ID, "error", err)
				return
			}
		}
	}
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
