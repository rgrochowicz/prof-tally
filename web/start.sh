#!/bin/bash

cat > /etc/nginx/conf.d/default.conf << EOF
server {
	listen 80;
	server_name _;

	root /app;

	location /api {

		gzip             on;
		gzip_proxied     any;
		gzip_types       text/css text/plain text/xml application/xml application/javascript application/x-javascript text/javascript application/json text/x-json;
		proxy_pass http://$API_HOST:$API_PORT;
	}

}
EOF

nginx