FROM debian:wheezy

RUN apt-get update -y && \
	apt-get install --no-install-recommends -y -q curl && \
	echo "deb http://nginx.org/packages/debian/ wheezy nginx" >> /etc/apt/sources.list && \
	curl http://nginx.org/packages/keys/nginx_signing.key | apt-key add - && \
	apt-get update && \
	apt-get install -y -q nginx && \
	echo "\ndaemon off;" >> /etc/nginx/nginx.conf

ADD . /app
RUN chmod +x /app/start.sh

WORKDIR /etc/nginx

CMD ["bash", "/app/start.sh"]