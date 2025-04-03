package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single websocket connection
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// Read pumps messages from the websocket connection to the manager
func (c *Client) Read(unregister chan<- *Client) {
	defer func() {
		unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// Write pumps messages from the manager to the websocket connection
func (c *Client) Write(unregister chan<- *Client) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The manager closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error writing message: %v", err)
				unregister <- c
				return
			}
		}
	}
}
