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


	</style>
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<p class="centered">Welcome, {{.UserName}}! [<a href="./logout">Logout</a>]</p>
			<h2>{{.ViewTitle}}</h2>
			<p><a href="./view">View All</a> - <a href="./new">New Item</a></p>
			<table width="100%">
			<tr>
			<th>Item ID</th>
			<th>Type</th>
			<th>Model#</th>
			<th>Description</th>
			<th>Total</th>
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