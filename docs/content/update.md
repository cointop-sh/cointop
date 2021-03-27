---
title: "Update"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Update

## Go

To update make sure to use the `-u` flag if installed via Go.

```bash
go get -u github.com/miguelmota/cointop
```

## Homebrew (macOS)

```bash
brew uninstall cointop && brew install cointop
```

## Snap (Ubuntu)

Use the `refresh` command to update snap.

```bash
sudo snap refresh cointop
```

## Copr (Fedora)

```bash
sudo dnf update cointop
```

## AUR (Arch Linux)

```bash
yay -S cointop
```

## XBPS (Void Linux)

```bash
sudo xbps-install -Su cointop
```

## Flatpak (Linux)

```bash
sudo flatpak uninstall com.github.miguelmota.Cointop
sudo flatpak install flathub com.github.miguelmota.Cointop
```

## NixOS (Linux)

```bash
nix-env -uA nixpkgs.cointop
```

## AppImage (Linux)

Use the same [install](/install/#appimage-linux) instructions to update AppImage.
