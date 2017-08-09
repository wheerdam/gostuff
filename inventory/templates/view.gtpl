<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<img style="float:right" src="./static/logo.png">			
			<p class="centered">Welcome, {{.UserName}}! [<a href="./logout">Logout</a>]</p>
			<hr>
			<h2>{{.ViewTitle}}</h2>
			<p>{{.ViewOps}}</p>
			<p>
			<form action="/item" method="post">
				<label for="one">View Item ID:</label><input id="one" size="20" type="text" name="id">
				<input type="submit" value="View / Add">
			</form>
			</p>
			<p>Types:
			{{range $str := .Types}}
			<a href="./view?type={{$str}}">{{$str}}</a> -
			{{end}}
			</p>
			<p>Manufacturers:
			{{range $str := .Manufacturers}}
			<a href="./view?manufacturer={{$str}}">{{$str}}</a> -
			{{end}}
			</p>
			<table width="100%">
			<tr>
			<th>Item ID</th>
			<th>Type</th>
			<th>Model#</th>
			<th>Description</th>
			<th>Total Qty</th>
			</tr>
			{{range .Data}}
			<tr>
			<td><a href="./item?id={{.ItemID}}">{{.ItemID}}</a></td>
			<td><a href="./view?type={{.Type}}">{{.Type}}</a></td>
			<td>{{.Model_number}}</td>
			<td>{{.Descriptive_name}}</td>
			<td>{{.TotalQty}}</td>
			</tr>
			{{end}}
			</table>
		</div>
    </body>
</html>