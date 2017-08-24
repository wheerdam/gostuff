<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="{{.Prefix}}/static/style.css">
    <title>BBI Inventory System</title>
	<script>
	function deleteItem() {
		var txt;
		if (confirm("Are you sure you want to delete this item?") == true) {
			window.location = "{{.Prefix}}/delete?id={{.Info.ItemID}}";
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
			<img style="float:right" src="{{.Prefix}}/static/logo.png" height="80px">			
			<p class="centered">Welcome, {{.UserName}}! [<a href="{{.Prefix}}/logout">Logout</a>]</p>
			<hr>
			<h2 style="vertical-align: middle"><img src="{{.Prefix}}/static/open.png" height="30px" style="vertical-align: middle"> Item #{{.Info.ItemID}}</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p><button onclick="window.location.href='./list'">List All</button>
			<button onclick="window.location.href='./browse'">Browse</button>
			<button onclick="window.location.href='./edit?id={{.Info.ItemID}}'">Edit this Item</button>
			<button onclick="deleteItem()">Delete Item</button>
			</div>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>
			<form action="{{.Prefix}}/item" method="post" style="margin: 0px; padding: 0px">
				<label for="one" style="margin-left: 5px; margin-right: 5px">View Item ID:</label>
				<input id="one" size="20" type="text" name="id"  style="margin:0px">
				<input type="submit" value="View" style="margin:0px">
				<a href="{{.Prefix}}/edit">Add</a>
			</form>			
			</p>
			</div>
			<table width="100%">
			<tr><td style="white-space: nowrap">Part Number</td><td width="100%">{{.Info.Model_number}}</td></tr>
			<tr><td style="white-space: nowrap">Manufacturer</td><td><a href="{{.Prefix}}/list?manufacturer={{.Info.Manufacturer}}">{{.Info.Manufacturer}}</a></td></tr>
			<tr><td style="white-space: nowrap">Type</td><td><a href="{{.Prefix}}/list?type={{.Info.Type}}">{{.Info.Type}}</a></td></tr>
			<tr><td style="white-space: nowrap">Sub-type</td><td><a href="{{.Prefix}}/list?subtype={{.Info.Subtype}}">{{.Info.Subtype}}</a></td></tr>
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
			<th>Location</th>
			<th>Quantity</th>
			<th>Delete</th>
			</tr>
			{{range .InvEntries}}
			<tr>
				<form action="{{$.Prefix}}/modify-qty" method="post">
					<td><input type="hidden" name="id" readonly="readonly" value="{{.ItemID}}">
					<input type="hidden" name="opt" value="set">
					<input size="20" type="text" name="location" readonly="readonly" value="{{.Location}}"></td>
					<td><input size="10" type="text" name="quantity" value="{{.Quantity}}">
					<input type="submit" value="Update"></td>
				</form>
				<form action="{{$.Prefix}}/delete-entry" method="post">
					<td>
					<input type="hidden" readonly="readonly" name="id" value="{{.ItemID}}">
					<input type="hidden" readonly="readonly" name="location" value="{{.Location}}">
					<input type="submit" value="Delete Location">
					</td>
				</form>
			</tr>
			{{end}}
			</table>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p>
				<form action="{{.Prefix}}/add-entry" method="post" style="margin: 0px; padding: 0px">
					<input type="hidden" readonly="readonly" name="id" value="{{.Info.ItemID}}">
					<label for="location">Add a Location:</label><input id="location" size="20" type="text" name="location">
					<input type="submit" value="Add">
				</form>
			</p>
			</div>			
		</div>
    </body>
</html>