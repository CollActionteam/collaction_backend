import boto3
d_client = boto3.client('dynamodb')


class ddb_ops():
    """
      Get Item
    """

    def get_item(self, table, crowdaction_id):
        self.table = table
        self.crowdaction_id = crowdaction_id

        res = d_client.get_item(
            TableName=table,
            Key={
                'pk': {'S': 'act'},
                'sk': {'S': crowdaction_id}
            }
        )
        return res

    """
      Query Items
    """

    def query(self, table, crowdaction_id):
        self.table = table
        self.crowdaction_id = crowdaction_id

        res = d_client.query(
            TableName=table,
            IndexName='invertedIndex',
            KeyConditionExpression="sk = :sk",
            ExpressionAttributeValues={
                ':sk':  {'S': f'prt#act#{crowdaction_id}'}
            },
        )
        return res
    """

      Update Items
    """

    def update(self, table, usr_id, reward, crowdaction_id):
        self.table = table
        self.usr_id = usr_id
        self.reward = reward
        self.crowdaction_id = crowdaction_id

        d_client.update_item(
            TableName=table,
            Key={
                'userid': {
                    'S': usr_id,
                }
            },
            AttributeUpdates={
                'reward': {
                    'Value': {
                        'L': [
                            {
                                "M": {
                                    "award": {
                                        "S": reward
                                    },
                                    "crowdactionID": {
                                        "S": crowdaction_id
                                    }
                                }
                            },
                        ],
                    },
                    'Action': 'ADD'  # this operations still pending
                }
            },
        )
