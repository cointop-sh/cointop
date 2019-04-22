VERSION = $$(git describe --abbrev=0 --tags)

all: build

version:
	@echo $(VERSION)

run:
	go run main.go

deps:
	GO111MODULE=on go mod vendor

debug:
	DEBUG=1 go run main.go

build:
	@go build main.go

# http://macappstore.org/upx
build/mac: clean/mac
	env GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos/cointop && upx bin/macos/cointop

build/linux: clean/linux
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux/cointop && upx bin/linux/cointop

build/multiple: clean
	env GOARCH=amd64 go build -ldflags "-s -w" -o bin/cointop64 && upx bin/cointop64 && \
	env GOARCH=386 go build -ldflags "-s -w" -o bin/cointop32 && upx bin/cointop32

clean/mac:
	go clean && \
	rm -rf bin/mac

clean/linux:
	go clean && \
	rm -rf bin/linux

clean:
	go clean && \
	rm -rf bin/

test:
	go test ./...

cointop/test:
	go run main.go -test

cointop/version:
	go run main.go -version

cointop/clean:
	go run main.go -clean

cointop/reset:
	go run main.go -reset

snap/clean:
	snapcraft clean
	rm -f cointop_*.snap

snap/stage:
	GO111MODULE=off snapcraft stage

snap/build: snap/clean snap/stage
	snapcraft snap

snap/deploy:
	snapcraft push cointop_*.snap --release stable

snap/remove:
	snap remove cointop

snap/build-and-deploy: snap/build snap/deploy snap/clean
	@echo "done"

flatpak/build:
	flatpak-builder --force-clean build-dir com.github.miguelmota.Cointop.json

flatpak/run/test:
	flatpak-builder --run build-dir com.github.miguelmota.Cointop.json cointop

flatpak/repo:
	flatpak-builder --repo=repo --force-clean build-dir com.github.miguelmota.Cointop.json

flatpak/add:
	flatpak --user remote-add --no-gpg-verify cointop-repo repo

flatpak/remove:
	flatpak --user remote-delete cointop-repo

flatpak/install:
	flatpak --user install cointop-repo com.github.miguelmota.Cointop

flatpak/run:
	flatpak run com.github.miguelmota.Cointop

rpm/install/deps:
	sudo dnf install -y rpm-build
	sudo dnf install -y dnf-plugins-core

rpm/cp/specs:
	cp .rpm/cointop.spec ~/rpmbuild/SPECS/

rpm/build:
	rpmbuild -ba ~/rpmbuild/SPECS/cointop.spec

rpm/lint:
	rpmlint ~/rpmbuild/SPECS/cointop.spec

rpm/dirs:
	mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}

rpm/download:
	wget https://github.com/miguelmota/cointop/archive/$(VERSION).tar.gz -O ~/rpmbuild/SOURCES/$(VERSION).tar.gz

copr/install/cli:
	sudo dnf install -y copr-cli

copr/create-project:
	copr-cli create cointop --chroot fedora-rawhide-x86_64

copr/build:
	copr-cli build cointop ~/rpmbuild/SRPMS/cointop-*.rpm
	rm -rf ~/rpmbuild/SRPMS/cointop-*.rpm

copr/deploy: rpm/dirs rpm/cp/specs rpm/download rpm/build copr/build

brew/clean: brew/remove
	brew cleanup --force cointop
	brew prune

brew/remove:
	brew uninstall --force cointop

brew/build: brew/remove
	brew install --build-from-source cointop.rb

brew/audit:
	brew audit --strict cointop.rb

brew/test:
	brew test cointop.rb

brew/tap:
	brew tap cointop/cointop https://github.com/miguelmota/cointop

brew/untap:
	brew untap cointop/cointop

git/rm/large:
	java -jar bfg.jar --strip-blobs-bigger-than 200K .

git/repack:
	git reflog expire --expire=now --all
	git fsck --full --unreachable
	git repack -A -d
	git gc --aggressive --prune=now

release:
	rm -rf dist
	goreleaser
