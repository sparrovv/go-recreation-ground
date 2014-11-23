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

	t, _ := template.New("index.html").Parse(HTMLTemplate)

	t.Execute(&html, page)
	return
}

func loadMD(fileName string) string {
	body, _ := ioutil.ReadFile(fileName)
	return string(body)
}

const HTMLTemplate string = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Marked in the browser</title>
  <script src="http://cdn.rawgit.com/chjj/marked/v0.3.2/lib/marked.js"></script>
	<link rel="stylesheet" type="text/css" href="http://cdn.rawgit.com/sindresorhus/github-markdown-css/v1.2.2/github-markdown.css">
	<style>
	   #content {
			 width: 90%;
			 margin: 0 auto;
			 padding: 30px;
			 border:  1px solid #ddd;
			 border-radius: 3px;
		 }
	</style>
</head>
<body>
  <div id="content" class="markdown-body"></div>
  <script>
    document.getElementById('content').innerHTML =
      marked('{{.Markdown}}');
  </script>
</body>
</html>
`
