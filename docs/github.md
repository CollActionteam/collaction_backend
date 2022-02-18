# GitHub related documentation
Documentation relating to various GitHub features used by this project.  
(If you are developing locally you probably do not need to pay attention to this)

## Workflow jobs
* `build` - builds the SAM project and caches the result
* `test` - runs GoLang unit tests
* `deploy` - deploys to _prod_ from branch "master" and to _dev_ from branch "development"

## Repository secrets
The following repository secrets are used by GitHub workflows:

### AWS
| Secret | Explanation
|---|---|
|`AWS_ACCESS_KEY_ID`|The AWS access key for the CI user.|
|`AWS_SAMCONFIG_DEV`|The base64 encoded `samconfig.toml` used for deployments to _dev_|
|`AWS_SAMCONFIG_PROD`|The base64 encoded `samconfig.toml` used for deployments to _prod_|
|`AWS_SECRET_ACCESS_KEY`|The AWS secret access key for the CICD user.|

