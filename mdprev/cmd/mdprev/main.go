package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/sparrovv/go-recreation-ground/mdprev"
)

func main() {
	port := flag.String("port", "9900", "port number on which server listens")
	flag.Parse()
	mdFileName := flag.Arg(0)

	if _, err := os.Stat(mdFileName); os.IsNotExist(err) {
		fmt.Printf("File doesn't exist: %s", mdFileName)
		return
	}

	mdPrev := mdprev.NewMdPrev(mdFileName)
	mdPrev.Watch()

	go mdPrev.RunServer(*port)
	go mdPrev.ListenAndBroadcastChanges()

	url := "http://localhost:" + *port + "/" + mdPrev.MdFile
	open.Run(url) // Opens in the default browser
	fmt.Println("Server listens on:", url)

	<-mdPrev.Exit
}
