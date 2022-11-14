# Flagon

*Query flags on the command line*

## Usage

```
> flagon state "some-flag-name" --user "$user_id" --attr "branch=$branch"
# { "name": "some-flag-name", "defaultValue": false, "value": true }
```

For example, in a CI system where you wish to select which deployment process is used:

```shell
email=$(git show --quiet --pretty="format:%ce")
branch=$(git rev-parse --abbrev-ref HEAD)

if flagon state "ci-replacement-deploy" --user "${email}" --attr "branch=${branch}" --silent; then
  ./build/replacement-deploy.sh "${branch}" "${commit}"
else
  ./build/deploy.sh "${branch}" "${commit}"
fi
```

## Configuration

### Common

| Flag        | Default         | Description                                                                 |
|-------------|-----------------|-----------------------------------------------------------------------------|
| `--backend` | `launchdarkly`  | The backend to query flags from                                             |
| `--output`  | `json`          | The output format to write to the console.  Currently only supports `json`  |
| `--silent`  | `false`         | Silence any console output                                                  |

### Telemetry

| EnvVar                                | Default           | Description                                                     |
|---------------------------------------|-------------------|-----------------------------------------------------------------|
| `OTEL_TRACE_EXPORTER`                 | ` `               | Which exporter to use: `otlp`, `stdout`, `stderr`               |
| `OTEL_EXPORTER_OTLP_ENDPOINT`         | `localhost:4317`  | Set the Exporter endpoint                                       |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`  | ` `               | Set the Exporter endpoint, takes priority over `OTEL_EXPORTER_OTLP_ENDPOINT` |
| `OTEL_EXPORTER_OTLP_HEADERS`          | ` `               | A Csv of Headers and Values to pass to the tracing service, for example `Authentication: Bearer 13213213,X-Environment: Production` |
| `OTEL_DEBUG`                          | `false`           | Print debug information from tracing to the console             |

### Backend: LaunchDarkly

| EnvVar              | Flag           | Default  | Description                                                                  |
|---------------------|----------------|----------|------------------------------------------------------------------------------|
| `FLAGON_LD_SDKKEY`  | `--ld-sdk-key` |          | The [project](https://app.launchdarkly.com/settings/projects) SDK Key to use |
| `FLAGON_LD_TIMEOUT` | `--ld-timeout` | `10s`    | How long to wait for successful connection                                   |
| `FLAGON_LD_DEBUG`   | `--ld-debug`   | `0`      | Set to `true` (or `1`) to see debug information from the LaunchDarkly client |
