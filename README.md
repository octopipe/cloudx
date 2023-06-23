# runtime-engine

The runtime-engine helps you to orchestrate applications and infrastructure no matter where they run.

## What is runtime-engine?

The runtime-engine manage your infrastructure in any cloud based in [IDP principles](https://internaldeveloperplatform.org/).


### FEATURES:
- Multiple plugin dependencies
- Destroy infrastructure
- Plugins diff
- Terraform execution per version
- No needed of terraform code customization
- Each execution with only credentials of target
- Can watch progress of partial executions
- Historical execution data
- Can re-run with new plugins versions
- Fix the problem of plugin having the files with same names
- When possible parallel execution is performed

### TODO: 
- Delete plugins using dependencies tree
- Add author per execution
- Add execution sensitives data in secrets

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
