VERSION = $$(git describe --abbrev=0 --tags)
VERSION_DATE = $$(git log -1 --pretty='%ad' --date=format:'%Y-%m-%d' $(VERSION))
COMMIT_REV = $$(git rev-list -n 1 $(VERSION))

all: build

version:
	@echo $(VERSION)

commit_rev:
	@echo $(COMMIT_REV)

start:
	go run main.go

deps-clean:
	go clean -modcache
	rm -rf vendor

deps-download:
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

deps: deps-clean deps-download
vendor: deps

debug:
	DEBUG=1 go run main.go

.PHONY: build
build:
	go build -ldflags "-X github.com/miguelmota/cointop/cointop.version=$(VERSION)" -o bin/cointop main.go

# http://macappstore.org/upx
build-mac: clean-mac
	env GOARCH=amd64 go build -ldflags "-s -w -X github.com/miguelmota/cointop/cointop.version=$(VERSION)" -o bin/macos/cointop && upx bin/macos/cointop

build-linux: clean-linux
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X github.com/miguelmota/cointop/cointop.version=$(VERSION)" -o bin/linux/cointop && upx bin/linux/cointop

build-multiple: clean
	env GOARCH=amd64 go build -ldflags "-s -w -X github.com/miguelmota/cointop/cointop.version=$(VERSION)" -o bin/cointop64 && upx bin/cointop64 && \
	env GOARCH=386 go build -ldflags "-s -w -X github.com/miguelmota/cointop/cointop.version=$(VERSION)" -o bin/cointop32 && upx bin/cointop32

install: build
	sudo mv bin/cointop /usr/local/bin

uninstall:
	sudo rm /usr/local/bin/cointop

clean-mac:
	go clean && \
	rm -rf bin/mac

clean-linux:
	go clean && \
	rm -rf bin/linux

clean:
	go clean && \
	rm -rf bin/

.PHONY: docs
docs:
	(cd docs && hugo)

docs-server:
	(cd docs && hugo serve -p 8080)

docs-deploy: docs
	netlify deploy --prod

test:
	go test ./...

cointop-test:
	go run main.go -test

cointop-version:
	go run main.go -version

cointop-clean:
	go run main.go -clean

cointop-reset:
	go run main.go -reset

snap-clean:
	snapcraft clean
	rm -f cointop_*.snap
	rm -f cointop_*.tar.bz2

snap-stage:
	# https://github.com/elopio/go/issues/2
	mv go.mod go.mod~ ;GO111MODULE=off GOFLAGS="-ldflags=-s -ldflags=-w -ldflags=-X=github.com/miguelmota/cointop/cointop.version=$(VERSION)" snapcraft stage; mv go.mod~ go.mod

snap-install:
	sudo apt install snapd
	sudo snap install snapcraft --classic
	sudo snap install core20

snap-install-arch:
	yay -S snapd
	sudo snap install snapcraft --classic
	sudo ln -s /var/lib/snapd/snap /snap # enable classic snap support
	sudo snap install hello-world

snap-install-local:
	sudo snap install --dangerous cointop_master_amd64.snap

snap-build: snap-clean snap-stage
	snapcraft snap

snap-deploy:
	snapcraft push cointop_*.snap --release stable

snap-remove:
	snap remove cointop

snap-build-and-deploy: snap-build snap-deploy snap-clean
	@echo "done"

snap: snap-build-and-deploy

flatpak-build:
	flatpak-builder --force-clean build-dir com.github.miguelmota.Cointop.json

flatpak-run-test:
	flatpak-builder --run build-dir com.github.miguelmota.Cointop.json cointop

flatpak-repo:
	flatpak-builder --repo=repo --force-clean build-dir com.github.miguelmota.Cointop.json

flatpak-add-repo:
	flatpak --user remote-add --no-gpg-verify cointop-repo repo

flatpak-add-flathub:
	sudo flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo

flatpak-remove:
	flatpak --user remote-delete cointop-repo

flatpak-install:
	flatpak --user install cointop-repo com.github.miguelmota.Cointop

flatpak-install-local:
	flatpak-builder --force-clean --install --install-deps-from=flathub --user build-dir com.github.miguelmota.Cointop.json

