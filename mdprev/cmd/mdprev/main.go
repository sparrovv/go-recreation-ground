package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/sparrovv/go-recreation-ground/mdprev"
)

func main() {
	flag.Parse()
	mdFileName := flag.Arg(0)
	if _, err := os.Stat(mdFileName); os.IsNotExist(err) {
		fmt.Printf("File doesn't exist: %s", mdFileName)
		return
	}

	mdPrev := mdprev.NewMdPrev(mdFileName)
	mdPrev.Watch()

	// Open it in the default browser
	open.Run("http://localhost:9900")

	mdPrev.RunServer()
}
