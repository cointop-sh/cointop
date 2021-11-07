---
title: "Colorschemes"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Colorschemes

cointop supports custom colorschemes (themes).

<img src="https://user-images.githubusercontent.com/168240/59164231-165b9c80-8abf-11e9-98cf-915ee37407ff.gif" alt="cointop colorschemes" width="880" />

To use standard colorschemes, clone the [colors](https://github.com/cointop-sh/colors) repository into the config directory:

```bash
$ cd ~/.config/cointop
$ git clone git@github.com:cointop-sh/colors.git
```

Note: depending on your system, this may not be the correct location. The "colors" directory needs to go in the same place as your config.toml file.

Then edit your config `~/.config/cointop/config.toml` and set the colorscheme you want to use:

```toml
colorscheme = "cointop"
```

The colorscheme name is the name of the colorscheme TOML file.

By default, the colorscheme files should go under `~/.config/cointop/colors/`

For example, if you have `matrix.toml` under `~/.config/cointop/colors/matrix.toml` then the `colorscheme` property in `config.toml` should be set to:

```toml
colorscheme = "matrix"
```

Alternatively, you can run cointop with the `--colorscheme` flag to set the colorscheme:

```bash
$ cointop --colorscheme matrix
```

To create your own colorscheme; simply copy an existing [colorscheme](https://github.com/cointop-sh/colors/blob/master/cointop.toml), rename it, and customize the colors.

The default `cointop` colorscheme is shown below:

```toml
colorscheme = "cointop"

base_fg = "white"
base_bg = "black"

chart_fg = "white"
chart_bg = "black"
chart_bold = false

marketbar_fg = "white"
marketbar_bg = "black"
marketbar_bold = false

marketbar_label_active_fg = "cyan"
marketbar_label_active_bg = "black"
marketbar_label_active_bold = false

menu_fg = "white"
menu_bg = "black"
menu_bold = false

menu_header_fg = "black"
menu_header_bg = "green"
menu_header_bold = false

menu_label_fg = "yellow"
menu_label_bg = "black"
menu_label_bold = false

menu_label_active_fg = "yellow"
menu_label_active_bg = "black"
menu_label_active_bold = true

searchbar_fg = "white"
searchbar_bg = "black"
searchbar_bold = false

statusbar_fg = "black"
statusbar_bg = "cyan"
statusbar_bold = false

table_column_price_fg = "cyan"
table_column_price_bg = "black"
table_column_price_bold = false

table_column_change_fg = "white"
table_column_change_bg = "black"
table_column_change_bold = false

table_column_change_down_fg = "red"
table_column_change_down_bg = "black"
table_column_change_down_bold = false

table_column_change_up_fg = "green"
table_column_change_up_bg = "black"
table_column_change_up_bold = false

table_header_fg = "black"
table_header_bg = "green"
table_header_bold = false

table_header_column_active_fg = "black"
table_header_column_active_bg = "cyan"
table_header_column_active_bold = false

table_row_fg = "white"
table_row_bg = "black"
table_row_bold = false

table_row_active_fg = "black"
table_row_active_bg = "cyan"
table_row_active_bold = false

table_row_favorite_fg = "yellow"
table_row_favorite_bg = "black"
table_row_favorite_bold = false
```

Supported colors are:

- `black`
- `blue`
- `cyan`
- `green`
- `magenta`
- `red`
- `white`
- `yellow`
- `default` - system default
