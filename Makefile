all:
	@echo "no default"

run:
	go run main.go

# http://macappstore.org/upx
build: clean
	go build -ldflags "-s -w" -o bin/cointop32 && upx bin/cointop32 && \
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/cointop64 && upx bin/cointop64

clean:
	go clean && \
	rm -f bin/cointop64 && \
	rm -f bin/cointop32
