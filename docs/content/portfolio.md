---
title: "Portfolio"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# Portfolio

<img src="https://user-images.githubusercontent.com/168240/50439364-a78ade00-08a6-11e9-992b-af63ef21100d.png" alt="portfolio screenshot" width="880" />

## View portfolio

To view your portfolio, press <kbd>P</kbd> (Shift+p)

## Exit portfolio

To exit out of the portfolio view press, <kbd>P</kbd> (Shift+p) again or <kbd>q</kbd> or <kbd>ESC</kbd>

## Add entry

To add a coin to your portfolio, press <kbd>e</kbd> on the highlighted coin, enter a value, and then press <kbd>Enter</kbd>

## Edit entry

To edit the holdings of coin in your portfolio, press <kbd>e</kbd> on the highlighted coin, enter the new value, and then press <kbd>Enter</kbd>

## Remove Entry

To remove an entry in your portfolio, press <kbd>e</kbd> on the highlighted coin and set the value to an empty value and press <kbd>Enter</kbd>

## Changing chart

To change the coin for the chart, press <kbd>Enter</kbd> on the highlighted coin. Pressing <kbd>Enter</kbd> again on the same highlighted row will show the global chart again.

# CLI

The portfolio holdings can be retrieved with the `holdings` command.

### Default holdings table view

```bash
$ cointop holdings
    name        symbol       price      holdings         balance         24h%   %holdings
Bitcoin         BTC       11833.16            10        118331.6        -1.02       74.14
Ethereum        ETH          394.9           100           39490         0.02       24.74
Dogecoin        DOGE    0.00355861        500000          1779.3         1.46        1.11
```

### Output as csv

```bash
$ cointop holdings --format csv
name,symbol,price,holdings,balance,24h%,%holdings
Bitcoin,BTC,11833.16,10,118331.6,-1.02,74.16
Ethereum,ETH,394.48,100,39448,-0.18,24.72
Dogecoin,DOGE,0.00355861,500000,1779.3,1.46,1.12
```

### Output as json

```bash
$ cointop holdings --format json
[{"%holdings":"74.16","24h%":"-1.02","balance":"118331.6","holdings":"10","name":"Bitcoin","price":"11833.16","symbol":"BTC"},{"%holdings":"24.72","24h%":"-0.18","balance":"39448","holdings":"100","name":"Ethereum","price":"394.48","symbol":"ETH"},{"%holdings":"1.12","24h%":"1.46","balance":"1779.3","holdings":"500000","name":"Dogecoin","price":"0.00355861","symbol":"DOGE"}]
```

### Human readable numbers

Adds comma and dollar signs:

```bash
$ cointop holdings -h
    name        symbol        price     holdings           balance        24h%  %holdings
Bitcoin         BTC      $11,833.16           10        $118,331.6      -1.02%     74.14%
Ethereum        ETH          $394.9          100           $39,490       0.02%     24.74%
Dogecoin        DOGE    $0.00355861      500,000          $1,779.3       1.46%      1.11%
```

### Filter coins based on name or symbol

```bash
$ cointop holdings --filter btc,eth
    name        symbol     price        holdings         balance         24h%   %holdings
Bitcoin         BTC     11833.16              10        118331.6        -1.02       74.16
Ethereum        ETH       394.48             100           39448        -0.18       24.72
```

### Filter columns

```bash
$ cointop holdings --cols symbol,holdings,balance
```

### Convert to a different fiat currency

```bash
$ cointop holdings -h --convert eur
    name        symbol    price holdings        balance    24h% %holdings
Ethereum        ETH     €278.49      100        €27,849 -15.87%   100.00%
```

### Combining flags

```bash
$ cointop holdings --total --filter btc,doge --format json -h
{"total":"$120,298.37"}
```

### Help

For all other options, see help command:

```bash
$ cointop holdings --help
```
