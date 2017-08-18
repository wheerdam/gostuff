<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>Login</title>
    </head>
    <body>
		<div class="container">
			<p>You must login to continue</p>
			<hr>
			<form action="/login" method="post">
				<p><label for="one">Username:</label><input id="one" size="20" type="text" name="username"></p>
				<p><label for="two">Password:</label><input id="two" size="20" type="password" name="password"></p>
				<p><input type="submit" value="Login"></p>
			</form>
		</div>
    </body>
</html>
