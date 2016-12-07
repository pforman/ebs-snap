all: build

build: ebs-snap binremote-snap clean-snaps

ebs-snap:
	go build -o bin/ebs-snap ./cmd/ebs-snap

remote-snap:
	go build -o bin/remote-snap ./cmd/remote-snap

clean-snaps:
	go build -o bin/clean-snaps ./cmd/clean-snaps

install:
	go install ./cmd/...

