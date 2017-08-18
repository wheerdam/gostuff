<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>{{.Header}}</title>
    </head>
    <body>
		{{.Options}}
		<hr />
		<p><a href="{{.Up}}">&#9652; Go Up</a></p>
		{{range .Dirs}}
			{{.}}
		{{end}}
		{{range .Others}}
			{{.}}
		{{end}}
		{{range .Medias}}
			{{.}}
		{{end}}		
    </body>
</html>
