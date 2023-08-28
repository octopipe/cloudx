# Cloudx

The Cloudx helps you to orchestrate applications and infrastructure no matter where they run.

## What is Cloudx?

The Cloudx manage your infrastructure in any cloud based in [IDP principles](https://internaldeveloperplatform.org/).

## Development

### Dependencies
- go >= 1.17
- jq

### Controller
Start the controller with the command:
```
$ make controller
```
### Runner
Start the runner locally, for this step you need to know the shared-infra-name and execution-id. You can get this information in the shared-infra applied.
```
$ sh start-runner.sh <shared-infra-name> <execution-id> <ACTION>
```
