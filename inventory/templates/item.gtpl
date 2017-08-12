<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>BBI Inventory System</title>
	<script>
	function deleteItem() {
		var txt;
		if (confirm("Are you sure you want to delete this item?") == true) {
			window.location = "./delete?id={{.Info.ItemID}}";
		}
	}
	</script>
	<style>
		td, th {
			border: 1px solid #ababab;
			padding: 5px;
		}
	</style>
    </head>
    <body>
		<div class="container">
			<img style="float:right" src="./static/logo.png" height="80px">			
			<p class="centered">Welcome, {{.UserName}}! [<a href="./logout">Logout</a>]</p>
			<hr>
			<h2 style="vertical-align: middle"><img src="./static/open.png" height="30px" style="vertical-align: middle"> Item #{{.Info.ItemID}}</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p><button onclick="window.location.href='./view'">View All</button>
			<button onclick="window.location.href='./edit?id={{.Info.ItemID}}'">Edit this Item</button>
			<button onclick="deleteItem()">Delete Item</button>
			</div>
			<table width="100%">
			<tr><td style="white-space: nowrap">Part Number</td><td width="100%">{{.Info.Model_number}}</td></tr>
			<tr><td style="white-space: nowrap">Manufacturer</td><td><a href="./view?manufacturer={{.Info.Manufacturer}}">{{.Info.Manufacturer}}</a></td></tr>
			<tr><td style="white-space: nowrap">Type</td><td><a href="./view?type={{.Info.Type}}">{{.Info.Type}}</a></td></tr>
			<tr><td style="white-space: nowrap">Sub-type</td><td>{{.Info.Subtype}}</td></tr>
			<tr><td style="white-space: nowrap">Description</td><td>{{.Info.Descriptive_name}}</td></tr>
			<tr><td style="white-space: nowrap">Physical Description</td><td>{{.Info.Phys_description}}</td></tr>
			<tr><td style="white-space: nowrap">Product Link</td><td><a href="{{.Info.ProductURL}}">{{.Info.ProductURL}}</a></td></tr>
			<tr><td style="white-space: nowrap">Datasheet Link</td><td><a href="{{.Info.DatasheetURL}}">{{.Info.DatasheetURL}}</a></td></tr>
			<tr><td style="white-space: nowrap">Seller 1</td><td><a href="{{.Info.Seller1URL}}">{{.Info.Seller1URL}}</a></td></tr>
			<tr><td style="white-space: nowrap">Seller 2</td><td><a href="{{.Info.Seller2URL}}">{{.Info.Seller2URL}}</a></td></tr>
			<tr><td style="white-space: nowrap">Seller 3</td><td><a href="{{.Info.Seller3URL}}">{{.Info.Seller3URL}}</a></td></tr>
			<tr><td style="white-space: nowrap">Unit Price</td><td>{{.Info.UnitPrice}}</td></tr>
			<tr><td style="white-space: nowrap">Notes</td><td>{{.Info.Notes}}</td></tr>
			<tr><td style="white-space: nowrap">Value</td><td>{{.Info.Value}}</td></tr>
			<tr><td style="white-space: nowrap">Total Qty</td><td>{{.Info.TotalQty}}</td></tr>
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
			<td><input size="10" type="text" name="id" readonly="readonly" value="{{.ItemID}}"></td>
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