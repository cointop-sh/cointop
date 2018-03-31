all:
	@echo "no default"

run:
	go run cointop.go keybindings.go navigation.go sort.go layout.go status.go chart.go table.go

build: clean
	go build -ldflags "-s -w" -o bin/cointop && upx bin/cointop

clean:
	go clean && rm -f bin/cointop
