# GitHub related documentation
Documentation relating to various GitHub features used by this project.  
(If you are developing locally you probably do not need to pay attention to this)

## Workflows
_(TODO)_

## Repository secrets
The following repository secrets are used by GitHub workflows:

### AWSs
| Secret | Explanation
|---|---|
|`AWS_ACCESS_KEY_ID`|The AWS access key for the CI user.|
|`AWS_SAMCONFIG_DEV`|The base64 encoded `samconfig.toml` used for deployments to `dev`|
|`AWS_SECRET_ACCESS_KEY`|The AWS secret access key for the CI user.|

