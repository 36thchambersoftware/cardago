build:
	cd cmd/cardano; GOOS=linux GOARCH=amd64 go build -o ~/go/bin/

install:
	cd cmd/cardano; go install