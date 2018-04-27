all:
	@echo "no default"

run:
	go run main.go

debug:
	DEBUG=1 go run main.go

# http://macappstore.org/upx
build/mac: clean
	env GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos/cointop && upx bin/macos/cointop

build/multiple: clean
	env GOARCH=amd64 go build -ldflags "-s -w" -o bin/cointop64 && upx bin/cointop64 && \
	env GOARCH=386 go build -ldflags "-s -w" -o bin/cointop32 && upx bin/cointop32

clean:
	go clean && \
	rm -rf bin/

test:
	go test ./...

snap/clean:
	snapcraft clean

snap/stage:
	snapcraft stage

snap/build: snap/clean snap/stage
	snapcraft snap

snap/deploy:
	snapcraft push cointop_*.snap --release stable

snap/remove:
	snap remove cointop

brew/clean: brew/remove
	brew cleanup --force cointop
	brew prune

brew/remove:
	brew uninstall --force cointop

brew/build: brew/remove
	brew install --build-from-source cointop.rb

brew/audit:
	brew audit --strict cointop.rb
