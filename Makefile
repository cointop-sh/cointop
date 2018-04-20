all:
	@echo "no default"

run:
	go run main.go

# http://macappstore.org/upx
build: clean
	env GOARCH=amd64 go build -ldflags "-s -w" -o bin/cointop64 && upx bin/cointop64 && \
	env GOARCH=386 go build -ldflags "-s -w" -o bin/cointop32 && upx bin/cointop32

clean:
	go clean && \
	rm -f bin/cointop64 && \
	rm -f bin/cointop32

test:
	go test ./...

snap:
	snapcraft clean && snapcraft stage && snapcraft snap

snap/deploy:
	snapcraft push <*.snap> --release stable
