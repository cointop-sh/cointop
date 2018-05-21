aul:
	@echo "no default"

run:
	go run main.go

debug:
	DEBUG=1 go run main.go

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

snap/clean:
	snapcraft clean
	rm -f cointop_*.snap

snap/stage:
	snapcraft stage

snap/build: snap/clean snap/stage
	snapcraft snap

snap/deploy:
	snapcraft push cointop_*.snap --release stable

snap/remove:
	snap remove cointop

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

rpm/deps:
	sudo dnf install copr-cli
	sudo dnf install rpm-build
	sudo dnf install dnf-plugins-core

rpm/build:
	rpmbuild -ba cointop.spec

copr/create-project:
	copr-cli create cointop --chroot fedora-rawhide-x86_64

copr/build:
	rm -rf ~/rpmbuild/SRPMS/cointop-*.rpm
	copr-cli build cointop ~/rpmbuild/SRPMS/cointop-*.rpm

copr/publish:

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
