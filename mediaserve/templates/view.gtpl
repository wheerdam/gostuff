<html>
    <head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" type="text/css" href="/static/style.css">
    <title>{{.Header}}</title>
	<style>
	form {
		margin: 0;
		padding: 0;
	}
	</style>
    </head>
    <body>
		<div style="width: 100%; color: #0077aa; background-color: #151515; vertical-align: middle; box-sizing: border-box; border-bottom: 1px solid #00ccff">
			<div style="float: right; padding: 5px; margin-top: 1px">
				<a href="./logout">Logout</a>
			</div>
			<div style="padding: 5px">
				<a href="{{.Up}}" style="margin-right: 15px">&#9652; Up</a>{{.Options}}
			</div>			
		</div>
		<p style="font-size: 9pt">{{.DirInfo}}</p>
		{{range .Dirs}}
			{{.}}
		{{end}}
		{{range .Others}}
			{{.}}
		{{end}}
		{{.MPre}}
		{{range .Medias}}
			{{.}}
		{{end}}	
		{{.MPost}}
		<div style="width: 100%; color: #0077aa; background-color: #151515; vertical-align: middle; box-sizing: border-box; border-top: 1px solid #00ccff">
			<div style="float: right; padding: 5px; margin-top: 1px">
				<a href="./logout">Logout</a>
			</div>
			<div style="padding: 5px">
				<a href="{{.Up}}" style="margin-right: 15px">&#9652; Up</a>
				<button onclick="toggleUpload()">Upload</button>
			</div>			
			<div id="upload" style="display: none; padding: 8px; box-sizing: border-box; width: 100%; border: 1px solid red">
				<form enctype="multipart/form-data" action="/upload" method="post">
					<input type="hidden" name="path" value="{{.Path}}">
					<label for="upload">Upload</label><br><input type="file" name="upload" style="width: 250px">
					<input type="submit" value="Upload" style="width: 80px">
				</form> 
				<form enctype="multipart/form-data" action="/get" method="post">
					<input type="hidden" name="path" value="{{.Path}}">
					<label for="url" style="margin-top: 8px">Web URL</label><br><input type="text" name="url" style="width: 250px;">
					<input type="submit" value="Get" style="width: 80px">
				</form>
			</div>
		</div>
    </body>
	<script>
		function toggleUpload() {
			var x = document.getElementById("upload");
			if(x.style.display === "none") {
				x.style.display = "block";
				window.scrollTo(0,document.body.scrollHeight);
			} else {
				x.style.display = "none";
			}
		}
	</script>
</html>
