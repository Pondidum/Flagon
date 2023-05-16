# Changelog

## [0.0.9] - 2023-05-16

## Added

- you can now use `pondidum/flagon/query` action to query a flag.  It assumes you have used `pondidum/flagon` to setup the tool first.

## [0.0.8] - 2023-01-27

## Added

- add support for `TRACEPARENT` environment variable; format matches the [w3 format](https://www.w3.org/TR/trace-context-1/).

## [0.0.7] - 2022-11-19

### Changed

- `state` exit codes have been updated: `0` for flag on, `1` for flag off, `2` for error.

## [0.0.6] - 2022-11-18

### Added

- printing supports go templates: `--output "template={{ .Value }}"` for example

## [0.0.5] - 2022-11-17

### Added

- Add tests around `state` command, fix default value handling

## [0.0.4] - 2022-11-15

### Added

- Add github action.yaml to the repository

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
