# This function should be responsible for point system calculation
# as well as the event rule deletion from the bus once the crowdaction
# terminates

import json
import logging
import boto3

e_client = boto3.client('events')
l_client = boto3.client('lambda')
d_client = boto3.client('dynamodb')


def compute_badge_award(points, reward_list):
    print('point system: ', points, reward_list)
    """
     There is an assumption about the order
     of the reward_list. This is taken into
     account a descending order
    """
    if points >= int(reward_list[3]):
        return "Diamond"
    elif points >= int(reward_list[2]):
        return "Golden"
    elif points >= int(reward_list[1]):
        return "Silver"
    elif points >= int(reward_list[0]):
        return "Bronze"
    else:
        return "No reward"


def ddb_query(table, usr_id, reward, crowdaction_id):
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
                'Action': 'ADD'
            }
        },
    )


def lambda_handler(event, context):
    prt_commit = {
        "no-beef": 5,
        "pescatarian": 10,
        "no-dairy": 20,
        "vegetarian": 30,
        "vegan": 70
    }
    badge_reward_list = []

    target_name = 'lambda_cron_test_end_date'
    single_table = 'collaction-dev-edreinoso-SingleTable-BAXICTFSQ4WV'
    profile_table = 'collaction-dev-edreinoso-ProfileTable-XQEJJNBK6UUY'
    # crowdaction_id = event['resources'][0].split('/')[1]
    crowdaction_id = 'sustainability#food#185f66fd'
    participant_sk = "prt#act#" + crowdaction_id
    print('Lambda Crontab!', event['resources']
          [0].split('/')[1].replace('_', '#'))

    # POINT CALCULATION LOGIC

    # 2. fetch the badge scale for crowdaction ✅
    badge_scale = d_client.get_item(
        TableName=single_table,
        Key={
            'pk': {'S': 'act'},
            'sk': {'S': crowdaction_id}
        }
    )
    # print(badge_scale)
    print(badge_scale['Item']['commitment_options'])
    # for reward in badge_scale['Item']['badges']['L']:
    #     badge_reward_list.append(reward['N'])
    # print(badge_reward_list)

    # 3. restructure the tree to a dictionary
    # 4. go through all participants
    # 5. validate their commitment level
    # 6. award badge

    return {
        'statusCode': 200,
        'body': json.dumps('Crowdaction Ended!')
    }
