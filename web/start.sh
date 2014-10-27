#!/bin/bash

cat > /etc/nginx/conf.d/default.conf << EOF
server {
	listen 80;
	server_name _;

	location / {
		root /app;
		index index.html;
	}

	location /api {
		proxy_pass http://$API_HOST:$API_PORT/;
	}

}
EOF

nginx