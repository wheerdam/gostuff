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
			<h2 style="vertical-align: middle"><img src="./static/search.png" height="30px" style="vertical-align: middle"> Search</h2>
			<div style="border: 1px solid #ababab; padding-left: 10px; margin-bottom: 10px;">
			<p style="line-height: 175%">
			<span style="color: #babaff">&#9656;</span><strong><a href="./list">LIST ALL</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong><a href="./browse">BROWSE</a></strong>&nbsp;
			<span style="color: #babaff">&#9656;</span><strong>SEARCH</strong>&nbsp;
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
			<form action="/search" method="post" style="margin: 0px; padding: 0px">
				<p><strong>Search for</strong></p>
				<p><label for="a1">Type</label><input id="a1" size="20" type="text" name="type">
				<p><label for="a2">Subtype</label><input id="a2" size="20" type="text" name="subtype">
				<p><label for="a3">Manufacturer</label><input id="a3" size="20" type="text" name="manufacturer">
				<p><label for="a4">Description</label><input id="a4" size="20" type="text" name="description">
				<p><label for="a5">Part Number</label><input id="a5" size="20" type="text" name="part_number">				
				<p><label for="a6">Physical Description</label><input id="a6" size="20" type="text" name="phys_description">	
				<p><input type="submit" value="Go" style="margin:0px"></p>
			</form>
			</div>
		</div>
		<div style="text-align: center; clear: both">
		<p><a href="https://golang.org"><img src="./static/goproject.png"></a></p>
		<p style="font-size:8pt">Powered by <a href="https://golang.org">Go</a> - Gopher art by <a href="https://golang.org/doc/gopher/README?m=text">Renee French</a> <a href="https://creativecommons.org/licenses/by/3.0/">CC-BY 3.0</a></p>
		</div>
    </body>
</html>