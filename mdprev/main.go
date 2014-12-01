package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"code.google.com/p/go.net/websocket"

	fsnotify "gopkg.in/fsnotify.v1"
)

var (
	mdContent string
	wsConn    *websocket.Conn
)

func main() {
	flag.Parse()

	mdFileName := flag.Arg(0)
	if _, err := os.Stat(mdFileName); os.IsNotExist(err) {
		fmt.Printf("File doesn't exist: %s", mdFileName)
		return
	}
	mdContent = loadMD(mdFileName)

	// observe file for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	watchFile(watcher, mdFileName)

	// Open it in the default browser
	go func() {
		cmd := "open http://localhost:9900"
		out, err := exec.Command("bash", "-lc", cmd).Output()
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%s", out)
	}()

	runServer()
}

func runServer() {
	var blockingCH chan string
	mdHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := toHTML(mdContent)

		fmt.Fprintf(w, html.String())
	}

	webSocketHandler := func(ws *websocket.Conn) {
		wsConn = ws
		log.Println("ws connection established.")
		// It blocks the connection from closing
		<-blockingCH
	}

	http.HandleFunc("/", mdHandler)
	http.Handle("/ws", websocket.Handler(webSocketHandler))

	log.Fatal(http.ListenAndServe(":9900", nil))
}

func toHTML(md string) (html bytes.Buffer) {
	mdJS, _ := Asset("assets/marked.min.js")
	ghCSS, _ := Asset("assets/github-markdown.css")
	page := struct {
		Markdown string
		JS       template.JS
		CSS      template.CSS
	}{md, template.JS(string(mdJS)), template.CSS(string(ghCSS))}

	t, _ := template.New("index.html").Parse(HTMLTemplate)
	t.Execute(&html, page)
	return
}

func loadMD(fileName string) string {
	body, _ := ioutil.ReadFile(fileName)
	return string(body)
}

func watchFile(watcher *fsnotify.Watcher, mdFileName string) {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				if event.Op&fsnotify.Write == fsnotify.Write {

					mdContent = loadMD(mdFileName)
					fmt.Fprintf(wsConn, mdContent)
					log.Println("modified file:", event.Name)

				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	err := watcher.Add(mdFileName)
	if err != nil {
		panic(err)
	}
}
