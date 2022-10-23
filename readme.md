# Flagon

*Query flags on the command line*

## Usage

```
> flagon state "some-flag-name" --user "$user_id" --prop "branch=$branch" --output json
# { "name": "some-flag-name", default: "off", state: "on" }
```