package main

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

		var ws = new WebSocket("ws://localhost:9900/ws");
		ws.onmessage = function(e) {
			console.log(event.data);
			document.getElementById('content').innerHTML = marked(event.data);
		};
  </script>
</body>
</html>
`
