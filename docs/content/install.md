---
title: "Install"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Install

There are multiple ways you can install cointop depending on the platform you're on.

## From source (always latest and recommended)

Make sure to have [go](https://golang.org/) (1.17+) installed, then do:

```bash
go install github.com/cointop-sh/cointop@latest
```

The cointop executable will be under your GOPATH so make sure `$GOPATH/bin` is added to the `$PATH` variable if not already.

Now you can run cointop:

```bash
cointop
```

## Binary (all platforms)

You can download the binary from the [releases](https://github.com/cointop-sh/cointop/releases) page.

```bash
curl -o- https://raw.githubusercontent.com/cointop-sh/cointop/master/install.sh | bash
```

```bash
wget -qO- https://raw.githubusercontent.com/cointop-sh/cointop/master/install.sh | bash
```

## Homebrew (macOS)

cointop is available via [Homebrew](https://formulae.brew.sh/formula/cointop) for macOS:

```bash
brew install cointop
```

Run

```bash
cointop
```

## Snap (Ubuntu)

cointop is available as a [snap](https://snapcraft.io/cointop) for Linux users.

```bash
sudo snap install cointop --stable
```

Running snap:

```bash
sudo snap run cointop
```

Note: snaps don't work in Windows WSL. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216).

## Copr (Fedora)

cointop is available as a [copr](https://copr.fedorainfracloud.org/coprs/miguelmota/cointop/) package.

First, enable the repository

```bash
sudo dnf copr enable miguelmota/cointop -y
```

Install cointop

```bash
sudo dnf install cointop
```

Run

```bash
cointop
```

## AUR (Arch Linux)

cointop is available as an [AUR](https://aur.archlinux.org/packages/cointop) package.

```bash
git clone https://aur.archlinux.org/cointop.git
cd cointop
makepkg -si
```

Using [yay](https://github.com/Jguer/yay)

```bash
yay -S cointop
```

## XBPS (Void Linux)

cointop is available as a [XBPS](https://voidlinux.org/packages/) package.

```bash
sudo xbps-install -S cointop
```

## Flatpak (Linux)

cointop is available as a [Flatpak](https://flatpak.org/) package via the [Flathub](https://flathub.org/apps/details/com.github.miguelmota.Cointop) registry.

Add the flathub repository (if not done so already)

```bash
sudo flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo
```

Install cointop flatpak

```bash
sudo flatpak install flathub com.github.miguelmota.Cointop
```

Run cointop flatpak

```bash
flatpak run com.github.miguelmota.Cointop
```

## NixOS (Linux)

cointop is available as a [nixpkg](https://search.nixos.org/packages?channel=unstable&show=cointop&from=0&size=30&sort=relevance&query=cointop).

```bash
nix-env -iA nixpkgs.cointop
```

## AppImage (Linux)

You can download the AppImage from the [releases](https://github.com/cointop-sh/cointop/releases) page.

```bash
VERSION=$(curl --silent "https://api.github.com/repos/cointop-sh/cointop/releases/latest" | grep -Po --color=never '"tag_name": ".\K.*?(?=")')
URL="https://github.com/cointop-sh/cointop/releases/download/v$VERSION/cointop-v$VERSION.glibc2.32-x86_64.AppImage"
wget $URL
```

Make AppImage executable:

```bash
chmod +x cointop-*.AppImage
```

Run AppImage:

```bash
./cointop-*.AppImage
```

## FreshPorts (FreeBSD / OpenBSD)

cointop is available as a [FreshPort](https://www.freshports.org/finance/cointop/) package.

```bash
sudo pkg install cointop
```

## Windows (PowerShell / WSL)

Install [Go](https://golang.org/doc/install) and [git](https://git-scm.com/download/win), then:

```powershell
go get -u github.com/cointop-sh/cointop
```

You'll need additional font support for Windows. Please see the [wiki](https://github.com/cointop-sh/cointop/wiki/Windows-Command-Prompt-and-WSL-Font-Support) for instructions.

## Docker

cointop is available on [Docker Hub](https://hub.docker.com/r/cointop/cointop).

```bash
docker run -it cointop/cointop
```

Note: the config is under `/root/.config/cointop` in container, so attach a volume to make it persistent in host:

```bash
docker run -v ~/.cache/cointop:/root/.config/cointop -it cointop/cointop
```

## Binaries

You can find pre-built binaries on the [releases](https://github.com/cointop-sh/cointop/releases) page.
