package mdprev

import (
	"bytes"
	"html/template"
)

func ToHTML(md string) (html bytes.Buffer) {
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

const HTMLTemplate string = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Marked in the browser</title>
  <script type="text/javascript">
		{{.JS}}
	</script>
	<style>
	   {{.CSS}}

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
    document.getElementById('content').innerHTML = marked('{{.Markdown}}');

		var ws = new WebSocket("ws://" + window.location.host + "/ws");

		ws.onerror = function(error){
			console.log('Error detected: ' + error);
		}
		ws.onclose = function(){
			console.log('Connection closed');
		}
		ws.onmessage = function(e) {
			console.log(event.data);
			if(event.data  == 'ping'){
				//ws.send('pong');
			}else{
				document.getElementById('content').innerHTML = marked(event.data);
			}
		};
  </script>
</body>
</html>
`
