---
title: "Config"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Config

The first time you run cointop, it'll create a config file in:

```
~/.config/cointop/config.toml
```

On Unix systems, the default config path is `$XDG_CONFIG_HOME/cointop/config.toml`

On macOS (darwin), the default config path is `$HOME/Library/Application Support/cointop/config.toml`

On Windows, the default config path is `%AppData%\cointop\config.toml`

_Note: The config directory is determined by [`os.UserConfigDir()`](https://pkg.go.dev/os#UserConfigDir)_

You may specify a different config file to use by using the `--config` flag:

```bash
cointop --config="/path/to/config.toml"
```

Alternatively, you can set the config file path via the environment variable `COINTOP_CONFIG`

```bash
export COINTOP_CONFIG="/path/to/config.toml"
cointop
```

## Key bindings

You can configure the actions you want for each key in `config.toml`:

```toml
currency = "USD"
default_view = ""
default_chart_range = "1Y"
api = "coingecko"
colorscheme = "cointop"
refresh_rate = 60

[shortcuts]
  "$" = "last_page"
  0 = "move_to_first_page_first_row"
  1 = "sort_column_1h_change"
  2 = "sort_column_24h_change"
  7 = "sort_column_7d_change"
  "?" = "help"
  "/" = "open_search"
  "[" = "previous_chart_range"
  "\\" = "toggle_table_fullscreen"
  "]" = "next_chart_range"
  "{" = "first_chart_range"
  "}" = "last_chart_range"
  "<" = "scroll_left"
  ">" = "scroll_right"
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
  "ctrl+j" = "enlarge_chart"
  "ctrl+k" = "shorten_chart"
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
  tab = "move_down_or_next_page"
  t = "sort_column_total_supply"
  u = "sort_column_last_updated"
  v = "sort_column_24h_volume"

[favorites]

[portfolio]

[coinmarketcap]
  pro_api_key = ""
```

## List of actions

This are the action keywords you may use in the config file to change what the shortcut keys do:

Action|Description
----|------|
`first_chart_range`|Select first chart date range (e.g. 24H)
`first_page`|Go to first page
`move_to_first_page_first_row`|Go to first row on the first page
`enlarge_chart`|Increase chart height
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
`move_down_or_next_page`|Move one row down or to next page if at last row
`move_up_or_previous_page`|Move one row up or to previous page if at first row
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
`scroll_left`|Scroll table to the left
`scroll_right`|Scroll table to the right
`shorten_chart`|Decrease chart height
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
`toggle_table_fullscreen`|Toggle table fullscreen
