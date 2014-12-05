package mdprev

import (
	"fmt"
	"io/ioutil"
	"log"

	"code.google.com/p/go.net/websocket"

	fsnotify "gopkg.in/fsnotify.v1"
)

type MdPrev struct {
	MdFile    string
	MdContent string
	MdChanges chan bool
	WSConns   []*websocket.Conn
}

func NewMdPrev(mdFile string) *MdPrev {
	ch := make(chan bool)
	cons := make([]*websocket.Conn, 0)
	mdPrev := &MdPrev{mdFile, "", ch, cons}
	mdPrev.loadContent()

	return mdPrev
}

func (m *MdPrev) loadContent() {
	body, _ := ioutil.ReadFile(m.MdFile)
	m.MdContent = string(body)
}

func (m *MdPrev) KeepWSConn(c *websocket.Conn) {
	m.WSConns = append(m.WSConns, c)
}

func (mdPrev *MdPrev) Watch() {
	// observe file for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	//defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					mdPrev.MdChanges <- true

					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(mdPrev.MdFile)
	if err != nil {
		panic(err)
	}
}

func (mdPrev *MdPrev) UpdateWSConnections() {
	go func() {
		for _ = range mdPrev.MdChanges {
			mdPrev.loadContent()

			for _, ws := range mdPrev.WSConns {
				fmt.Fprint(ws, mdPrev.MdContent)
			}
		}
	}()
}
