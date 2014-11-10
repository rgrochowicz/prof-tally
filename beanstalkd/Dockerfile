FROM debian:wheezy

RUN apt-get update -y && apt-get install --no-install-recommends -y -q beanstalkd

CMD ["/usr/bin/beanstalkd", "-z", "12000000"]