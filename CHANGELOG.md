# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.5] - 2019-06-30
### Fixed
- Fixed table sorting issues

## [1.3.4] - 2019-06-30
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
