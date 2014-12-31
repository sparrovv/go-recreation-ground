package mdprev

import (
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func newHub(broadcast chan []byte, exit chan bool) *hub {
	return &hub{
		broadcast:   broadcast,
		exit:        exit,
		register:    make(chan connector),
		unregister:  make(chan connector),
		connections: make(map[connector]bool),
	}
}

type connector interface {
	writer()
	closeCh()
	sendMsg(msg []byte)
}

// inspired by: http://gary.burd.info/go-websocket-chat
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Only listen to get EOF, so we can remove ws connection
func (c *connection) unregisterOnEOF(h *hub) {
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
	h.unregister <- c
	c.ws.Close()
}

func (c *connection) closeCh() {
	close(c.send)
}

func (c *connection) sendMsg(msg []byte) {
	c.send <- msg
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

type hub struct {
	// Registered connections.
	connections map[connector]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan connector

	// Unregister requests from connections.
	unregister chan connector

	// send when all connections are closed after unregistering the last one
	exit chan bool
}

// all the routing of messages happen here
func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				c.closeCh()
			}
			if len(h.connections) == 0 {
				// mitigate closing after page reload
				go isItReallyDead(h)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				c.sendMsg(m)
			}
		}
	}
}

func isItReallyDead(h *hub) {
	time.Sleep(time.Second * 2)

	if len(h.connections) == 0 {
		h.exit <- true
	}
}
