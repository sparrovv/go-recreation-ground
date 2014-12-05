package mdprev

import (
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

func (mdPrev *MdPrev) RunServer() {
	mdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := ToHTML(mdPrev.MdContent)

		fmt.Fprintf(w, html.String())
	}

	var blockChan chan bool
	webSocketHandler := func(ws *websocket.Conn) {
		mdPrev.KeepWSConn(ws)
		<-blockChan
	}

	http.HandleFunc("/", mdHandler)
	http.Handle("/ws", websocket.Handler(webSocketHandler))
	// waits for file changes and update ws connections
	mdPrev.UpdateWSConnections()

	log.Fatal(http.ListenAndServe(":9900", nil))
}
