<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>{{.Header}}</title>
    </head>
    <body>
		<div style="width: 100%; color: #0077aa; background-color: #151515; vertical-align: middle; box-sizing: border-box; border-bottom: 1px solid #00ccff">
			<div style="float: right; padding: 5px; margin-top: 1px">
				<a href="./logout">Logout</a>
			</div>
			<div style="padding: 5px">
				<a href="{{.Up}}" style="margin-right: 15px">&#9652; Up</a>{{.Options}}
			</div>			
		</div>
		<p style="font-size: 9pt">{{.DirInfo}}</p>
		{{range .Dirs}}
			{{.}}
		{{end}}
		{{range .Others}}
			{{.}}
		{{end}}
		{{.MPre}}
		{{range .Medias}}
			{{.}}
		{{end}}	
		{{.MPost}}
		<div style="width: 100%; color: #0077aa; background-color: #151515; vertical-align: middle; box-sizing: border-box; border-top: 1px solid #00ccff">
			<div style="float: right; padding: 5px; margin-top: 1px">
				<a href="./logout">Logout</a>
			</div>
			<div style="padding: 5px">
				<a href="{{.Up}}" style="margin-right: 15px">&#9652; Up</a>{{.Options}}
			</div>			
		</div>
    </body>
</html>
