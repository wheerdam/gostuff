<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body {
			/*background-color: black;
			color: white;*/
			font-family: verdana, sans-serif;
			font-size: 16px;
		}
		
		div.container {
			margin: 0 auto;
			text-align: left;
		}

		
		@media only screen and (min-width: 800px) {
			div.center {
				text-align: center;
				width: 100%;
			}

			div.container {
				width: 750px;
			}

			img.optimized {
				width: 600px;
			}

			img.optimized-500 {
				width: 500px;
			}

			body {
				font-size: 14px;
			}
		}
		
		table, th, td {
			border-collapse: collapse;
			border: 1px solid black;
			font-family: monospace;
			font-size: 14px;
			margin-bottom: 10px;
		}

		th, td {
			padding: 3px 3px 3px 3px;
		}

		th {
			background-color: #cccccc;
			text-align: center;
		}

		table.apitable td.a {
			background-color: #dddddd;
			border: 0px;
		}

		table.apitable td.b {
			padding-left: 30px;
			border: 0px;
			font-family: verdana, sans-serif;
			padding-bottom: 15px;
		}

		table.apitable {
			border: 0px;
		}

		form {
		  width: 100%;
		}
		label {
		  display: inline-block;
		  width: 150;
		}
		input {
		  display: inline-block;
		}

	</style>
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<p class="centered">Welcome, {{.UserName}}! [<a href="./logout">Logout</a>]</p>
			<h2>{{.Header}}</h2>
			<form action="/commit" method="post">
				<p><label for="id">ID#</label><input id="id" size="60" type="text" name="id" readonly="readonly" value="{{.Info.ItemID}}"></p>
				<p><label for="model">Model#:</label><input id="model" size="60" type="text" name="model" value="{{.Info.Model_number}}"></p>
				<p><label for="mfct">Manufacturer:</label><input id="mfct" size="60" type="text" name="manufacturer" value="{{.Info.Manufacturer}}"></p>
				<p><label for="type">Type:</label><input id="type" size="60" type="text" name="type" value="{{.Info.Type}}"></p>
				<p><label for="subtype">Subtype:</label><input id="subtype" size="60" type="text" name="subtype" value="{{.Info.Subtype}}"></p>
				<p><label for="description">Description:</label><input id="description" size="60" type="text" name="description" value="{{.Info.Descriptive_name}}"></p>
				<p><label for="phsydescr">Phys. Description:</label><input id="phsydescr" size="60" type="text" name="phys_description" value="{{.Info.Phys_description}}"></p>
				<p><label for="produrl">Product Link:</label><input id="produrl" size="60" type="text" name="productURL" value="{{.Info.ProductURL}}"></p>
				<p><label for="dataurl">Datasheet Link:</label><input id="dataurl" size="60" type="text" name="datasheetURL" value="{{.Info.DatasheetURL}}"></p>
				<p><label for="seller1url">Seller 1 Link:</label><input id="seller1url" size="60" type="text" name="seller1URL" value="{{.Info.Seller1URL}}"></p>
				<p><label for="seller2url">Seller 2 Link:</label><input id="seller2url" size="60" type="text" name="seller2URL" value="{{.Info.Seller2URL}}"></p>
				<p><label for="seller3url">Seller 3 Link:</label><input id="seller3url" size="60" type="text" name="seller3URL" value="{{.Info.Seller3URL}}"></p>
				<p><label for="price">Unit Price:</label><input id="price" size="60" type="text" name="unitprice" value="{{.Info.UnitPrice}}"></p>
				<p><label for="notes">Notes:</label><input id="notes" size="60" type="text" name="notes" value="{{.Info.Notes}}"></p>
				<p><input type="submit" value="Commit"> {{.Footer}}</p>
			</form>
			</table>
			<h2>Current Inventory Locations:</h2>
			<table width="100%">
			{{range .InvEntries}}
			<form action="/delete-entry" method="post">
			<tr>
			<td><input size="10" type="text" readonly="readonly" name="id" value="{{.ItemID}}"></td>
			<td><input size="30" type="text" readonly="readonly" name="location" value="{{.Location}}"></td>
			<td><input type="submit" value="Delete Location for this Item"></td>
			</tr>
			</form>
			{{end}}
			</table>
			<form action="/add-entry" method="post">
				<p><input size="10" type="text" readonly="readonly" name="id" value="{{.Info.ItemID}}"></p>
				<p><label for="location">Add a Location:</label><input id="location" size="60" type="text" name="location">
				<input type="submit" value="Add"></p>
				</p>
			</form>
		</div>
    </body>
</html>