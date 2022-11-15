# Changelog

## [0.0.3] - 2022-11-15

### Fixed

- Fix the build's upload binary call, which was truncating the binary file

### Changed

- Statically link the binary in github actions

## [0.0.2] - 2022-11-14

### Fixed

- Releases in Github now publish the binary too

## [0.0.1] - 2022-11-14

### Added

- Exit with code `0` if a flag is `true`, and `1` otherwise
- Add `--silent` flag, to suppress console information

### Changed

- Expand what information is written to traces

## [0.0.0] - 2022-11-11

### Added

- Initial Version
- Read a flag from LaunchDarkly
