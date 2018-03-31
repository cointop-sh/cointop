all:
	@echo "no default"

run:
	go run main.go

# http://macappstore.org/upx
build: clean
	go build -ldflags "-s -w" -o bin/cointop && upx bin/cointop

clean:
	go clean && rm -f bin/cointop
