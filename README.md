<h1 align="center">
  <br />
  <img src="https://user-images.githubusercontent.com/168240/39561871-51cda852-4e5d-11e8-926b-7692d43143e8.png" alt="cointop" width="400" />
  <br />
  <br />
  <br />
</h1>

> Coin tracking for hackers

[![License](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/miguelmota/cointop/master/LICENSE.md) [![Build Status](https://travis-ci.org/miguelmota/cointop.svg?branch=master)](https://travis-ci.org/miguelmota/cointop) [![Go Report Card](https://goreportcard.com/badge/github.com/miguelmota/cointop?)](https://goreportcard.com/report/github.com/miguelmota/cointop) [![GoDoc](https://godoc.org/github.com/miguelmota/cointop?status.svg)](https://godoc.org/github.com/miguelmota/cointop) [![Mentioned in Awesome Terminals](https://awesome.re/mentioned-badge.svg)](https://github.com/k4m4/terminals-are-sexy)

[`cointop`](https://github.com/miguelmota/cointop) is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time.

The interface is inspired by [`htop`](https://en.wikipedia.org/wiki/Htop) and shortcut keys are inspired by [`vim`](https://en.wikipedia.org/wiki/Vim_(text_editor)).

<img src="https://user-images.githubusercontent.com/168240/39569578-7ce9f3b6-4e7a-11e8-82a9-8a18b91b1bd5.png" alt="cointop screenshot" width="880" />

<img src="https://user-images.githubusercontent.com/168240/39569662-bcbdbcc0-4e7a-11e8-8a8f-8ff45868a8ae.png" alt="help menu" width="880" />

<img src="https://user-images.githubusercontent.com/168240/41806841-043c0ca6-767a-11e8-9c51-df9fc64b3b5c.png" alt="currency convert menu" width="880" />

In action

<img src="https://user-images.githubusercontent.com/168240/39569570-75b1547c-4e7a-11e8-8eac-552abaa431f0.gif" alt="screencast" width="880" />

## Table of Contents

- [Features](#features)
- [Installing](#install)
- [Updating](#updating)
- [Getting started](#getting-started)
- [Shortcuts](#shortcuts)
- [Config](#config)
- [FAQ](#faq)
- [Development](#development)
- [License](#license)

## Features

- Quick sort shortcuts
- Custom key bindings config
- Vim inspired shortcut keys
- Fast pagination
- Charts for coins and global market graphs
- Quick chart date range change
- Fuzzy searching for finding coins
- Currency conversion
- Save and view favorite coins
- Portfolio tracking of holdings
- Color support
- Help menu
- Offline cache
- Works on macOS, Linux, and Windows
- It's very lightweight; can be left running indefinitely

## Installing

There are multiple ways you can install cointop depending on the platform you're on.

### From source (always latest)

Make sure to have [go](https://golang.org/) (1.9+) installed, then do:

```bash
go get -u github.com/miguelmota/cointop
```

Make sure `$GOPATH/bin` is added to the `$PATH` variable.

Now you can run cointop:

```bash
cointop
```

### Homebrew (macOS)

cointop is available via [Homebrew](https://formulae.brew.sh/formula/cointop) for macOS:

```bash
brew install cointop
```

Run

```bash
cointop
```

### Flatpak (Linux)

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

cointop flathub [project page](https://flathub.org/apps/details/com.github.miguelmota.Cointop)

### Copr (Fedora)

cointop is available as a [copr](https://copr.fedorainfracloud.org/) package.

First, enable the respository

```bash
sudo dnf copr enable miguelmota/cointop
```

Install cointop

````bash
sudo dnf install cointop
````

Run

```bash
cointop
```

cointop copr [project page](https://copr.fedorainfracloud.org/coprs/miguelmota/cointop/).

### Snap (Ubuntu)

cointop is available as a [snap](https://snapcraft.io/cointop) for Linux users.

```bash
sudo snap install cointop --stable
```

Running snap:

```bash
sudo snap run cointop
```

cointop snapcraft [project page](https://snapcraft.io/cointop).

Note: snaps don't work in Windows WSL. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216).

### Windows WSL (Windows)

Recommended to install using Go (instructions above).

You'll need additional font support for Windows WSL. Please see the [wiki](https://github.com/miguelmota/cointop/wiki/Windows-Command-Prompt-and-WSL-Font-Support) for instructions.

### Binaries

You can find pre-built binaries on the [releases](https://github.com/miguelmota/cointop/releases) page.

## Updating

To update make sure to use the `-u` flag if installed via Go.

```bash
go get -u github.com/miguelmota/cointop
```

### Homebrew (macOS)

```bash
brew uninstall cointop && brew install cointop
```

### Flatpak (Linux)

```bash
sudo flatpak uninstall com.github.miguelmota.Cointop
sudo flatpak install flathub com.github.miguelmota.Cointop
```

### Copr (Fedora)

```bash
sudo dnf update cointop
```

### Snap (Ubuntu)

Use the `refresh` command to update snap.

```bash
sudo snap refresh cointop --stable
```

## Getting started

Just run the `cointop` command to get started:

```bash
$ cointop
```

### Navigation

- Easiest way to navigate up and down is using the arrow keys <kbd>↑</kbd> and <kbd>↓</kbd>, respectively
- To go the next and previous pages, use <kbd>→</kbd> and <kbd>←</kbd>, respectively
- To go to the top and bottom of the page, use <kbd>g</kbd> and <kbd>G</kbd> (Shift+g), respectively
- Check out the rest of [shortcut](#shortcuts) keys for vim-inspired navigation

### Favorites

- To toggle a coin as a favorite, press <kbd>Space</kbd> on the highlighted coin
- To view all your favorite coins, press <kbd>F</kbd> (Shift+f)
- To exit out of the favorites view, press <kbd>F</kbd> (Shift+f) again or <kbd>q</kbd>

### Portfolio

- To add a coin to your portfolio, press <kbd>e</kbd> on the highlighted coin
- To edit the holdings of coin in your portfolio, press <kbd>e</kbd> on the highlighted coin
- To view your portfolio, press <kbd>P</kbd> (Shift+p)
- To exit out of the portfolio view press, <kbd>P</kbd> (Shift+p) again or <kbd>q</kbd>

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
<kbd>Enter</kbd>|Toggle chart for highlighted coin
<kbd>Esc</kbd>|Quit view
<kbd>Space</kbd>|Toggle coin as favorite
<kbd>Ctrl</kbd>+<kbd>c</kbd>|Quit application
<kbd>Ctrl</kbd>+<kbd>d</kbd>|Jump page down (vim inspired)
<kbd>Ctrl</kbd>+<kbd>f</kbd>|Search
<kbd>Ctrl</kbd>+<kbd>n</kbd>|Go to next page
<kbd>Ctrl</kbd>+<kbd>p</kbd>|Go to previous page
<kbd>Ctrl</kbd>+<kbd>r</kbd>|Force refresh data
<kbd>Ctrl</kbd>+<kbd>s</kbd>|Save config
<kbd>Ctrl</kbd>+<kbd>u</kbd>|Jump page up (vim inspired)
<kbd>Alt</kbd>+<kbd>↑</kbd>|Sort current column in ascending order
<kbd>Alt</kbd>+<kbd>↓</kbd>|Sort current column in descending order
<kbd>Alt</kbd>+<kbd>←</kbd>|Sort column to the left
<kbd>Alt</kbd>+<kbd>→</kbd>|Sort column to the right
<kbd>F1</kbd>|Show help|
<kbd>F5</kbd>|Force refresh data|
<kbd>0</kbd>|Go to first page (vim inspired)
<kbd>1</kbd>|Sort table by *[1] hour change*
<kbd>2</kbd>|Sort table by *[2]4 hour change*
<kbd>7</kbd>|Sort table by *[7] day change*
<kbd>a</kbd>|Sort table by *[a]vailable supply*
<kbd>b</kbd>|Sort table by *[b]alance*
<kbd>c</kbd>|Show currency convert menu
<kbd>C</kbd>|Show currency convert menu
<kbd>e</kbd>|Show portfolio edit holdings menu
<kbd>E</kbd> (Shift+e)|Show portfolio edit holdings menu
<kbd>f</kbd>|Toggle coin as favorite
<kbd>F</kbd> (Shift+f)|Toggle show favorites
<kbd>g</kbd>|Go to first line of page  (vim inspired)
<kbd>G</kbd> (Shift+g)|Go to last line of page (vim inspired)
<kbd>h</kbd>|Go to previous page (vim inspired)
<kbd>h</kbd>|Sort table by *[h]oldings* (portfolio view only)
<kbd>H</kbd> (Shift+h)|Go to top of table window (vim inspired)
<kbd>j</kbd>|Move down (vim inspired)
<kbd>k</kbd>|Move up (vim inspired)
<kbd>l</kbd>|Go to next page (vim inspired)
<kbd>L</kbd> (Shift+l)|Go to last line of visible table window (vim inspired)
<kbd>m</kbd>|Sort table by *[m]arket cap*
<kbd>M</kbd> (Shift+m)|Go to middle of visible table window (vim inspired)
<kbd>n</kbd>|Sort table by *[n]ame*
<kbd>o</kbd>|[o]pen link to highlighted coin on [CoinMarketCap](https://coinmarketcap.com/)
<kbd>p</kbd>|Sort table by *[p]rice*
<kbd>P</kbd> (Shift+p)|Toggle show portfolio
<kbd>r</kbd>|Sort table by *[r]ank*
<kbd>s</kbd>|Sort table by *[s]ymbol*
<kbd>t</kbd>|Sort table by *[t]otal supply*
<kbd>u</kbd>|Sort table by *last [u]pdated*
<kbd>v</kbd>|Sort table by *24 hour [v]olume*
<kbd>q</kbd>|Quit view
<kbd>$</kbd>|Go to last page (vim inspired)
<kbd>?</kbd>|Show help|
<kbd>/</kbd>|Search (vim inspired)|
<kbd>]</kbd>|Next chart date range|
<kbd>[</kbd>|Previous chart date range|
<kbd>}</kbd>|Last chart date range|
<kbd>{</kbd>|First chart date range|

## Config

The first time you run cointop, it'll create a config file in:

```
~/.cointop/config
```

You can then configure the actions you want for each key:

(default `~/.cointop/config`)

```toml
currency = "USD"
defaultView = "default"

[shortcuts]
  "$" = "last_page"
  0 = "first_page"
  1 = "sort_column_1h_change"
  2 = "sort_column_24h_change"
  7 = "sort_column_7d_change"
  "?" = "help"
  "/" = "open_search"
  "[" = "previous_chart_range"
  "]" = "next_chart_range"
  "{" = "first_chart_range"
  "}" = "last_chart_range"
  C = "show_currency_convert_menu"
  E = "show_portfolio_edit_menu"
  G = "move_to_page_last_row"
  H = "move_to_page_visible_first_row"
  L = "move_to_page_visible_last_row"
  M = "move_to_page_visible_middle_row"
  O = "open_link"
  P = "toggle_portfolio"
  a = "sort_column_available_supply"
  "alt+down" = "sort_column_desc"
  "alt+left" = "sort_left_column"
  "alt+right" = "sort_right_column"
  "alt+up" = "sort_column_asc"
  down = "move_down"
  left = "previous_page"
  right = "next_page"
  up = "move_up"
  c = "show_currency_convert_menu"
  b = "sort_column_balance"
  "ctrl+c" = "quit"
  "ctrl+d" = "page_down"
  "ctrl+f" = "open_search"
  "ctrl+n" = "next_page"
  "ctrl+p" = "previous_page"
  "ctrl+r" = "refresh"
  "ctrl+s" = "save"
  "ctrl+u" = "page_up"
  e = "show_portfolio_edit_menu"
  end = "move_to_page_last_row"
  enter = "toggle_row_chart"
  esc = "quit"
  f = "toggle_favorite"
  F = "toggle_show_favorites"
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
  q = "quit_view"
	Q = "quit_view"
  r = "sort_column_rank"
  s = "sort_column_symbol"
  space = "toggle_favorite"
  t = "sort_column_total_supply"
  u = "sort_column_last_updated"
  v = "sort_column_24h_volume"
```

## List of actions

This are the action keywords you may use in the config file to change what the shortcut keys do.

Action|Description
----|------|
`first_chart_range`|Select first chart date range (e.g. 1H)
`first_page`|Go to first page
`help`|Show help
`hide_currency_convert_menu`|Hide currency convert menu
`last_chart_range`|Select last chart date range (e.g. All Time)
`last_page`|Go to last page
`move_to_page_first_row`|Move to first row on page
`move_to_page_last_row`|Move to last row on page
`move_to_page_visible_first_row`|Move to first visible row on page
`move_to_page_visible_last_row`|Move to last visible row on page
`move_to_page_visible_middle_row`|Move to middle visible row on page
`move_up`|Move one row up
`move_down`|Move one row down
`next_chart_range`|Select next chart date range (e.g. 3D → 7D)
`next_page`|Go to next page
`open_link`|Open row link
`open_search`|Open search field
`page_down`|Move one row down
`page_up`|Scroll one page up
`previous_chart_range`|Select previous chart date range (e.g. 7D → 3D)
`previous_page`|Go to previous page
`quit`|Quit application
`quit_view`|Quit view
`refresh`|Do a manual refresh on the data
`save`|Save config
`show_currency_convert_menu`|Show currency convert menu
`show_favorites`|Show favorites
`sort_column_1h_change`|Sort table by column *1 hour change*
`sort_column_24h_change`|Sort table by column *24 hour change*
`sort_column_24h_volume`|Sort table by column *24 hour volume*
`sort_column_7d_change`|Sort table by column *7 day change*
`sort_column_asc`|Sort highlighted column by ascending order
`sort_column_available_supply`|Sort table by column *available supply*
`sort_column_balance`|Sort table by column *balance*
`sort_column_desc`|Sort highlighted column by descending order
`sort_column_holdings`|Sort table by column *holdings*
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
`toggle_show_currency_convert_menu`|Toggle show currency convert menu
`toggle_show_favorites`|Toggle show favorites
`toggle_portfolio`|Toggle portfolio view
`toggle_show_portfolio`|Toggle show portfolio view
`show_portfolio_edit_menu`|Show portfolio edit holdings menu

## FAQ

Frequently asked questions:

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

- Q: Where is the config file located?

  - A: The default configuration file is located under `~/.cointop/config`

- Q: What format is the configuration file in?

  - A: The configuration file is in [TOML](https://en.wikipedia.org/wiki/TOML) format.

- Q: How do I search?

  - A: The default key to open search is <kbd>/</kbd>. Type the search query after the `/` in the field and hit <kbd>Enter</kbd>.

- Q: How do I exit search?

  - A: Press <kbd>ESC</kbd> to exit search.

- Q: Does this work on the Raspberry Pi?

  - A: Yes, cointop works on the Rasperry Pi including the RPi Zero.

- Q: How do I add/remove a favorite?

  - A: Press the <kbd>f</kbd> key to toggle a coin as a favorite.

- Q: How do I view all my favorites?

  - A: Press <kbd>F</kbd> (Shift+f) to toggle view all your favorites.

- Q: How do I save my favorites?

  - A: Press <kbd>ctrl</kbd>+<kbd>s</kbd> to save your favorites.

- Q: What does the yellow asterisk in the row mean?

  - A: The yellow asterisk or star means that you've selected that coin to be a favorite.

- Q: How do I add a coin to my portfolio?

  - Press <kbd>e</kbd> on the highlighted coin to enter holdings and add to your portfolio.

- Q: How do I edit the holdings of a coin in my portfolio?

  - Press <kbd>e</kbd> on the highlighted coin to edit the holdings.

- Q: How do I remove a coin in my portfolio?

  - Press <kbd>e</kbd> on the highlighted coin to edit the holdings and set the value to any empty string (blank value). Set it to `0` if you want to keep the coin without a value.

- Q: How do I view my portfolio?

  - A: Press <kbd>P</kbd> (Shift+p) to toggle view your portfolio.

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

  - A: Unfortunately Windows WSL doesn't support `snapd` which is required for snaps to run. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216).

- Q: How do I show the help menu?

  - A: Press <kbd>?</kbd> to toggle the help menu. Press <kbd>q</kbd> to close help menu.

- Q: I'm getting the error: `new gocui: termbox: error while reading terminfo data: EOF` when trying to run.

  - A: Try setting the environment variable `TERM=screen-256color`

- Q: Does cointop work inside an emacs shell?

  - A: Yes, but it's slightly buggy.

- Q: My shortcut keys are messed or not correct.

  - A: Delete the cointop config directory and rerun cointop.

    ```bash
    rm -rf ~/.cointop
    ```

- Q: How do I display the chart for the highlighted coin?

  - A: Press <kbd>Enter</kbd> to toggle the chart for the highlighted coin.

- Q: How do I change the chart date range?

  - A: Press <kbd>]</kbd> to cycle to the next date range.

    Press <kbd>[</kbd> to cycle to the previous date range.

    Press <kbd>{</kbd> to select the first date range.

    Press <kbd>}</kbd> to selected the last date range.

- Q: What chart date ranges are supported?

  - A: Supported date ranges are `All Time`, `YTD`, `1Y`, `6M`, `3M`, `1M`, `7D`, `3D`, `24H`, `6H`, `1H`.

    <sup><sub>YTD = Year-to-date<sub></sup>

- Q: How do I change the fiat currency?

  - A: Press <kbd>c</kbd> to show the currency convert menu, and press the corresponding key to select that as the fiat currency.

- Q: Which currencies can I convert to?

  - A: The supported fiat currencies for conversion are `AUD`, `BRL`, `CAD`, `CFH`, `CLP`, `CNY`, `CZK`, `DKK`, `EUR`, `GBP`, `HKD`, `HUF`, `IDR`, `ILS`, `INR`, `JPY`, `KRW`, `MXN`, `MYR`, `NOK`, `NZD`, `PLN`, `PHP`, `PKR`, `RUB`, `SEK`, `SGD`, `THB`, `TRY`, `TWD`, `USD`, and `ZAR`.

    The supported crypto currencies for conversion are `BTC` and `ETH`.

- Q: How do I save the selected currency to convert to?

  - A: Press <kbd>ctrl</kbd>+<kbd>s</kbd> to save the selected currency to convert to.

- Q: What does saving do?

  - A: The save command (<kbd>ctrl</kbd>+<kbd>s</kbd>) saves your selected currency, selected favorite coins, and portfolio coins to the cointop config file.

- Q: The data isn't refreshing!

  - A: The CoinMarketCap API has rate limits, so make sure to keep manual refreshes to a minimum. If you've hit the rate limit then wait about half an hour to be able to fetch the data again. Keep in mind that CoinMarketCap updates prices every 5 minutes so constant refreshes aren't necessary.

- Q: How do I quit the application?

  - A: Press <kbd>ctrl</kbd>+<kbd>c</kbd> to quit the application.

- Q: How do I quit the open view/window?

  - A: Press <kbd>q</kbd> to quit the open view/window.

- Q: How do I set the favorites view to be the default view?

  - A: In `~/.cointop/config`, set `defaultView = "favorites"`

- Q: How do I set the portfolio view to be the default view?

  - A: In `~/.cointop/config`, set `defaultView = "portfolio"`

- Q: How do I set the table view to be the default view?

  - A: In `~/.cointop/config`, set `defaultView = "default"`

- Q: I'm getting the error `open /dev/tty: no such device or address`.

    -A: Usually this error occurs when cointop is running as a daemon or slave which means that there is no terminal allocated, so `/dev/tty` doesn't exist for that process. Try running it with the following environment variables:

    ```bash
    DEV_IN=/dev/stdout DEV_OUT=/dev/stdout cointop
    ```

- Q: What is the size of the binary?

  - A: The executable is only ~2MB in size.

## Development

### Go

Running cointop from source

```
make run
```

### Homebrew

Installing from source

```bash
make brew/build
```

### Flatpak

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
make flatpak/build
```

### Copr

Install dependencies

```bash
make copr/install/cli
make rpm/install/deps
make rpm/dirs
```

Build package

```bash
make rpm/cp/specs
make rpm/download
make rpm/build
make copr/build
```

### Snap

Building snap

```bash
make snap/build
```

## Deployment

See this [wiki](https://github.com/miguelmota/cointop/wiki/Deployment).

## License

Released under the [Apache 2.0](./LICENSE.md) license.
