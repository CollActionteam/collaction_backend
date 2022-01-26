# CollAction backend
Built using [AWS SAM](https://github.com/aws/serverless-application-model).

An interactive documentation of the API can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/CollActionteam/collaction_backend/development/docs/api.yaml).

## Dependencies
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Go installed](https://golang.org/doc/install)
* Docker - [Install Docker community edition](https://hub.docker.com/search/?type=edition&offering=community)

## Architecture
The project will use a Hexagonal/Port-Adapter architecture.

## Run locally
Build and run the entire application using the following commands:
```bash
sam build
sam local start-api --parameter-overrides
```
In order to connect from another device append `--host 0.0.0.0` to the second command.

You can also run a single function using an event file.
```bash
sam local invoke SomeFunction --event event_examples/some_event.json
```

## Unit tests
Run the tests from the root directory of the project thus:
```bash
go test ./...
```
or
```
go test ./... -v
```
