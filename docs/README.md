# CollAction backend
Built using [AWS SAM](https://github.com/aws/serverless-application-model).

An interactive documentation of the API can be found [here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/CollActionteam/collaction_backend/development/docs/api.yaml).

## Dependencies
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Go installed](https://golang.org/doc/install)
* Docker - [Install Docker community edition](https://hub.docker.com/search/?type=edition&offering=community)

## Project structure
- `auth/` - Package for extracting user information from requests
- `email-contact/` - Contact form Lambda function
- `hello-world/` - Example Lambda function.
- `events/` - Invocation events that you can use to invoke the function.
- `process-profile-picture/` - Lambda function for moderating uploaded profile pictures and making them public. 
- `upload-profile-picture/` - Lambda function for generating an upload URL for the profile picture. 
- `whoami/` - Example Lambda for authenticating Firebase users. 
- `profile-crud/` - Lambda function for creating, reading, and updating profile information
- `template.yaml` - A template that defines the application's AWS resources.

For additional information, refer to `README.md` files in these directories.

## Run locally
Build and run the entire application using the following commands:
```bash
sam build
sam local start-api --parameter-overrides
```
In order to connect from another device append `--host 0.0.0.0` to the second command.

You can also run a single function using an event file.
```bash
sam local invoke HelloWorldFunction --event events/event.json
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
