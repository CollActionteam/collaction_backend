# Email contact endpoint
Start an email conversation.

An email will be sent to the central inbox. The `reply-to` header is set accordingly to allow responding to the user seeking contact by simply replying to the email.

## Parameters:
Application parameters and schema are defined by global variables in `./app.py`.

The internal email address is provided by [AWS SSM Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html).  
A corresponding parameter must be created for each deployment stage.

## Usage:
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
