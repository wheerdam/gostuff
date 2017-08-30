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
			<h2 style="vertical-align: middle"><img src="{{.Prefix}}/static/scan.png" height="30px" style="vertical-align: middle"> Scan</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p style="line-height: 175%">
			<span style="color: #babaff">&#9656;</span><strong><a href="{{.Prefix}}/list">LIST ALL</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong><a href="{{.Prefix}}/browse">BROWSE</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong>SEARCH</strong>&nbsp;
			</p>
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
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
				<p>{{.Message}}</p>
			</div>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<form action="{{.Prefix}}/modify-qty" method="post" style="margin: 0px; padding: 0px">
				<p><strong>Update Quantity</strong></p>
				<p><label for="a1">Item ID</label><input autofocus="autofocus" id="a1" size="20" type="text" name="id" value="{{.PrevID}}"></p>
				<p><label for="a2">Location</label><input id="a2" size="20" type="text" name="location" value="{{.PrevLocation}}"></p>
				<p><label for="a3">Quantity</label><input id="a3" size="20" type="text" name="quantity" value="{{.PrevQty}}"></p>
				<p><label for="a4">Operation</label><input id="a4" size="20" type="text" name="op" value="{{.PrevOperation}}"></p>
				<p><input id="a5" type="checkbox" name="keepid" value="yes" {{.IDChecked}}><label for="a5">Keep ID</label></p>
				<input type="hidden" name="done" value="{{.Prefix}}/scan">
				<input type="hidden" name="prev" value="yes">
				<input type="hidden" name="getresults" value="yes">
				<p><input type="submit" value="Go" style="margin:0px"></p>
			</form>
			</div>
		</div>
		<div style="text-align: center; clear: both">
		<p><a href="https://golang.org"><img src="{{.Prefix}}/static/goproject.png"></a></p>
		<p style="font-size:8pt">Powered by <a href="https://golang.org">Go</a> - Gopher art by <a href="https://golang.org/doc/gopher/README?m=text">Renee French</a> <a href="https://creativecommons.org/licenses/by/3.0/">CC-BY 3.0</a></p>
		</div>
    </body>
</html>