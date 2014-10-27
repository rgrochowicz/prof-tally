Web
===

Provides the web ui for the app.  This is the **static** portion.  It's pretty much useless without the API running.

Attempts to follow [this style guide](https://github.com/gocardless/angularjs-style-guide).

nginx configuration:
--------------------
Here's an example nginx config for this app:
```
server {
	#listen on port 80
	listen 80;

	#process for all hosts
	server_name _;

	#place where the static files exist
	root /<path to repo>/prof-tally/web;

	#send everything that's going to /api to the node process
	location /api {
		proxy_pass http://127.0.0.1:3000;
	}
}
```