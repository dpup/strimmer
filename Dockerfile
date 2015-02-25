# Copied from golang:onbuild to optimize customization of apt-get steps.
# https://github.com/docker-library/golang/blob/master/1.3/onbuild/Dockerfile
FROM golang:1.3.3

# Install dependencies necessary to compile JSX
RUN apt-get update && apt-get -y upgrade && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install nodejs npm
RUN npm config set registry http://registry.npmjs.org
RUN npm install -g react-tools
RUN ln -s /usr/bin/nodejs /usr/bin/node # Fix stupid naming difference

# START golang:onbuild
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
CMD ["go-wrapper", "run"]
ONBUILD COPY . /go/src/app
ONBUILD RUN go-wrapper download
ONBUILD RUN go-wrapper install
# END golang:onbuild

# Do the JSX compilation and make compiled output available in the same dir.
COPY web/ /go/src/app/web/
RUN jsx --extension jsx web/ web/

# Update index.html to point at the compiled JS file.
COPY prod-scripts.js /go/src/app/
RUN node prod-scripts.js

EXPOSE 80

CMD ["go-wrapper", "run", "-port", "80", "-self", "bridge-vdqd7cwzzr.elasticbeanstalk.com"]
