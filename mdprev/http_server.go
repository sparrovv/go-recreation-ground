package mdprev

import (
	"fmt"
	"log"
	"net/http"
)

func (mdPrev *MdPrev) RunServer(portNumber string) {
	http.Handle("/"+mdPrev.MdFile, mdFileHandler(mdPrev))
	http.Handle("/", staticFileHandler(mdPrev))

	h := newHub(mdPrev.Broadcast, mdPrev.Exit)
	http.Handle("/ws", wsHandler(h))
	go h.run()

	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}

func mdFileHandler(mdPrev *MdPrev) http.Handler {
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
		//h.register(ws)
		c := &connection{send: make(chan []byte, 256), ws: ws}
		h.register <- c
		go c.writer()
		c.unregisterOnEOF(h)
	})
}
