# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.6.2] - 2021-02-12
### Added
- Add config option to keep row focus on sort

## [1.6.1] - 2021-02-12
### Fixed
- Chart data interpolation
- CoinMarketCap graph data endpoint

### Added
- Multiple coin support in price command

## [1.6.0] - 2021-02-12
### Fixed
- Coin chart lookup
- Dynamic column widths

### Added
- Configurable table columns
- Basic price alerts

## [1.5.5] - 2020-11-15
### Fixed
- Termux cache directory
- Open command on Windows

### Added
- Currency convesion option to holdings command
- Sort by percent holdings shortcut

## [1.5.4] - 2020-08-24
### Fixed
- Rank order for low market cap coins

### Added
- Colorschemes directory flag

## [1.5.3] - 2020-08-14
### Fixed
- Build error

## [1.5.2] - 2020-08-13
### Fixed
- `XDG_CONFIG_HOME` config path

### Added
- Holdings command with sorting and filter options
- Bitcoin dominance command

## [1.5.1] - 2020-08-05
### Fixed
- Version typo

## [1.5.0] - 2020-08-05
### Fixed
- Use version string from go build info

## [1.4.8] - 2020-08-03
### Added
- No cache flag

## [1.4.7] - 2020-08-02
### Added
- SSH server

### Fixed
- Config flag

## [1.4.6] - 2020-05-23
### Fixed
- Decimals places for BTC and ETH currency conversion
- Increase number of page results from CoinGecko

## [1.4.5] - 2020-02-18
### Fixed
- Convert to chosen currency for market data

### Added
- VND currency conversion

## [1.4.4] - 2019-12-31
### Fixed
- Flathub app release version

## [1.4.3] - 2019-12-29
### Fixed
- Chart update bug fixes
- Marketbar currency bug fixes

### Added
- Tab keybinding

## [1.4.2] - 2019-12-29
### Fixed
- Fix keybinding issue on FreeBSD

## [1.4.1] - 2019-11-17
### Fixed
- Fix version ldflags

## [1.4.0] - 2019-11-17
### Added
- Keyboard shortcuts to enlarge and shorten chart

## [1.3.6] - 2019-09-15
### Fixed
- Fixed various navigation and view switching issues

## [1.3.5] - 2019-09-08
### Fixed
- Fixed table sorting issues

## [1.3.4] - 2019-07-05
### Fixed
- Fixed Windows path

## [1.3.3] - 2019-06-30
### Added
- Added price command

## [1.3.2] - 2019-06-30
### Added
- Toggle table fullscreen shortcut and hide view flags

## [1.3.1] - 2019-06-26
### Added
- Show only table option

### Fixed
- CoinGecko prices

## [1.3.0] - 2019-06-09
### Added
- Added colorscheme support

## [1.2.2] - 2019-06-01
### Fixed
- Market bar background color

## [1.2.1] - 2019-06-01
### Fixed
- Added mutex lock when filecaching

## [1.2.0] - 2019-05-12
### Added
- Added CoinGecko API support

### Changed
- Default API from CoinMarketCap to CoinGecko

## [1.1.6] - 2019-04-23
### Added
- Prompt for CoinMarketCap Pro API Key

## [1.1.5] - 2019-04-22
### Fixed
- Release archive to contain latest source code

## [1.1.4] - 2019-04-21
### Changed
- CoinMarketCap legacy V2 API to Pro V1 API

### Added
- Config option to use CoinMarketCap Pro V1 API KEY

## [1.1.3] - 2019-02-25
### Fixed
- Vendor dependencies

## [1.1.2] - 2018-12-30
### Fixed
- Paginate CoinMarketCap V1 API responses due to their backward-incompatible update

### Added
- `-clean` flag to clean cache
- `-reset` flag to clean cache and delete config
- `-config` flag to use a different specified config file

## [1.1.1] - 2018-12-26
### Changed
- Use go modules instead of dep

## [1.1.0] - 2018-12-25
### Added
- Basic portfolio functionality
- `P` keyboard shortcut to toggle portfolio view
- `e` keyboard shortcut to edit portfolio holdings
- `[portfolio]` TOML config holdings list
