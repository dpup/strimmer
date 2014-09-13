FROM google/debian:wheezy
MAINTAINER Dan Pupius <dan@pupi.us>

RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates

EXPOSE 80

RUN mkdir -p /usr/local/strimmer/
WORKDIR /usr/local/strimmer/

ADD out/strimmer.linux.amd64 /usr/local/strimmer/strimmer
ADD index.html /usr/local/strimmer/index.html
ADD styles.css /usr/local/strimmer/styles.css
ADD app.jsx /usr/local/strimmer/app.jsx

CMD ./strimmer --self strimmer.endoflow.com --port 80
