# Email contact endpoint
Start an email conversation.

An email will be sent to the central inbox. The `reply-to` header is set accordingly to allow responding to the user seeking contact by simply replying to the email.

## Parameters:
Application parameters and schema are defined by global variables in `./app.py`.

The internal email address is provided by [AWS SSM Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html).  
A corresponding parameter must be created for each deployment stage.

## Usage:
Start the local server first
```bash
sam local start-api
```

```bash
curl POST http://path/to/my/endpoint \
   -H 'Content-Type: application/json' \
   -d @- << EOF
{
    "email": "test@example.org",
    "subject": "Hello world",
    "message": "Please respond to this email :)",
    "app_version": "android 1.0.1+1"
}
EOF
```
or
```bash
curl -H "Content-Type: application/json" \
-d "{\"email\": \"test@example.org\", \"subject\": \"Hello world\", \"message\": \"Please respond to this email :)\",\"app_version\": \"android 1.0.1+1\"}" \http://127.0.0.1:3000/contact
```
or (without local server)
```bash
sam local invoke EmailContactFunction --event event_examples/email.json
```