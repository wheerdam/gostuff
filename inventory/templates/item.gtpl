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
			<h2>Item #{{.Info.ItemID}}</h2>
			<p><a href="./view">View All</a> - <a href="./edit?id={{.Info.ItemID}}">Edit this Item</a>
			- <a href="./delete?id={{.Info.ItemID}}">Delete this Item</a></p>
			<table width="100%">
			<tr><td>Model#</td><td>{{.Info.Model_number}}</td></tr>
			<tr><td>Manufacturer</td><td><a href="./view?manufacturer={{.Info.Manufacturer}}">{{.Info.Manufacturer}}</a></td></tr>
			<tr><td>Type</td><td><a href="./view?type={{.Info.Type}}">{{.Info.Type}}</a></td></tr>
			<tr><td>Sub-type</td><td>{{.Info.Subtype}}</td></tr>
			<tr><td>Description</td><td>{{.Info.Descriptive_name}}</td></tr>
			<tr><td>Physical Description</td><td>{{.Info.Phys_description}}</td></tr>
			<tr><td>Product Link</td><td><a href="{{.Info.ProductURL}}">{{.Info.ProductURL}}</a></td></tr>
			<tr><td>Datasheet Link</td><td><a href="{{.Info.DatasheetURL}}">{{.Info.DatasheetURL}}</a></td></tr>
			<tr><td>Seller 1</td><td><a href="{{.Info.Seller1URL}}">{{.Info.Seller1URL}}</a></td></tr>
			<tr><td>Seller 2</td><td><a href="{{.Info.Seller2URL}}">{{.Info.Seller2URL}}</a></td></tr>
			<tr><td>Seller 3</td><td><a href="{{.Info.Seller3URL}}">{{.Info.Seller3URL}}</a></td></tr>
			<tr><td>Unit Price</td><td>{{.Info.UnitPrice}}</td></tr>
			<tr><td>Notes</td><td>{{.Info.Notes}}</td></tr>
			<tr><td>Total Qty</td><td>{{.Info.TotalQty}}</td></tr>
			</table>
			<h2>Inventory Details:</h2>
			<table width="100%">
			<tr>
			<th>Item ID</th>
			<th>Location</th>
			<th>Quantity</th>
			</tr>
			{{range .InvEntries}}
			<tr>
			<form action="/modify-qty" method="post">
			<td><input size="10" tyle="text" name="id" readonly="readonly" value="{{.ItemID}}"></td>
			<td><input size="30" type="text" name="location" readonly="readonly" value="{{.Location}}"></td>
			<td><input size="10" type="text" name="quantity" value="{{.Quantity}}">
			<input type="submit" value="Update"></td>
			</form>
			</tr>
			{{end}}
			</table>
		</div>
    </body>
</html>