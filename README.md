# cointop

> Coin tracking for hackers

[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/miguelmota/cointop/master/LICENSE.md) [![Build Status](https://travis-ci.org/miguelmota/cointop.svg?branch=master)](https://travis-ci.org/miguelmota/cointop) [![Go Report Card](https://goreportcard.com/badge/github.com/miguelmota/cointop?)](https://goreportcard.com/report/github.com/miguelmota/cointop) [![GoDoc](https://godoc.org/github.com/miguelmota/cointop?status.svg)](https://godoc.org/github.com/miguelmota/cointop)

<img src="./assets/screenshot-001.gif" width="880" />

[`cointop`](https://github.com/miguelmota/cointop) is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time. The interface is inspired by [`htop`](https://en.wikipedia.org/wiki/Htop).

## Table of Contents

- [Features](#features)
- [Installing](#install)
- [Updating](#updating)
- [Usage](#usage)
- [Shortcuts](#shortcuts)
- [Config](#config)
- [FAQ](#faq)
- [License](#license)

## Features

- Quick sort shortcuts
- Custom key bindings config
- Vim style keys
- Fast pagination
- 7 day charts for coins and global market graphs
- Fuzzy searching for finding coins
- Save and view favorite coins
- Color coding
- Works on macOS, Linux, and Windows

## Installing

Make sure to have [go](https://golang.org/) (1.4+) installed, then do:

```bash
go get -u github.com/miguelmota/cointop
```

### Windows WSL

You'll need additional font support for Windows WSL. Please see the [wiki](https://github.com/miguelmota/cointop/wiki/Windows-Command-Prompt-and-WSL-Font-Support) for instructions.

### Snap

cointop is also available as a [snap](https://snapcraft.io/cointop) for Linux users.

```bash
sudo snap install cointop --stable
```

Running snap:

```bash
sudo snap run cointop 
```

Note: snaps don't work in Windows WSL. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216/16).

## Updating

To update make sure to use the `-u` flag.

```bash
go get -u github.com/miguelmota/cointop
```

### snap

Use the `refresh` command to update snap.

```bash
sudo snap refresh cointop --stable
```

<!--
#### Alternatively (without go)

```
sudo curl -s "https://raw.githubusercontent.com/miguelmota/cointop/master/install.sh?$(date +%s)" | bash
```
-->

## Usage

Just run the `cointop` command to get started:

```bash
$ cointop
```

## Shortcuts

List of default shortcut keys:

Key|Action
----|------|
<kbd>↑</kbd>|Move up
<kbd>↓</kbd>|Move down
<kbd>→</kbd>|Go to next page
<kbd>←</kbd>|Go to previous page
<kbd>Page Up</kbd>|Jump page up
<kbd>Page Down</kbd>|Jump page down
<kbd>Home</kbd>|Go to first line of page
<kbd>End</kbd>|Go to last line of page
<kbd>Enter</kbd>|Toggle [c]hart for highlighted coin
<kbd>Esc</kbd>|Quit
<kbd>Space</kbd>|Toggle coin as favorite
<kbd>Ctrl</kbd>+<kbd>c</kbd>|Alias to quit
<kbd>Ctrl</kbd>+<kbd>d</kbd>|Jump page down (vim style)
<kbd>Ctrl</kbd>+<kbd>n</kbd>|Go to next page
<kbd>Ctrl</kbd>+<kbd>p</kbd>|Go to previous page
<kbd>Ctrl</kbd>+<kbd>r</kbd>|Force refresh
<kbd>Ctrl</kbd>+<kbd>s</kbd>|Save config
<kbd>Ctrl</kbd>+<kbd>u</kbd>|Jump page up (vim style)
<kbd>Alt</kbd>+<kbd>↑</kbd>|Sort current column in ascending order
<kbd>Alt</kbd>+<kbd>↓</kbd>|Sort current column in descending order
<kbd>Alt</kbd>+<kbd>←</kbd>|Sort column to the left
<kbd>Alt</kbd>+<kbd>→</kbd>|Sort column to the right
<kbd>F1</kbd>|Show help|
<kbd>0</kbd>|Go to first page (vim style)
<kbd>1</kbd>|Sort table by *[1] hour change*
<kbd>2</kbd>|Sort table by *[2]4 hour change*
<kbd>7</kbd>|Sort table by *[7] day change*
<kbd>a</kbd>|Sort table by *[a]vailable supply*
<kbd>c</kbd>|Toggle [c]hart for highlighted coin
<kbd>f</kbd>|Toggle show favorites
<kbd>g</kbd>|Go to first line of page  (vim style)
<kbd>G</kbd>|Go to last line of page (vim style)
<kbd>h</kbd>|Go to previous page (vim style)
<kbd>H</kbd>|Go to top of table window (vim style)
<kbd>j</kbd>|Move down (vim style)
<kbd>k</kbd>|Move up (vim style)
<kbd>l</kbd>|Go to next page (vim style)
<kbd>L</kbd>|Go to last line of visible table window (vim style)
<kbd>m</kbd>|Sort table by *[m]arket cap*
<kbd>M</kbd>|Go to middle of visible table window (vim style)
<kbd>n</kbd>|Sort table by *[n]ame*
<kbd>o</kbd>|[o]pen link to highlighted coin on [CoinMarketCap](https://coinmarketcap.com/)
<kbd>p</kbd>|Sort table by *[p]rice*
<kbd>r</kbd>|Sort table by *[r]ank*
<kbd>s</kbd>|Sort table by *[s]ymbol*
<kbd>t</kbd>|Sort table by *[t]otal supply*
<kbd>u</kbd>|Sort table by *last [u]pdated*
<kbd>v</kbd>|Sort table by *24 hour [v]olume*
<kbd>q</kbd>|[q]uit
<kbd>$</kbd>|Go to last page (vim style)
<kbd>?</kbd>|Show help|
<kbd>/</kbd>|Search (vim style)|

## Config

The first time you run cointop, it'll create a config file in:

```
~/.cointop/config
```

You can then configure the actions you want for each key:

(default `~/.cointop/config`)

```toml
[shortcuts]
  "$" = "last_page"
  0 = "first_page"
  1 = "sort_column_1h_change"
  2 = "sort_column_24h_change"
  7 = "sort_column_7d_change"
  "?" = "help"
  "/" = "open_search"
  G = "move_to_page_last_row"
  H = "move_to_page_visible_first_row"
  L = "move_to_page_visible_last_row"
  M = "move_to_page_visible_middle_row"
  a = "sort_column_available_supply"
  "alt+down" = "sort_column_desc"
  "alt+left" = "sort_left_column"
  "alt+right" = "sort_right_column"
  "alt+up" = "sort_column_asc"
  down = "move_down"
  left = "previous_page"
  right = "next_page"
  up = "move_up"
  c = "toggle_row_chart"
  "ctrl+c" = "quit"
  "ctrl+d" = "page_down"
  "ctrl+n" = "next_page"
  "ctrl+p" = "previous_page"
  "ctrl+r" = "refresh"
  "ctrl+s" = "save"
  "ctrl+u" = "page_up"
  end = "move_to_page_last_row"
  enter = "toggle_row_chart"
  esc = "quit"
  f = "toggle_show_favorites"
  F1 = "help"
  g = "move_to_page_first_row"
  h = "previous_page"
  home = "move_to_page_first_row"
  j = "move_down"
  k = "move_up"
  l = "next_page"
  m = "sort_column_market_cap"
  n = "sort_column_name"
  o = "open_link"
  p = "sort_column_price"
  pagedown = "page_down"
  pageup = "page_up"
  q = "quit"
  r = "sort_column_rank"
  s = "sort_column_symbol"
  space = "toggle_favorite"
  t = "sort_column_total_supply"
  u = "sort_column_last_updated"
  v = "sort_column_24h_volume"
```

## List of actions

Action|Description
----|------|
`first_page`|Go to first page
`help`|Show help
`last_page`|Go to last page
`move_to_page_first_row`|Move to first row on page
`move_to_page_last_row`|Move to last row on page
`move_to_page_visible_first_row`|Move to first visible row on page
`move_to_page_visible_last_row`|Move to last visible row on page
`move_to_page_visible_middle_row`|Move to middle visible row on page
`move_up`|Move one row up
`move_down`|Move one row down
`next_page`|Go to next page
`open_link`|Open row link
`open_search`|Open search field
`page_down`|Move one row down
`page_up`|Scroll one page up
`previous_page`|Go to previous page
`quit`|Quit application
`refresh`|Do a manual refresh on the data
`sort_column_1h_change`|Sort table by column *1 hour change*
`sort_column_24h_change`|Sort table by column *24 hour change*
`sort_column_24h_volume`|Sort table by column *24 hour volume*
`sort_column_7d_change`|Sort table by column *7 day change*
`sort_column_asc`|Sort highlighted column by ascending order
`sort_column_available_supply`|Sort table by column *available supply*
`sort_column_desc`|Sort highlighted column by descending order
`sort_column_last_updated`|Sort table by column *last updated*
`sort_column_market_cap`|Sort table by column *market cap*
`sort_column_name`|Sort table by column *name*
`sort_column_price`|Sort table by column *price*
`sort_column_rank`|Sort table by column *rank*
`sort_column_symbol`|Sort table by column *symbol*
`sort_column_total_supply`|Sort table by column *total supply*
`sort_left_column`|Sort the column to the left of the highlighted column
`sort_right_column`|Sort the column to the right of the highlighted column
`toggle_row_chart`|Toggle the chart for the highlighted row
`toggle_favorite`|Toggle coin as favorite
`toggle_show_favorites`|Toggle show favorites
`save`|Save config

## FAQ

- Q: Where is the data from?

  - A: The data is from [Coin Market Cap](https://coinmarketcap.com/).

- Q: What coins does this support?

  - A: This supports any coin listed on [Coin Market Cap](https://coinmarketcap.com/).

- Q: How often is the data polled?

  - A: Data gets polled once every minute by default. You can press <kbd>Ctrl</kbd>+<kbd>r</kbd> to force refresh.

- Q: I installed cointop without errors but the command is not found.

  - A: Make sure your `GOPATH` and `PATH` is set correctly.

    ```bash
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    ```

- Q: What is the size of the binary?

  - A: The executable is only ~1.9MB in size.

- Q: How do I search?

  - A: The default key to open search is <kbd>/</kbd>. Type the search query after the `/` in the field and hit <kbd>Enter</kbd>.

- Q: Does this work on the Raspberry Pi?

  - A: Yes, cointop works on the Rasperry Pi including the RPi Zero.

- Q: How do I add/remove a favorite?

  - A: Press the <kbd>space</kbd> key to toggle a coin as a favorite.

- Q: How do I view all my favorites?

  - A: Press <kbd>f</kbd> to toggle view all your favorites.

- Q: How do I save my favorites?

  - A: Press <kbd>ctrl</kbd>+<kbd>s</kbd> to save your favorites.

- Q: I'm getting question marks or weird symbols instead of the correct characters.

  - A: Make sure that your terminal has the encoding set to UTF-8 and that your terminal font supports UTF-8.

    You can also try running cointop with the following environment variables:

    ```bash
    LANG=en_US.utf8 TERM=xterm-256color cointop
    ```

    If you're on Windows WSL, please see the [wiki](https://github.com/miguelmota/cointop/wiki/Windows-Command-Prompt-and-WSL-Font-Support) for font support instructions.

- Q: How do I install Go on Ubuntu?

  - A: There's instructions on installing Go on Ubuntu in the [wiki](https://github.com/miguelmota/cointop/wiki/Installing-Go-on-Ubuntu).

- Q: I'm getting errors installing the snap in Windows WSL.

  - A: Unfortunately Windows WSL doesn't support `snapd` which is required for snaps to run. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216/16).

## License

Released under the MIT license.
