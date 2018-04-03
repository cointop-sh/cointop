# cointop

> Coin tracking for hackers

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/miguelmota/cointop/master/LICENSE.md) [![Go Report Card](https://goreportcard.com/badge/github.com/miguelmota/cointop?)](https://goreportcard.com/report/github.com/miguelmota/cointop) [![GoDoc](https://godoc.org/github.com/miguelmota/cointop?status.svg)](https://godoc.org/github.com/miguelmota/cointop)

<img src="./assets/screenshot-001.gif" width="880" />

[`cointop`](https://github.com/miguelmota/cointop) is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time. The interface is inspired by [`htop`](https://en.wikipedia.org/wiki/Htop).

## Features

- Quick sort shortcuts
- Vim style keys
- Pagination
- Color coded

#### Future releases

- Advanced search
- "Favorites" list
- Coin charts
- Currency conversion (i.e. Euro, Yen)
- Markets/Exchanges
- CryptoCompare API

## Install

Make sure to have [go](https://golang.org/) (1.9+) installed, then do:

```bash
go get -u github.com/miguelmota/cointop
```

<!--
#### Alternatively (without go)

```
sudo curl -s "https://raw.githubusercontent.com/miguelmota/cointop/master/install.sh?$(date +%s)" | bash
```
-->

## Usage

```bash
$ cointop
```

### Table commands

List of shortcuts:

|Key|Action|
|----|------|
|`<up>`|navigate up|
|`<down>`|navigate down|
|`<right>`|next page|
|`<left>`|previous page|
|`<enter>`|visit highlighted coin on CoinMarketCap|
|`<esc>`|alias to quit|
|`<space>`|alias to `<enter>`|
|`<ctrl-c>`|alias to quit|
|`<ctrl-d>`|page down|
|`<ctrl-n>`|alias to next page|
|`<ctrl-u>`|page up|
|`<ctrl-p>`|alias to previous page|
|`<ctrl-r>`|force refresh|
|`1`|sort by *[1] hour change*|
|`2`|sort by *[2]4 hour change*|
|`7`|sort by *[7] day change*|
|`a`|sort by *[a]vailable supply*|
|`j`|alias to `<down>`|
|`k`|alias to `<up>`|
|`l`|sort by *[l]ast updated*|
|`m`|sort by *[m]arket cap*|
|`n`|sort by *[n]ame*|
|`p`|sort by *[p]rice*|
|`r`|sort by *[r]ank*|
|`s`|sort by *[s]ymbol*|
|`t`|sort by *[t]otal supply*|
|`v`|sort by *24 hour [v]olume*|
|`q`|[q]uit|

<!--
|`h`|toggle [h]elp|
|`?`|alias to help|
-->

## FAQ

- Q: Where is the data from?

  - A: The data is from [Coin Market Cap](https://coinmarketcap.com/).

- Q: What coins does this support?

  - A: This supports any coin listed on [Coin Market Cap](https://coinmarketcap.com/).

- Q: How often is the data polled?

  - A: Data gets polled once every minute by default.

- Q: I installed cointop without errors but the command is not found.

  - A: Make sure your `GOPATH` and `PATH` is set correctly.
    ```bash
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    ```

- Q: What is the size of the binary?

  - A: The executable is only ~1.9MB in size.

## Authors

- [Miguel Mota](https://github.com/miguelmota)

## License

Released under the MIT license.
