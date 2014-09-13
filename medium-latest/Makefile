
deploy: out/strimmer.linux.amd64

out:
	mkdir -p ./out

out/strimmer.linux.amd64: *.go out
	GOOS=linux GOARCH=amd64 go build -o out/strimmer.linux.amd64 ./medium-latest.go
