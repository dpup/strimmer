This is a hack in progress.  It allows a websocket client to get updates via
websockets from a PubSubHubbub subscription.

```
go run strimmer.go --self my.host.com:3100 --port 3100
go run strimmer.go --self demo.revproxy.com --port 3100
```

## Developing

You need the server to be publicly accessible so that the PuSH hub can facilitate
the subscription and send updates. The easiest way to do this is to set up a
t2.micro on EC2 and then use [SSH port forwarding](https://medium.com/dev-tricks/reverse-port-forwarding-220030f3c84a).
