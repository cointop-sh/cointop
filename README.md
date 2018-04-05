# cointop

> Coin tracking for hackers

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/miguelmota/cointop/master/LICENSE.md) [![Build Status](https://travis-ci.org/miguelmota/cointop.svg?branch=master)](https://travis-ci.org/miguelmota/cointop) [![Go Report Card](https://goreportcard.com/badge/github.com/miguelmota/cointop?)](https://goreportcard.com/report/github.com/miguelmota/cointop) [![GoDoc](https://godoc.org/github.com/miguelmota/cointop?status.svg)](https://godoc.org/github.com/miguelmota/cointop)

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
- Custom shortcuts

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

### Cointop commands

List of default shortcuts:

Key|Action
----|------|
<kbd>↑</kbd>|move up
<kbd>↓</kbd>|move down
<kbd>→</kbd>|go to next page
<kbd>←</kbd>|go to previous page
<kbd>Page Up</kbd>|jump page up
<kbd>Page Down</kbd>|jump page down
<kbd>Home</kbd>|go to first line of page
<kbd>End</kbd>|go to last line of page
<kbd>Enter</kbd>|visit highlighted coin on CoinMarketCap
<kbd>Esc</kbd>|alias to quit
<kbd>Space</kbd>|alias to enter key
<kbd>Ctrl</kbd>+<kbd>c</kbd>|alias to quit
<kbd>Ctrl</kbd>+<kbd>d</kbd>|jump page down (vim style)
<kbd>Ctrl</kbd>+<kbd>n</kbd>|go to next page (vim style)
<kbd>Ctrl</kbd>+<kbd>p</kbd>|go to previous page (vim style)
<kbd>Ctrl</kbd>+<kbd>r</kbd>|force refresh
<kbd>Ctrl</kbd>+<kbd>u</kbd>|jump page up (vim style)
<kbd>0</kbd>|go to first page (vim style)
<kbd>1</kbd>|sort table by *[1] hour change*
<kbd>2</kbd>|sort table by *[2]4 hour change*
<kbd>7</kbd>|sort table by *[7] day change*
<kbd>a</kbd>|sort table by *[a]vailable supply*
<kbd>g</kbd>|go to first line of page  (vim style)
<kbd>G</kbd>|go to last line of page (vim style)
<kbd>h</kbd>|go to previous page (vim style)
<kbd>H</kbd>|go to top of table window (vim style)
<kbd>j</kbd>|move down (vim style)
<kbd>k</kbd>|move up (vim style)
<kbd>l</kbd>|go to next page (vim style)
<kbd>L</kbd>|go to last line of table window (vim style)
<kbd>m</kbd>|sort table by *[m]arket cap*
<kbd>M</kbd>|go to middle of table window (vim style)
<kbd>n</kbd>|sort table by *[n]ame*
<kbd>p</kbd>|sort table by *[p]rice*
<kbd>r</kbd>|sort table by *[r]ank*
<kbd>s</kbd>|sort table by *[s]ymbol*
<kbd>t</kbd>|sort table by *[t]otal supply*
<kbd>u</kbd>|sort table by *last [u]pdated*
<kbd>v</kbd>|sort table by *24 hour [v]olume*
<kbd>q</kbd>|[q]uit
<kbd>$</kbd>|go to last page (vim style)

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
