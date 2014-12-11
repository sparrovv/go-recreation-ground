package mdprev

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

// inspired by: http://gary.burd.info/go-websocket-chat
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (mdPrev *MdPrev) RunServer(portNumber string) {
	var h = hub{
		broadcast:   mdPrev.Broadcast,
		exit:        mdPrev.Exit,
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}

	http.Handle("/"+mdPrev.MdFile, indexHandler(mdPrev))
	http.Handle("/", staticFileHandler(mdPrev))
	http.Handle("/ws", wsHandler(&h))
	go h.run()

	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}

func indexHandler(mdPrev *MdPrev) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := ToHTML(mdPrev.MdContent)

		fmt.Fprintf(w, html.String())
	})
}

func staticFileHandler(mdPrev *MdPrev) http.Handler {
	return http.FileServer(http.Dir(mdPrev.MdDirPath()))
}

func wsHandler(h *hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := &connection{send: make(chan []byte, 256), ws: ws}
		h.register <- c
		defer func() { h.unregister <- c }()
		go c.writer()
		c.reader(h)
	})
}

// Only listen to get EOF, so we can remove ws connection
func (c *connection) reader(h *hub) {
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
	c.ws.Close()
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
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	// send when all connections are closed after unregistering the last one
	exit chan bool
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			if len(h.connections) == 0 {
				// mitigate closing after page reload
				go isItReallyDead(h)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
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
