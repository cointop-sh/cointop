all:
	@echo "no default"

run:
	go run cointop.go keybindings.go navigation.go sort.go

build:
	go build
