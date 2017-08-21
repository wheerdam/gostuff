<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="{{.Prefix}}/static/style.css">
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<img style="float:right" src="{{.Prefix}}/static/logo.png" height="80px">			
			<p>Welcome, {{.UserName}}! [<a href="{{.Prefix}}/logout">Logout</a>] - CSV: <a href="{{.Prefix}}/download-items">Items</a> - <a href="{{.Prefix}}/download-inventory">Inventory</a></p>
			<hr>
			<h2 style="vertical-align: middle"><img src="{{.Prefix}}/static/browse.png" height="30px" style="vertical-align: middle"> Browse</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p style="line-height: 175%">
			<span style="color: #babaff">&#9656;</span><strong><a href="{{.Prefix}}/list">LIST ALL</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong>BROWSE</strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong><a href="{{.Prefix}}/search">SEARCH</a></strong>&nbsp;
			</p>
			</div>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>
			<form action="{{.Prefix}}/item" method="post" style="margin: 0px; padding: 0px">
				<label for="one" style="margin-left: 5px; margin-right: 5px">View Item ID:</label>
				<input id="one" size="20" type="text" name="id"  style="margin:0px">
				<input type="submit" value="View / Add" style="margin:0px">
			</form>
			</p>
			</div>
			<table style="border: 0px" width="100%">
			<tr>
			<th>Manufacturers</th>
			<th>Types</th>
			<th>Types, Subtypes</th>
			<th>Types, Manufacturers</th>
			</tr>
			<tr>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ range .Manufacturers }}
				<p><span style="color: #babaff">&#9656;</span>
					<a href="{{$.Prefix}}/list?manufacturer={{.}}">{{.}}</a>
				</p>
			{{end}}
			</td>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ $types := .Types }}
			{{ range $type := $types }}
				<p><span style="color: #babaff">&#9656;</span>
					<a href="{{$.Prefix}}/list?type={{$type.Name}}">{{.Name}}</a>
				</p>
			{{end}}
			</td>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ $types := .Types }}
			{{ range $type := $types }}
				<div style="margin-bottom: 20px">
				<p style="font-size: 8pt">
					{{.Name}}
				</p>
				{{range $type.Subtypes}}
					<p style="margin-left: 10px">
						<span style="color: #babaff">&#9656;</span>
						<a href="{{$.Prefix}}/list?type={{$type.Name}}&subtype={{.}}">{{.}}</a>
				{{end}}
				</div>
			{{end}}
			</td>
			<td style="padding-left: 10px; padding-right: 10px; vertical-align: top">
			{{ $types := .Types }}
			{{ range $type := $types }}
				<div style="margin-bottom: 20px">
				<p style="font-size: 8pt">
					{{.Name}}
				</p>
				{{range $type.Manufacturers}}
					<p style="margin-left: 10px">
						<span style="color: #babaff">&#9656;</span>
						<a href="{{$.Prefix}}/list?type={{$type.Name}}&manufacturer={{.}}">{{.}}</a>
				{{end}}
				</div>
			{{end}}
			</td>			
			</tr>
			</table>
		</div>
		<div style="text-align: center; clear: both">
		<p><a href="https://golang.org"><img src="{{.Prefix}}/static/goproject.png"></a></p>
		<p style="font-size:8pt">Powered by <a href="https://golang.org">Go</a> - Gopher art by <a href="https://golang.org/doc/gopher/README?m=text">Renee French</a> <a href="https://creativecommons.org/licenses/by/3.0/">CC-BY 3.0</a></p>
		</div>
    </body>
</html>