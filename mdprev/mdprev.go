package mdprev

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	fsnotify "gopkg.in/fsnotify.v1"
)

type MdPrev struct {
	MdFile    string
	MdContent string
	MdChanges chan bool
	WSConns   []io.Writer
}

func NewMdPrev(mdFile string) *MdPrev {
	ch := make(chan bool)
	cons := make([]io.Writer, 0)
	mdPrev := &MdPrev{mdFile, "", ch, cons}
	mdPrev.loadContent()

	return mdPrev
}

func (m *MdPrev) loadContent() {
	body, _ := ioutil.ReadFile(m.MdFile)
	m.MdContent = string(body)
}

func (m *MdPrev) KeepWSConn(c io.Writer) {
	m.WSConns = append(m.WSConns, c)
}

// observe file for changes
func (mdPrev *MdPrev) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					mdPrev.MdChanges <- true
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
