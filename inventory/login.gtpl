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
		form {
		  width: 500px;
		}
		label {
		  display: inline-block;
		  width: 100px;
		}
		input {
		  display: inline-block;
		}

	</style>
    <title>BBI Inventory System</title>
    </head>
    <body>
		<div class="container">
			<form action="/login" method="post">
				<p><label for="one">Username:</label><input id="one" size="20" type="text" name="username"></p>
				<p><label for="two">Password:</label><input id="two" size="20" type="password" name="password"></p>
				<p><input type="submit" value="Login"></p>
			</form>
		</div>
    </body>
</html>