FROM debian:wheezy

ENV NODE_VERSION 0.11.14

RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl python build-essential git ca-certificates libpq-dev
RUN mkdir /nodejs && curl http://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz | tar xvzf - -C /nodejs --strip-components=1
ENV PATH $PATH:/nodejs/bin

WORKDIR /app
ADD . /app

RUN npm install
CMD ["/nodejs/bin/npm", "start"]