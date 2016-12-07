all: build

build: ebs-snap remote-snap clean-snaps

ebs-snap:
	go build -o bin/ebs-snap ./cmd/ebs-snap

remote-snap:
	go build -o bin/remote-snap ./cmd/remote-snap

clean-snaps:
	go build -o bin/clean-snaps ./cmd/clean-snaps

install:
	go install ./cmd/...

test:
	go test ./cmd/...

