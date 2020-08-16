package main

const form = `
<html>
<head>
<!-- Latest compiled and minified CSS -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<!-- Optional theme -->
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<!-- Latest compiled and minified JavaScript -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
</head>
	<body style="background-color:black;color:white;">

		<main class="col-md-6 col-md-offset-3">
			<form action="/submit" method="POST">

				<div class="form-group" style="text-align:center;">
					<BR><BR><H1>Github repo cloner</H1><BR><BR>
				</div>
				<div class="form-group">
					<label for="giturl">Enter github url to download:</label>
					<input type="text" name="giturl" id="giturl" value="%v" class="form-control">
				</div>
				<div class="form-group">
					<label for="directory">Enter directory in which to be cloned</label>
					<input type="text" name="directory" id="directory" value="%v" class="form-control">
				</div>
				<div class="form-group">	
					<input type="submit" formaction="/clone" value="Clone" name="value" style="color: black; background-color: white; font-weight: bold;">					
					<input type="submit"  formaction="/clear" value="Clear" name="value" style="color: black; background-color: white; font-weight: bold;">
				</div>
			</form>
			<div>
				<H1><BR><BR>
					Result: %v %v %v %v %v
				</H1>
			</div>

		</main>
	</body>
</html>
`
