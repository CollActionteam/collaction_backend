import json


def lambda_handler(event, context):
    claims = {}
    if 'authorizer' in event['requestContext']:
        claims = event['requestContext']['authorizer']['jwt']['claims']
    return {
        "statusCode": 200,
        "body": json.dumps({
            "message": f'claims from token are {str(claims)}',
        }),
    }
