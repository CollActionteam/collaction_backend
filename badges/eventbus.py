# This function should be responsible for point system calculation
# as well as the event rule deletion from the bus once the crowdaction
# terminates

import json
import logging
import boto3

e_client = boto3.client('events')
l_client = boto3.client('lambda')
d_client = boto3.client('dynamodb')

commit_dict = {}  # this may be global for now


def compute_badge_award(points, reward_list):
    """
     There is an assumption about the order
     of the reward_list. This is taken into
     account a descending order
    """
    if int(points) >= int(reward_list[3]):
        return "Diamond"
    elif int(points) >= int(reward_list[2]):
        return "Golden"
    elif int(points) >= int(reward_list[1]):
        return "Silver"
    elif int(points) >= int(reward_list[0]):
        return "Bronze"
    else:
        return "No reward"


def ddb_get_item(table, crowdaction_id):
    pass


def ddb_query(table, crowdaction_id):
    pass


def ddb_update(table, usr_id, reward, crowdaction_id):
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


def tree_recursion(tree):
    for i in range(0, len(tree)):
        t = tree[i]['M']
        commit_key = t['id']['S']
        commit_points = t['points']['N']
        commit_dict[commit_key] = commit_points
        if 'requires' in t:
            tree_recursion(t['requires']['L'])


def lambda_handler(event, context):
    badge_reward_list = []

    target_name = 'lambda_cron_test_end_date'
    single_table = 'collaction-dev-edreinoso-SingleTable-BAXICTFSQ4WV'
    profile_table = 'collaction-dev-edreinoso-ProfileTable-XQEJJNBK6UUY'
    # crowdaction_id = event['resources'][0].split('/')[1]
    crowdaction_id = 'sustainability#food#185f66fd'
    participant_sk = "prt#act#" + crowdaction_id
    print('Lambda Crontab!', event['resources']
          [0].split('/')[1].replace('_', '#'))

    """
      POINT CALCULATION LOGIC
    """

    # 1. fetch the badge scale for crowdaction ✅
    badge_scale = d_client.get_item(
        TableName=single_table,
        Key={
            'pk': {'S': 'act'},
            'sk': {'S': crowdaction_id}
        }
    )
    tree = badge_scale['Item']['commitment_options']['L']
    for reward in badge_scale['Item']['badges']['L']:
        badge_reward_list.append(reward['N'])
    print(badge_reward_list)

    # 2. restructure the tree to a dictionary ✅
    tree_recursion(tree)
    print(commit_dict)  # verifying the dictionary convertion

    # 3. go through all participants ✅
    participant_list = d_client.query(
        TableName=single_table,
        IndexName='invertedIndex',
        KeyConditionExpression="sk = :sk",
        ExpressionAttributeValues={
            ':sk':  {'S': f'prt#act#{crowdaction_id}'}
        },
    )

    # 4. validate their commitment level ✅
    user_prt_list = []  # list required to store individual participations
    for i in range(0, len(participant_list['Items'])):
        prt_details = participant_list['Items'][i]
        usr_id = prt_details['userID']['S']
        prt_lvl = prt_details['commitments']['L'][0]['S']
        usr_obj = {
            "userid": usr_id,
            "prt": prt_lvl,
            "points": commit_dict[prt_lvl]
        }
        if prt_lvl in commit_dict:
            usr_obj['badge'] = compute_badge_award(
                commit_dict[prt_lvl], badge_reward_list)
        user_prt_list.append(usr_obj)

    # print(user_prt_list)

    # 5. award badge ⏰
    for usr in user_prt_list:
        ddb_update(profile_table, usr['userid'], usr['badge'], crowdaction_id)

    return {
        'statusCode': 200,
        'body': json.dumps('Crowdaction Ended!')
    }
