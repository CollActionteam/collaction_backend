import boto3

client = boto3.client('events')

response = client.list_rules()

for rules in response['Rules']:
    client.delete_rule(
        Name=rules['Name']
    )
