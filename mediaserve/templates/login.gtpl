<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>Login</title>
    </head>
    <body>
		<div class="container" style="text-align: center;" align="center">
		<div style="width: 100%; color: #0077aa; background-color: #151515; box-sizing: border-box; padding: 5px; border-bottom: 1px solid #00ccff">
			You must log in to continue
		</div>
		<form action="/login" method="post" style="margin-top:20px;">
			<p><span style="text-align: left;"><label for="one">Username:</label></span><input id="one" size="15" type="text" name="username"></p>
			<p><span style="text-align: left;"><label for="two">Password:</label></span><input id="two" size="15" type="password" name="password"></p>
			<p><input type="submit" value="Login" style="margin-top: 20px;"></p>
		</form>
		</div>
    </body>
</html>
