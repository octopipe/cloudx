# runtime-engine

The runtime-engine helps you to orchestrate applications and infrastructure no matter where they run.

## What is runtime-engine?

The runtime-engine manage your infrastructure in any cloud based in [IDP principles](https://internaldeveloperplatform.org/).

## How to run?

Apply a shared infra manifest in the cluster, example:
```yaml
apiVersion: commons.cloudx.io/v1alpha1
kind: SharedInfra
metadata:
  name: shared-infra-1
spec:
  author: Author
  description: SharedInfra for example 1
  plugins:
  - name: test-1
    ref: registry/plugin:test-1
    depends: []
    type: terraform
    outputs: []
    inputs:
    - key: name
      value: test-1
  - name: test-2
    ref: registry/plugin:test-1
    type: aws/secrets-manager
    depends: 
    - test-1
    outputs: []
    inputs:
    - key: name
      value: "test-2-{{ test-1.outputs.uid }}"
```


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
