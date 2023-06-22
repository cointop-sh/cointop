---
title: "Development"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Development

## Go

Running cointop from source

```bash
make run
```

## Update vendor dependencies

```bash
make deps
```

## Homebrew

Installing from source

```bash
make brew-build
```

## Flatpak

Install the freedesktop runtime (if not done so already)

```bash
sudo flatpak install flathub org.freedesktop.Platform//1.6 org.freedesktop.Sdk//1.6
```

Install golang extension

```bash
sudo flatpak install flathub org.freedesktop.Sdk.Extension.golang
```

Building flatpak package

```bash
make flatpak-build
```

## Copr

Install dependencies

```bash
make copr-install-cli
make rpm-install-deps
make rpm-dirs
```

Build package

```bash
make rpm-cp-specs
make rpm-download
make rpm-build
make copr-build
```

## Snap

Building snap

```bash
make snap-build
```

## Docker

Build Docker image

```bash
make docker-build
```
