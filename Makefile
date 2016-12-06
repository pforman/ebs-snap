build:
	go build -o bin/ebs-snap ./cmd/ebs-snap
	go build -o bin/remote-snap ./cmd/remote-snap

install:
	go install ./cmd/...
