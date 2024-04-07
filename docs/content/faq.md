---
title: "Frequently asked questions"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Frequently asked questions

## Where is the data from?

  By default, the data is from [CoinGecko](https://www.coingecko.com/). Data from [CoinMarketCap](https://coinmarketcap.com/) is another option.

## What APIs does it support?

  APIs currently supported are [CoinMarketCap](https://coinmarketcap.com/) and [CoinGecko](https://www.coingecko.com/).

## What coins does this support?

  This supports any coin supported by the API being used to fetch coin information.  There is, however, a limit on the number of coins that
  cointop fetches by default.  You can increase this by passing `--max-pages` and `--per-page` arguments on the command line.

## How do I set the API to use?

  You can use the `--api` flag, eg. `--api coingecko`. You can also set the API choice in the config file.

  ```toml
  api = "coingecko"
  ```

  Options are: `coinmarketcap`, `coingecko`

## How do I change the colorscheme (theme)?

  You can use the `--colorscheme` flag, eg. `--colorscheme matrix`. You can also set the colorscheme choice in the config file.

  ```toml
  colorscheme = "<colorscheme>"
  ```

  For more instructions, visit the [colors](https://github.com/cointop-sh/colors) repository.

## How do I create a custom colorscheme?

  Copy an existing [colorscheme](https://github.com/cointop-sh/colors/blob/master/cointop.toml) to `~/.config/cointop/colors/` and customize the colors. Then run cointop with `--colorscheme <colorscheme>` to use the colorscheme.

  You can use any of the 250-odd X11 colors by name. See https://en.wikipedia.org/wiki/X11_color_names (use lower-case and without spaces).   You can also include 24-bit colors by using the #rrggbb hex code.

  You can also define values in the colorscheme file, and reference them from throughout the file, using the following syntax:

  ```toml
  define_base03 = "#002b36"
  menu_header_fg = "$base03"
  ```

## How do I make the background color transparent?

  Change the background color options in the colorscheme file to `default` to use the system default color, eg. `base_bg = "default"`

## Where is the config file located?

  The default configuration file is located under `~/.config/cointop/config.toml`

  However it may vary depending on your operating system since the config directory is determinedby by [`os.UserConfigDir()`](https://pkg.go.dev/os#UserConfigDir).

  On Unix systems, the default config path is `$XDG_CONFIG_HOME/cointop/config.toml`

  On macOS (darwin), the default config path is `$HOME/Library/Application Support/cointop/config.toml`

  On Windows, the default config path is `%AppData%\cointop\config.toml`

  Note: Previous versions of cointop used `~/.cointop/config` or `~/.cointop/config.toml` as the default config filepath. Cointop will use those config filepaths respectively if they exist.

## What format is the configuration file in?

  The configuration file is in [TOML](https://en.wikipedia.org/wiki/TOML) format.

## Will you be supporting more coin API's in the future?

  Yes supporting more coin APIs is planned.

## How often is the data polled?

  Data gets polled once 60 seconds by default. You can press <kbd>Ctrl</kbd>+<kbd>r</kbd> to force refresh. You can configure the refresh rate with the flag `--refresh-rate <seconds>`

## How can I change the refresh rate?

  Run cointop with the flag `--refresh-rate 60` where the value is the number of seconds that it will fetch for data. You can also set the refresh rate in the config file:

  ```toml
  refresh_rate = 60
  ```

## I ran cointop for the first time and don't see any data?

  Running cointop for the first time will fetch the data and populate the cache which may take a few seconds.

## I'm no longer seeing any data!

  Run cointop with the `--clean` flag to delete the cache. If you're still not seeing any data, then please [submit an issue](https://github.com/cointop-sh/cointop/issues/new).

## How do I get a CoinMarketCap Pro (Paid) API key?

  Create an account on [CoinMarketCap](https://pro.coinmarketcap.com/signup) and visit the [Account](https://pro.coinmarketcap.com/account) page to copy your Pro API key.

## How do I add my CoinMarketCap Pro API key?

  Add the API key in the cointop config file:

  ```toml
  [coinmarketcap]
    pro_api_key = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  ```

  Alternatively, you can export the environment variable `CMC_PRO_API_KEY` containing the API key in your `~/.bashrc`

  ```bash
  export CMC_PRO_API_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  ```

  You may also set the API key on start:

  ```bash
  cointop --coinmarketcap-api-key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  ```

## How do I add my CoinGecko Demo (Free) API key?

  Add the API key in the cointop config file:

  ```toml
  [coingecko]
    api_key = "CG-xxxxxxxxxxxxxxxxxxxxxxxx"
  ```

  Alternatively, you can export the environment variable `COINGECKO_API_KEY` containing the API key in your `~/.bashrc`

  ```bash
  export COINGECKO_API_KEY=CG-xxxxxxxxxxxxxxxxxxxxxxxx
  ```

  You may also set the API key on start:

  ```bash
  cointop --coingecko-api-key=CG-xxxxxxxxxxxxxxxxxxxxxxxx
  ```

## How do I add my CoinGecko Pro (Paid) API key?

  Add the API key in the cointop config file:

  ```toml
  [coingecko]
    pro_api_key = "CG-xxxxxxxxxxxxxxxxxxxxxxxx"
  ```

  Alternatively, you can export the environment variable `COINGECKO_PRO_API_KEY` containing the API key in your `~/.bashrc`

  ```bash
  export COINGECKO_PRO_API_KEY=CG-xxxxxxxxxxxxxxxxxxxxxxxx
  ```

  You may also set the API key on start:

  ```bash
  cointop --coingecko-pro-api-key=CG-xxxxxxxxxxxxxxxxxxxxxxxx
  ```

## I can I add my own API to cointop?

  Fork cointop and add the API that implements the API [interface](https://github.com/cointop-sh/cointop/blob/master/cointop/common/api/interface.go) to [`cointop/cointop/common/api/impl/`](https://github.com/cointop-sh/cointop/tree/master/cointop/common/api/impl). You can use the CoinGecko [implementation](https://github.com/cointop-sh/cointop/blob/master/cointop/common/api/impl/coingecko/coingecko.go) as reference.

## I installed cointop without errors but the command is not found.

  Make sure your `GOPATH` and `PATH` is set correctly.

  ```bash
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```

## How do I search?

  The default key to open search is <kbd>/</kbd>. Type the search query after the `/` in the field and hit <kbd>Enter</kbd>.
  Each search starts from the current cursor position. To search for the same term again, hit <kbd>/</kbd> then <kbd>Enter</kbd>.

  The default behaviour will start to search by symbol first, then it will continues searching by name if there is no result. To search by only symbol, type the search query after `/s:`. To search by only name, type the search query after `/n:`.

## How do I exit search?

  Press <kbd>Esc</kbd> to exit search.

## Does this work on the Raspberry Pi?

  Yes, cointop works on the Rasperry Pi including the RPi Zero.

## How do I add/remove a favorite?

  Press the <kbd>f</kbd> key to toggle a coin as a favorite.

## How do I view all my favorites?

  Press <kbd>F</kbd> (Shift+f) to toggle view all your favorites.

## How do I save my favorites?

  Favorites are autosaved when setting them. You can also press <kbd>ctrl</kbd>+<kbd>s</kbd> to manually save your favorites to the config file.

## What does the yellow asterisk in the row mean?

  The yellow asterisk `*` star means that you've selected that coin to be a favorite.

## My favorites aren't being saved?

  Try running cointop with `--clean` flag to clear the cache which might be causing the problem.

## How do I add a coin to my portfolio?

  Press <kbd>e</kbd> on the highlighted coin to enter holdings and add to your portfolio.

  This dialog supports basic expressions including `+` `-` `*` `/` etc.

## How do I edit the holdings of a coin in my portfolio?

  Press <kbd>e</kbd> on the highlighted coin to edit the holdings.

## How do I remove a coin in my portfolio?

  Press <kbd>e</kbd> on the highlighted coin to edit the holdings and set the value to any empty string (blank value). Set it to `0` if you want to keep the coin without a value.

## How do I view my portfolio?

  Press <kbd>P</kbd> (Shift+p) to toggle view your portfolio.

## How do I save my portfolio?

  Your portfolio is autosaved after you edit holdings. You can also press <kbd>ctrl</kbd>+<kbd>s</kbd> to manually save your portfolio holdings to the config file.

## How do I include buy/cost price in my portfolio?

  Currently there is no UI for this. If you want to include the cost of your coins in the Portfolio screen, you will need to edit your config.toml

  Each coin consists of four values: coin name, coin amount, cost-price, cost-currency.

  For example, the following configuration includes 100 ALGO at USD1.95 each; and 0.1 BTC at AUD50100.83 each.

  ```toml
   holdings = [["Algorand", "100", "1.95", "USD"], ["Bitcoin", "0.1", "50100.83", "AUD"]]
   ```

  With this configuration, four new columns are useful:

  - `cost_price` the price and currency that the coins were purchased at
  - `cost` the cost (in the current currency) of the coins
  - `pnl` the PNL of the coins (current value vs original cost)
  - `pnl_percent` the PNL of the coins as a fraction of the original cost

  With the holdings above, and the currency set to GBP (British Pounds) cointop will look something like this:

  ![portfolio profit and loss](https://user-images.githubusercontent.com/122371/138361142-8e1f32b5-ca24-471d-a628-06968f07c65f.png)

## How do I hide my portfolio balances (private mode)?

  You can run cointop with the `--hide-portfolio-balances` flag to hide portfolio balances or use the keyboard shortcut <kbd>Ctrl</kbd>+<kbd>space</kbd> on the portfolio page to toggle hide/show.

  <img width="880" alt="hide portfolio balances" src="https://user-images.githubusercontent.com/122371/133578568-153af3cc-350d-4744-ac89-dd8f5e317d1d.png" />

## I'm getting question marks or weird symbols instead of the correct characters.

  Make sure that your terminal has the encoding set to UTF-8 and that your terminal font supports UTF-8.

  You can also try running cointop with the following environment variables:

  ```bash
  LANG=en_US.utf8 TERM=xterm-256color cointop
  ```

  If you're on Windows (PowerShell, Command Prompt, or WSL), please see the [wiki](https://github.com/cointop-sh/cointop/wiki/Windows-Command-Prompt-and-WSL-Font-Support) for font support instructions.

## How do I install Go on Ubuntu?

  There's instructions on installing Go on Ubuntu in the [wiki](https://github.com/cointop-sh/cointop/wiki/Installing-Go-on-Ubuntu).

## I'm getting errors installing the snap in Windows WSL.

  Unfortunately Windows WSL doesn't support `snapd` which is required for snaps to run. See this [issue thread](https://forum.snapcraft.io/t/windows-subsystem-for-linux/216).

## How do I fix my GOPATH on Windows?

  Go to Control Panel -> Under _System_ click _Edit the system environment variables_ -> then click the _Environment Variables..._ button -> check the `GOPATH` variable.

  Check the environment variable in PowerShell:

  ```bash
  $ Get-ChildItem Env:GOPATH

  Name                           Value
  ----                           -----
  GOPATH                         C:\Users\alice\go
  ```

## How do I manually build the cointop executable on Windows?

  Here's how to build the executable and run it:

  ```powershell
  > md C:\Users\Josem\go\src\github.com\cointop-sh -ea 0
  > git clone https://github.com/cointop-sh/cointop.git
  > go build -o cointop.exe main.go
  > cointop.exe
  ```

## How do I show the help menu?

  Press <kbd>?</kbd> to toggle the help menu. Press <kbd>q</kbd> to close help menu.

## I'm getting the error: `new gocui: termbox: error while reading terminfo data: EOF` or the error `termbox: error while reading terminfo data: termbox: unsupported terminal` when trying to run.

  Try setting the environment variable `TERM=screen-256color` when starting cointop. E.g. `TERM=screen-256color cointop`

## Does cointop work inside an emacs shell?

  Yes, but it's slightly buggy.

## My shortcut keys are messed or not correct.

  Delete the cointop config directory and rerun cointop.

  ```bash
  rm -rf ~/.cointop
  ```

## How do I display the chart for the highlighted coin?

  Press <kbd>Enter</kbd> to toggle the chart for the highlighted coin.

## How do I change the chart date range?

  Press <kbd>]</kbd> to cycle to the next date range.

  Press <kbd>[</kbd> to cycle to the previous date range.

  Press <kbd>{</kbd> to select the first date range.

  Press <kbd>}</kbd> to selected the last date range.

## What chart date ranges are supported?

  Supported date ranges are `All Time`, `YTD`, `1Y`, `6M`, `3M`, `1M`, `7D`, `3D`, `24H`.

  <sup><sub>YTD = Year-to-date<sub></sup>

## How do I expand or shorten the chart?

  Press <kbd>Ctrl</kbd>+<kbd>j</kbd> to increase the chart height.

  Press <kbd>Ctrl</kbd>+<kbd>k</kbd> to decrease the chart height.

## How do I change the fiat currency?

  Press <kbd>c</kbd> to show the currency convert menu, and press the corresponding key to select that as the fiat currency.

## Which currencies can I convert to?

  The supported fiat currencies for conversion are `AUD`, `BRL`, `CAD`, `CFH`, `CLP`, `CNY`, `CZK`, `DKK`, `EUR`, `GBP`, `HKD`, `HUF`, `IDR`, `ILS`, `INR`, `JPY`, `KRW`, `MXN`, `MYR`, `NOK`, `NZD`, `PLN`, `PHP`, `PKR`, `RUB`, `SEK`, `SGD`, `THB`, `TRY`, `TWD`, `USD`,  `VND`, and `ZAR`.

  The supported crypto currencies for conversion are `BTC` and `ETH`.

  Please note that some APIs may have limited support for certain conversion formats.

## How do I save the selected currency to convert to?

  The selected currency conversion is autosaved. You can also press <kbd>ctrl</kbd>+<kbd>s</kbd> to manually save the selected currency conversion.

## How do change the format for currency and numeric values?

  Cointop uses the `LC_MONETARY` and `LC_NUMERIC` environment variables to format
  currency and numeric values.

  For example, to see cointop display currency with Bengalese numbers
  and numeric values in arabic run `/usr/bin/env LC_MONETARY=bn LC_NUMERIC=ar cointop`.

  For more information about how to check and configure your locale settings see
  [`locale(1)`](https://www.unix.com/man-page/linux/1/locale/) and
  [`locale(5)`](https://www.unix.com/man-page/linux/5/locale/).

## What does saving do?

  The save command (<kbd>ctrl</kbd>+<kbd>s</kbd>) saves your selected currency, selected favorite coins, and portfolio coins to the cointop config file.

## The data isn't refreshing!

  The coin APIs have rate limits, so make sure to keep manual refreshes to a minimum. If you've hit the rate limit then wait about half an hour to be able to fetch the data again. Keep in mind that some coin APIs, such as CoinMarketCap, update prices every 5 minutes so constant refreshes aren't necessary.

## How do I quit the application?

  Press <kbd>ctrl</kbd>+<kbd>c</kbd> to quit the application.

## How do I quit the open view/window?

  Press <kbd>q</kbd> to quit the open view/window.

## How do I set the favorites view to be the default view?

  In the config file, set `default_view = "favorites"`

## How do I set the portfolio view to be the default view?

  In the config file, set `default_view = "portfolio"`

## How do I set the table view to be the default view?

  In the config file, set `default_view = "default"`

## How do I set the default chart range?

  In the config file, set `default_chart_range = "3M"`

  Supported date ranges are `All Time`, `YTD`, `1Y`, `6M`, `3M`, `1M`, `7D`, `3D`, `24H`.

## How do I set the table columns?

  In the config file, set `columns` value in one of the favorites, portfolio, table sections. The list of available columns can be seen on the help menu - look for "sort_column_XXX".

  For example:

  ```toml
  [table]
    columns = ["rank", "name", "symbol", "price", "1h_change", "24h_change", "7d_change", "24h_volume", "market_cap", "available_supply", "total_supply", "last_updated"]
  ```

## What price-change columns are available?

  Supported columns relating to price change are `1h_change`, `24h_change`, `7d_change`, `30d_change`, `1y_change`

## How can I use K (thousand), M (million), B (billion), T (trillion) suffixes for shorter numbers?

  There is a setting at the top-level of the configuration file called `compact_notation=true` which changes the marketbar values `market cap`, `volume` and `portfolio total value`.

  The same setting can be applied at in the `[table]` section to impact the `24h_volume`, `market_cap`, `total_supply`, `available_supply` columns in the main coin view; and in the `[favorites]` section to change the same columns.   The setting also changes the column names to be shorter.

## How can use a different config file other than the default?

  Run cointop with the `--config` flag, eg `cointop --config="/path/to/config.toml"`, to use the specified file as the config.

  Alternatively, you can set the config file path via the environment variable `COINTOP_CONFIG`, eg `export COINTOP_CONFIG="/path/to/config.toml"`

## I'm getting the error `open /dev/tty: no such device or address`.

  Usually this error occurs when cointop is running as a daemon or slave which means that there is no terminal allocated, so `/dev/tty` doesn't exist for that process. Try running it with the following environment variables:

  ```bash
  DEV_IN=/dev/stdout DEV_OUT=/dev/stdout cointop
  ```

## I can only view the first page, why isn't the pagination is working?

  Sometimes the coin APIs will make updates and break things. If you see this problem please [submit an issue](https://github.com/cointop-sh/cointop/issues/new).

## How can run cointop with just the table?

  Run cointop with the `--only-table` flag.

    <img width="880" alt="table view only" src="https://user-images.githubusercontent.com/168240/60208658-b0387e80-980d-11e9-8819-8039fb11218f.png" />

## How do I toggle the table to go fullscreen?

  Press <kbd>\\</kbd> to toggle the table fullscreen mode.

## How can I hide the top marketbar?

  Run cointop with the `--hide-marketbar` flag.

## How can I hide the chart?

  Run cointop with the `--hide-chart` flag.

## How can I hide the bottom statusbar?

  Run cointop with the `--hide-statusbar` flag.

## How do I scroll the table horizontally left or right?

  Use the keys <kbd><</kbd> to scroll the table to the left and <kbd><</kbd> to scroll the table to the right.

## How can I delete the cache?

  Run `cointop clean` to delete the cache files. Cointop will generate new cache files after fetching data.

## How can I reset cointop?

  Run the command `cointop reset` to delete the config files and cache. Cointop will generate a new config when starting up. You can run `cointop --reset` to reset before running cointop.

## Why aren't <kbd>Home</kbd> or <kbd>End</kbd> keys working for me?

  Make sure to not manually set `TERM` in your `~/.bashrc`, `~/.zshrc`, or any where else. The `TERM` environment variable should be automatically set by your terminal, otherwise this may cause the terminal emulator to send escape codes when these keys are pressed. See [this Arch wiki](https://wiki.archlinux.org/index.php/Home_and_End_keys_not_working) for more info.

  Use the `cointop price` command. Here are some examples:

## What is the size of the binary?

  The Go build size is ~8MB but packed with UPX it's only a ~3MB executable binary.

## How much memory does cointop use?

  Cointop uses ~15MB of memory so you can run it on a Raspberry Pi Zero if you wanted to (one reason why cointop was built using Go instead of Node.js or Python).

## How does cointop differ from *rate.sx*?

  [rate.sx](https://rate.sx/) is great for one-off queries or fetching data for bash scripts because it doesn't require installing anything. Cointop differs in that it is interactive and also supports more currencies.

## How can I get just the coin price with cointop?

  Use the `cointop price` command. Here are some examples:

  ```bash
  $ cointop price --coin ethereum
  $277.76

  $ cointop price -c ethereum --currency btc
  Ƀ0.02814

  $ cointop -c ethereum -f eur
  €245.51

  $ cointop price -c ethereum -f usd --api coinmarketcap
  $276.37
  ```

## Does cointop do mining?

  Cointop does not do any kind of cryptocurrency mining.

## How can I run the cointop SSH server on port 22?

  Port 22 is a privileged port so you need to run with `sudo`:

  ```bash
  sudo cointop server -p 22
  ```

## Does using the public SSH instance track me?

  No, there's no tracking on the public SSH server instance. It's running the same code that you see on the public repo.

## Why does the public SSH server kick me off?

  The public SSH instance is for demo use or short-term usage and will disconnect users after an idle timeout to allow other users to try it out.

## How do I fix the error `SSH key is required to start server` when trying to run SSH server?

  Generate the SSH key if the host machine doesn't have one:

  ```bash
  $ ssh-keygen
  ```

  Make sure the host key flag points to the location of the SSH key:

  ```bash
  cointop server -k ~/.ssh/id_rsa [...]
  ```

## How do I fix the error `no matching host key type found. Their offer: ssh-rsa` when trying to SSH?

Use the following flag when connecting to the SSH server:

  ```bash
  ssh -oHostKeyAlgorithms=+ssh-rsa cointop.sh
  ```

You can also add this config to the `~/.ssh/config` file so you don't have to use the flag every time:

```
Host cointop.sh
  HostName cointop.sh
  HostKeyAlgorithms=+ssh-rsa
```

## Why doesn't the version number work when I install with `go get`?

  The version number is read from the git tag during the build process but this requires the `GO111MODULE` environment variable to be set in order for Go to read the build information:

  ```bash
  GO111MODULE=on go get github.com/cointop-sh/cointop
  ```

## How can I get more information when something is going wrong?

  Cointop creates a logfile at `/tmp/cointop.log`. Normally nothing is written to this, but if you set the environment variable
  `DEBUG=1` cointop will write a lot of output describing its operation.  Furthermore, if you also set `DEBUG_HTTP=1` it will
  emit lots about every HTTP request that cointop makes to coingecko (backend).  Developers may ask for this information
  to help diagnose any problems you may experience.

  ```bash
  DEBUG=1 DEBUG_HTTP=1 cointop
  ```

  If you set environment variable `DEBUG_FILE` you can explicitly provide a logfile location, rather than `/tmp/cointop.log`
