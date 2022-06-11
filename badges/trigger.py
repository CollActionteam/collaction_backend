"""
 This function should be responsible for keeping track of the end date
 of the crowdaction, to record it in the event bus for later point
 system calculation
"""

from datetime import datetime
import logging
import boto3
import json
import random
import string

e_client = boto3.client('events')
l_client = boto3.client('lambda')


def lambda_handler(event, context):
    if(event['Records'][0]['eventName'] == 'INSERT'):
        record = event['Records'][0]['dynamodb']['NewImage']

        # TODO: test this!
        # checking for records that are only for crowdactions
        if (record['pk']['S'] != "act"):
            return {
                'statusCode': 200,
                'body': json.dumps('Event is not crowdaction')
            }

        target_arn = 'arn:aws:lambda:eu-central-1:156764677614:function:lambda_cron_test_end_date'
        target_name = 'lambda_cron_test_end_date'
        action = 'lambda:InvokeFunction'

        title = record['title']['S']
        description = record['description']['S']
        date_end = record['date_end']['S']
        crowdaction_id = record['sk']['S'].replace('#', '_')
        # may have to handle date exception
        date_end_expr = datetime_to_cron(
            datetime.strptime(date_end, "%Y-%m-%d %H:%M:%S"))

        # PUT TARGET
        e_client.put_rule(
            Name=crowdaction_id,
            ScheduleExpression=date_end_expr,
            State='ENABLED',
            Description='event rule for ' + crowdaction_id,
        )

        # PUT TARGET
        e_client.put_targets(
            Rule=crowdaction_id,
            Targets=[
                {
                    'Id': crowdaction_id,
                    'Arn': target_arn,
                },
            ]
        )

        # ADD PERMISSIONS
        l_client.add_permission(
            FunctionName=target_name,
            StatementId=crowdaction_id,
            Action=action,
            Principal='events.amazonaws.com',
            SourceArn='arn:aws:events:eu-central-1:156764677614:rule/'+crowdaction_id,
        )

        print('rule has been placed on the bus')

    return {
        'statusCode': 200,
        'body': json.dumps('End day has been scheduled!')
    }


"""
  Function that would convert datetime into cronjobs
"""


def datetime_to_cron(dt):
    return f"cron({dt.minute} {dt.hour} {dt.day} {dt.month} ? {dt.year})"
