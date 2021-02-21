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

Then edit your config `~/.config/cointop/config.toml` and set the colorscheme you want to use:

```toml
colorscheme = "<colorscheme>"
```

The colorscheme name is the name of the colorscheme TOML file.

For example, if you have `matrix.toml` in `~/.cointop/colors/` then the `colorscheme` property should be set to:

```toml
colorscheme = "matrix"
```

Alternatively, you can run cointop with the `--colorscheme` flag to set the colorscheme:

```bash
$ cointop --colorscheme matrix
```

To create your own colorscheme; simply copy an existing [colorscheme](https://github.com/cointop-sh/colors/blob/master/cointop.toml), rename it, and customize the colors.
