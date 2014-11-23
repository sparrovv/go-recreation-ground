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
