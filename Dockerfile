FROM golang:onbuild
EXPOSE 80
CMD ["go-wrapper", "run", "-port", "80", "-self", "bridge-vdqd7cwzzr.elasticbeanstalk.com"]