flatpak-run:
	flatpak run com.github.miguelmota.Cointop

flatpak-update-version:
	xmlstarlet ed --inplace -u '/component/releases/release/@version' -v $(VERSION) .flathub/com.github.miguelmota.Cointop.appdata.xml
	xmlstarlet ed --inplace -u '/component/releases/release/@date' -v $(VERSION_DATE) .flathub/com.github.miguelmota.Cointop.appdata.xml

rpm-install-deps:
	sudo dnf install -y rpm-build
	sudo dnf install -y dnf-plugins-core

rpm-cp-specs:
	cp .rpm/cointop.spec ~/rpmbuild/SPECS/

rpm-build:
	rpmbuild --nodeps -ba ~/rpmbuild/SPECS/cointop.spec

rpm-lint:
	rpmlint ~/rpmbuild/SPECS/cointop.spec

rpm-dirs:
	mkdir -p ~/rpmbuild
	mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	chmod -R a+rwx ~/rpmbuild

rpm-download:
	wget https://github.com/miguelmota/cointop/archive/$(VERSION).tar.gz -O ~/rpmbuild/SOURCES/$(VERSION).tar.gz

copr-install-cli:
	sudo dnf install -y copr-cli

copr-deps: copr-install-cli rpm-install-deps

copr-create-project:
	copr-cli create cointop --chroot fedora-rawhide-x86_64

copr-build:
	copr-cli build cointop ~/rpmbuild/SRPMS/cointop-*.rpm
	rm -rf ~/rpmbuild/SRPMS/cointop-*.rpm

.PHONY: copr
copr: rpm-dirs rpm-cp-specs rpm-download rpm-build copr-build

brew-clean: brew-remove
	brew cleanup --force cointop
	brew prune

brew-remove:
	brew uninstall --force cointop

brew-build: brew-remove
	brew install --build-from-source cointop.rb

brew-audit:
	brew audit --strict cointop.rb

brew-test:
	brew test cointop.rb

brew-tap:
	brew tap cointop/cointop https://github.com/miguelmota/cointop

brew-untap:
	brew untap cointop/cointop

git-rm-large:
	java -jar bfg.jar --strip-blobs-bigger-than 200K .

git-repack:
	git reflog expire --expire=now --all
	git fsck --full --unreachable
	git repack -A -d
	git gc --aggressive --prune=now

release:
	rm -rf dist
	VERSION=$(VERSION) goreleaser

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t cointop/cointop .

docker-run:
	docker run -it cointop/cointop

docker-push:
	docker push cointop/cointop:latest

docker-build-and-push: docker-build docker-push

docker-run-ssh:
	docker run -p 2222:22 -v ~/.ssh/demo:/keys -v ~/.cache/cointop:/tmp/cointop_config --entrypoint cointop -it cointop/cointop server -k /keys/id_rsa

ssh-server:
	go run cmd/cointop/cointop.go server -p 2222

ssh-client:
	ssh localhost -p 2222

mp3:
	cat <(printf "package notifier\nfunc Mp3() string {\nreturn \`" "") <(xxd -p media/notification.mp3 | tr -d "\n") <(printf "\`\n}" "") > pkg/notifier/mp3.go

pkg2appimage-install:
	wget -c https://github.com/$(wget -q https://github.com/AppImage/pkg2appimage/releases -O - | grep "pkg2appimage-.*-x86_64.AppImage" | head -n 1 | cut -d '"' -f 2)
	chmod +x pkg2appimage-*.AppImage

appimage-clean-workspace:
	rm -rf .appimage_workspace

appimage-clean: appimage-clean-workspace
	rm -rf *.AppImage

.PHONY: appimage
appimage: appimage-clean-workspace
	( \
		mkdir -p .appimage_workspace && \
		mkdir -p dist/appimage && \
		cd .appimage_workspace && \
		../pkg2appimage-*.AppImage ../.appimage/cointop.yml && \
		cp out/cointop-*.AppImage ../dist/appimage/ \
	)

appimage-run:
	./dist/appimage/cointop-*.AppImage
