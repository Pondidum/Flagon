# Flagon

*Query flags on the command line*

## Why?

The main reason for developing this is so that I can use feature flags in the CI processes; to allow migration between different steps and ways of doing different tasks, with control over which users/branches/changes see which versions of a task.

The second reason is for debugging:  Is a particular flag on or off for a given user?  By having this in a CLI tool, this information can be embedded in practically any system.

## Usage

```bash
> flagon state "some-flag-name" --user "${user_id}" --attr "branch=${branch}"
# { "name": "some-flag-name", "defaultValue": false, "value": true }
```

Flagon's exit codes are as follows:

- `0` the flag queried is on (`true`)
- `1` the flag queried is off (`false`)
- `2` an error occurred querying the flag

If you need `flagon state`` to always succeed, use `|| true`:

```bash
json=$(flagon state "some-flag-name" --user "${user_id}" --attr "branch=${branch}" || true)
```

By default, the output is the json of the [flag struct](./backends/backend.go#10).  You can also use `--output template=<GO TEMPLATE>` to customise the output, which is useful when exporting the status as environment variables (or outputs) in CI systems:

```bash
> flagon state "some-flag-name" --output "template={{ .Value }}" || true
# true
```

In CI systems, it is often useful to control flags based on the committer, or the branch they are pushing.  For example, this can be done by querying git:

```bash
email=$(git show --quiet --pretty="format:%ce")
branch=$(git rev-parse --abbrev-ref HEAD)

if flagon state "ci-replacement-deploy" --user "${email}" --attr "branch=${branch}" --silent; then
  ./build/replacement-deploy.sh "${branch}" "${commit}"
else
  ./build/deploy.sh "${branch}" "${commit}"
fi
```


The `--user` flag should always map to the identifier for the user in the flag backend, for example, in LaunchDarkly, this maps to the user's `key` property:

```bash
flagon state "ci-replacement-deploy" --user "${user_id}"
```

You can also pass in other attributes:

```bash
flagon state "ci-replacement-deploy" --user "${user_id}" --attr "email=${email}"
```


## Github Actions

Add `pondidum/flagon` as a step in your job, and the `flagon` binary will be available on your `$PATH` in subsequent steps:

```yaml
steps:
- name: Configure Flagon
  uses: pondidum/flagon@main

- name: Run Script
  run: |
    if flagon state "enable-script" false --attr branch="${{ github.ref_name }}"; then
      ./some-script.sh
    fi
```

Optionally, you can specify which version of flagon to use, otherwise the latest release is used:

```yaml
steps:
- name: Configure Flagon
  uses: pondidum/flagon@main
  with:
    version: 0.0.5
```

You can also use flagon to control if jobs (or steps) run.  This uses the `--output template` feature to print just the value of the flag's `.Value` property:

```yaml
jobs:
  flags:
    runs-on: ubuntu-latest

    outputs:
      enabled: ${{ steps.query.outputs.enabled }}

    steps:
    - name: Configure Flagon
      uses: pondidum/flagon@main

    - name: Query
      id: query
      run: echo "enabled=$(flagon state "enable-extra-job" false --output "template={{ .Value }}")" >> "${GITHUB_OUTPUT}"

  controlled:
    runs-on: ubuntu-latest
    if: ${{ needs.flags.outputs.enabled }}
    needs:
      - flags

    steps:
    - name: Print
      run: echo "${{ needs.flags.outputs.enabled }}"
```

## Backends

Currently, this only supports [LaunchDarkly] as a backend.  I am open to Pull Requests or suggestions of other backends to add.

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


[LaunchDarkly]: https://launchdarkly.com
