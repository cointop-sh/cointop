# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
