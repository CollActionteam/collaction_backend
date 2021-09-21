import boto3
from cerberus import Validator
import json
from re import match
from string import Template

KEY_EMAIL = 'email'
KEY_SUBJECT = 'subject'
KEY_MESSAGE = 'message'
KEY_APP_VERSION = 'app_version'
PATTERN_EMAIL = '^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$'
PATTERN_APP_VERSION = '^(?:ios|android) [0-9]+\.[0-9]+\.[0-9]+\+[0-9]+$'
SCHEMA = {
    KEY_EMAIL: {'type': 'string', 'required': True, 'regex': PATTERN_EMAIL},
    KEY_SUBJECT: {'type': 'string', 'required': True, 'maxlength': 50},
    KEY_MESSAGE: {'type': 'string', 'required': True, 'maxlength': 500},
    KEY_APP_VERSION: {'type': 'string',
                      'required': True, 'regex': PATTERN_APP_VERSION}
}
CHARSET = 'UTF-8'
SSM_PARAMETER_EMAIL = Template('/collaction/$stage/contact/email')
EMAIL_SENDER = Template('CollAction app user <$email>')
EMAIL_SUBJECT = Template('$subject ($email via CollAction app)')
EMAIL_BODY = Template(
    'email: $email\n\n[message]\n$message\n\napp version: $app_version')

validator = Validator(SCHEMA)
aws_ses = boto3.client('ses')
aws_ssm = boto3.client('ssm')


def send_message(sender, recipient, reply_to, subject, message, app_version):
    aws_ses.send_email(
        Destination={
            'ToAddresses': [
                recipient,
            ]
        },
        ReplyToAddresses=[
            reply_to,
        ],
        Message={
            'Body': {
                'Text': {
                    'Charset': CHARSET,
                    'Data': EMAIL_BODY.substitute(email=reply_to, message=message, app_version=app_version)
                },
            },
            'Subject': {
                'Charset': CHARSET,
                'Data': EMAIL_SUBJECT.substitute(subject=subject, email=reply_to)
            },
        },
        Source=EMAIL_SENDER.substitute(email=sender),
    )


def lambda_handler(event, context):
    stage = event['requestContext']['stage']
    parameter_name = SSM_PARAMETER_EMAIL.substitute(stage=stage)
    internal_email = aws_ssm.get_parameter(
        Name=parameter_name, WithDecryption=True)['Parameter']['Value']
    assert(match(PATTERN_EMAIL, internal_email))

    body_json = {}
    try:
        body_json = json.loads(event['body'])
    except:
        pass
    if not validator(body_json):
        return {
            'statusCode': 400,
            'body': json.dumps({'message': 'Invalid payload'})
        }
    else:
        reply_to = body_json[KEY_EMAIL]
        subject = body_json[KEY_SUBJECT]
        message = body_json[KEY_MESSAGE]
        app_version = body_json[KEY_APP_VERSION]
        send_message(internal_email, internal_email,
                     reply_to, subject, message, app_version)
        return {
            'statusCode': 200,
            'body': json.dumps({'message': f"Message sent"})
        }
