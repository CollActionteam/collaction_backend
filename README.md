# CollAction backend
Built using [AWS SAM](https://github.com/aws/serverless-application-model).

## Dependencies
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Python 3 installed](https://www.python.org/downloads/)
* Docker - [Install Docker community edition](https://hub.docker.com/search/?type=edition&offering=community)

## Project structure
- `email-contact/` - Contact form Lambda function
- `hello_world/` - Example Lambda function.
- `events/` - Invocation events that you can use to invoke the function.
- `tests/` - Unit tests for the application code. 
- `template.yaml` - A template that defines the application's AWS resources.

For additional information, refer to `README.md` files in these directories.

## Run locally
Build and run the entire application using the following commands:
```bash
sam build
sam local start-api
```
In order to connect from another device append `--host 0.0.0.0` to the second command.

You can also run a single function using an event file.
```bash
sam local invoke HelloWorldFunction --event events/event.json
```

## Unit tests
Tests are defined in the `tests` folder in this project. Use PIP to install [pytest](https://pypi.org/project/pytest/) and [pytest-mock](https://pypi.org/project/pytest-mock/) before running the tests.
```bash
python -m pytest tests/
```
