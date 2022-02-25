# CollAction backend
The backend for the [CollAction app](https://github.com/CollActionteam/collaction_app).

## API
An interactive documentation of the API can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/CollActionteam/collaction_backend/development/docs/api.yaml).

❗ Currently the API is being overhauled (see `./docs/api2.yml`)  
The new version will conform to [JSend](https://github.com/omniti-labs/jsend).

## Dependencies
* [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [GoLang](https://golang.org/doc/install)
* [Docker](https://hub.docker.com/search/?type=edition&offering=community)

## Architecture
The project follows the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  
The file structure is as follows:
```
repository/
├─ docs/             👉 Documentation
├─ internal/         👉 Contains folders for business logic and (unit) tests for each service in a corresponding folder
│  ├─ constants/     👉 (Shared) constant values
│  ├─ models/        👉 Definitions for structs (❗ No logic)
├─ pkg/
│  ├─ handler/
│  |  ├─ aws/        👉 Lambdas (seperate folders each)
│  ├─ mocks/
│  |  ├─ repository/ 👉 Mocks for external repositories
│  ├─ repository/
│  |  ├─ aws/        👉 External repositories (e.g. AWS SSM/Dynamo)
├─ go.mod            👉 Go dependencies
├─ template.yaml     👉 CloudFormation template
```

## Run locally
⚠ Not all features of the API can be run locally!  
(To use the full range of AWS services, deploy a stack for testing using `sam deploy -g`)

Build and run the entire application using the following commands:
```bash
sam build
sam local start-api
```

You can also run a single function using an event file.
```bash
sam local invoke SomeFunction --event event_examples/some_event.json
```

## Unit tests
Run the tests from the root directory using:
```bash
go test ./...
```

## DevOps
GitHub topics (such as GitHub Actions) are documented [here](https://github.com/CollActionteam/collaction_backend/blob/development/docs/github.md)
