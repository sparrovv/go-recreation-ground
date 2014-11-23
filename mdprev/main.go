package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	md := loadMD(os.Args[1])

	output := toHTML(md)
	tfile, _ := ioutil.TempFile("", "")
	// @FIX - file is removed before the browser has chance to open it
	//defer os.Remove(tfile.Name())

	ioutil.WriteFile(tfile.Name(), output.Bytes(), os.ModeTemporary)

	// @FIX - choose the default browser
	cmd := "open -a 'Google Chrome' " + tfile.Name()
	out, err := exec.Command("bash", "-lc", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}

func toHTML(md string) (html bytes.Buffer) {
	page := struct {
		Markdown string
	}{md}

	t, _ := template.New("index.html").ParseFiles("index.html.tpl")

	t.ExecuteTemplate(&html, "index.html.tpl", page)
	return
}

func loadMD(fileName string) string {
	body, _ := ioutil.ReadFile(fileName)
	return string(body)
}
