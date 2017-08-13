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
			<h2 style="vertical-align: middle"><img src="./static/list.png" height="30px" style="vertical-align: middle"> {{.ViewTitle}}</h2>			
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p style="line-height: 175%">
			<span style="color: #babaff">&#9656;</span><strong><a href="./list">LIST ALL</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong><a href="./browse">BROWSE</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong><a href="./search">SEARCH</a></strong>&nbsp;
			{{range $str := .Types}}
			<span style="color: #babaff">&#9656;</span><a href="./list?type={{$str}}">{{$str}}</a>&nbsp;&nbsp;
			{{end}}
			</p>
			</div>
			
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>
			<form action="/item" method="post" style="margin: 0px; padding: 0px">
				<label for="one" style="margin-left: 5px; margin-right: 5px">View Item ID:</label>
				<input id="one" size="20" type="text" name="id"  style="margin:0px">
				<input type="submit" value="View / Add" style="margin:0px">
			</form>
			</p>
			</div>
			
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>{{.ViewOps}}</p>
			</div>

			<table width="100%">
			<tr>
			<th></th>
			<th>Qty</th>			
			<th>Type</th>
			<th>Manufacturer</th>
			<th>Part</th>
				
			
			</tr>
			{{range .Data}}
			<tr>			
			
			<td style="padding: 5px; white-space: nowrap"><p style="font-size: 9pt">			
			{{.ItemID}}</p>
			<p style="font-size: 9pt; margin-left: 5px">			
			<a href="./item?id={{.ItemID}}"><img src="./static/open.png" height="12px" title="View Entry"></a>			
			<a href="./edit?id={{.ItemID}}"><img src="./static/edit.png" height="12px" title="Edit Entry"></a>
			<a href="{{.Seller1URL}}"><img src="./static/buy.png" height="12px" title="Seller Link"></a>
			<a href="https://www.google.com/search?q=%22{{.Manufacturer}}%22 %22{{.Model_number}}%22"><img src="./static/goog.png" height="12px" title="Search Google"></a>
			</p>
			</td>			
			
			<td style="padding: 10px; text-align: right; background-color: #ddffdd">{{.TotalQty}}</td>
			
			<td style="padding: 5px; white-space: nowrap">
			<p><a href="./list?type={{.Type}}&subtype={{.Subtype}}">{{.Subtype}}</a></p>
			<p style="font-size: 9pt; margin-left: 5px"><span style="color: #babaff">&#9656;</span>
			<a href="./list?type={{.Type}}">{{.Type}}</a></p></td>
			
			<td style="white-space: nowrap; padding: 5px; vertical-align: middle">
			<p style="margin-left: 5px"><span style="color: #babaff">&#9656;</span> <a href="./list?manufacturer={{.Manufacturer}}">{{.Manufacturer}}</a></p>			
			</td>
			
			<td style="width: 100%; padding: 5px">
			<p><a href="./item?id={{.ItemID}}">{{.Model_number}}</a></p>
			<p style="font-size: 9pt">
			{{.Descriptive_name}}</p>
			</td>

			</tr>
			{{end}}
			</table>
		</div>
		<div style="text-align: center">
		<p><a href="https://golang.org"><img src="./static/goproject.png"></a></p>
		<p style="font-size:8pt">Powered by <a href="https://golang.org">Go</a> - Gopher art by <a href="https://golang.org/doc/gopher/README?m=text">Renee French</a> <a href="https://creativecommons.org/licenses/by/3.0/">CC-BY 3.0</a></p>
		</div>
    </body>
</html>