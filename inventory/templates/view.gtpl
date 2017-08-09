<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<img style="float:right" src="./static/logo.png" height="80px">			
			<p class="centered">Welcome, {{.UserName}}! [<a href="./logout">Logout</a>]</p>
			<hr>
			<h2>{{.ViewTitle}}</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>{{.ViewOps}}</p>
			</div>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>
			<form action="/item" method="post">
				<label for="one">View Item ID:</label><input id="one" size="20" type="text" name="id">
				<input type="submit" value="View / Add">
			</form>
			</p>
			</div>
			<div style="border: 1px solid #ababab; padding: 0px; margin-bottom: 10px;">
			<p class="mono">
			{{range $str := .Types}}
			&#9656; <a href="./view?type={{$str}}">{{$str}}</a>
			{{end}}
			</p>
			</div>
			<div style="border: 1px solid #ababab; padding: 0px; margin-bottom: 10px;">
			<p class="mono">
			{{range $str := .Manufacturers}}
			&#9656; <a href="./view?manufacturer={{$str}}">{{$str}}</a>
			{{end}}
			</p>
			</div>
			<table width="100%">
			<tr>
			<th>Type</th>
			<th>Subtype</th>
			<th>Model#</th>			
			<th>Manufacturer</th>
			<th>ID</th>
			<th>Qty</th>
			</tr>
			{{range .Data}}
			<tr>			
			<td><a href="./view?type={{.Type}}">{{.Type}}</a></td>
			<td><a href="./view?type={{.Type}}&subtype={{.Subtype}}">{{.Subtype}}</a></td>			
			<td>{{.Model_number}}</td>
			<td><a href="./view?manufacturer={{.Manufacturer}}">{{.Manufacturer}}</a></td>
			<td><a href="./item?id={{.ItemID}}"><img src="./static/open.png" height="20" alt="View Item"></a> <span style="font-size: 7pt">{{.ItemID}}</span></td>
			<td>{{.TotalQty}}</td>
			</tr>
			{{end}}
			</table>
		</div>
    </body>
</html>