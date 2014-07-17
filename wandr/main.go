package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var runningJob bool

func main() {
	flag.Parse()
	fileToWatch := flag.Arg(0)
	commandToRun := flag.Arg(1)
	validateInput(fileToWatch, commandToRun)

	log.Println("Watching:", fileToWatch, "running", commandToRun)

	executeCommand := func() {
		runningJob = true

		defer func() {
			runningJob = false
		}()

		cmd := exec.Command("sh", "-l", "-c", commandToRun)
		out, err := cmd.Output()

		if err != nil {
			println(err.Error())
			return
		}
		print(string(out))
	}

	for {
		changed, err := isFileChanged(fileToWatch)
		if err != nil {
			fmt.Println(err)
		}

		if changed && runningJob == false {
			go executeCommand()
		}
	}
}

// Naive checking if file has been changed
func isFileChanged(filePath string) (bool, error) {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return false, err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			return true, err
		}
		time.Sleep(1 * time.Second)
	}

	return false, err
}

func validateInput(arg1 string, arg2 string) {
	if arg1 == "" || arg2 == "" {
		fmt.Println("Looks like you didn't pass all required args.")
		fmt.Println("this how you use wandr:")
		fmt.Println("wandr 'file_to_watch' 'command_to_run'")
		os.Exit(1)
	}
}
