package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Stripping out {$date =>} from mongoexport
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(replace(scanner.Text())) // Println will add back the final '\n'
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func replace(str string) string {
	re, err := regexp.Compile(`{\s"\$date" : "([0-9\-:T\+\.]+)" \}`)
	if err != nil {
		panic(err)
	}

	return re.ReplaceAllString(str, "\"${1}\"")
}
