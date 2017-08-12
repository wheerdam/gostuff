<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<img style="float:right" src="./static/logo.png" height="80px">			
			<p>Welcome, {{.UserName}}! [<a href="./logout">Logout</a>] - CSV: <a href="./download-items">Items</a> - <a href="./download-inventory">Inventory</a></p>
			<hr>
			<h2>Browse</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p style="line-height: 175%">
			<span style="color: #babaff">&#9656;</span><strong><a href="./view">LIST ALL</a></strong>&nbsp;
			</p>
			</div>
			<div align="center">
			<table style="border: 0px">
			<tr>
			<th>Manufacturers</th>
			<th>Types, Subtypes</th>
			<th>Types, Manufacturers</th>
			</tr>
			<tr>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ range .Manufacturers }}
				<p><span style="color: #babaff">&#9656;</span>
					<a href="./view?manufacturer={{.}}">{{.}}</a>
				</p>
			{{end}}
			</td>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ $types := .Types }}
			{{ range $type := $types }}
				<div style="margin-bottom: 20px">
				<p><span style="color: #babaff">&#9656;</span>
					<a href="./view?type={{$type.Name}}">{{.Name}}</a>
				</p>
				{{range $type.Subtypes}}
					<p style="margin-left: 30px">
						<span style="color: #babaff">&#9656;</span>
						<a href="./view?type={{$type.Name}}&subtype={{.}}">{{.}}</a>
				{{end}}
				</div>
			{{end}}
			</td>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ $types := .Types }}
			{{ range $type := $types }}
				<div style="margin-bottom: 20px">
				<p><span style="color: #babaff">&#9656;</span>
					<a href="./view?type={{$type.Name}}">{{.Name}}</a>
				</p>
				{{range $type.Manufacturers}}
					<p style="margin-left: 30px">
						<span style="color: #babaff">&#9656;</span>
						<a href="./view?type={{$type.Name}}&manufacturer={{.}}">{{.}}</a>
				{{end}}
				</div>
			{{end}}
			</td>			
			</tr>
			</table>
			</div>
		</div>
		<div style="text-align: center; clear: both">
		<p><a href="https://golang.org"><img src="./static/goproject.png"></a></p>
		<p style="font-size:8pt">Powered by <a href="https://golang.org">Go</a> - Gopher art by <a href="https://golang.org/doc/gopher/README?m=text">Renee French</a> <a href="https://creativecommons.org/licenses/by/3.0/">CC-BY 3.0</a></p>
		</div>
    </body>
</html>