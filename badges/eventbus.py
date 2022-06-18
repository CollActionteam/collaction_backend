# This function should be responsible for point system calculation
# as well as the event rule deletion from the bus once the crowdaction
# terminates

import json
import logging
import boto3
import os
from dynamodb import *

e_client = boto3.client('events')
l_client = boto3.client('lambda')

ddb = ddb_ops()

commit_dict = {}  # this may be global for now


def compute_badge_award(points, reward_list):
    print('venezuela')
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
    crowdaction_id = event['resources'][0].split(
        '/')[1].replace('_', '#')  # prod
    # crowdaction_id = 'sustainability#food#185f66fd' # test
    print('Lambda Crontab!', crowdaction_id)

    """
      POINT CALCULATION LOGIC
    """

    # 1. fetch the badge scale for crowdaction ✅
    badge_scale = ddb.get_item(single_table, crowdaction_id)
    # print(badge_scale)

    tree = badge_scale['Item']['commitment_options']['L']
    for reward in badge_scale['Item']['badges']['L']:
        badge_reward_list.append(reward['N'])
    print(badge_reward_list)  # verfying the badge reward list

    # 2. restructure the tree to a dictionary ✅
    tree_recursion(tree)
    print(commit_dict)  # verifying the dictionary convertion

    # 3. go through all participants ✅
    participant_list = ddb.query(single_table, crowdaction_id)
    # print(participant_list)

    # 4. map user commitment level ✅
    user_prt_list = []  # list required to store individual participations
    for i in range(0, len(participant_list['Items'])):
        prt_details = participant_list['Items'][i]
        usr_id = prt_details['userID']['S']
        prt_lvl = prt_details['commitments']['L']
        usr_prt_counter = 0
        for n in range(0, len(prt_lvl)):
            usr_prt_counter += int(commit_dict[prt_lvl[n]['S']])
        usr_obj = {
            "userid": usr_id,
            "prt": prt_lvl,
            "points": usr_prt_counter
        }
        print(usr_obj)
        # if prt_lvl in commit_dict: # would I be assuming that a user would always have a participation
        print('helloworld')
        usr_obj['badge'] = compute_badge_award(
            usr_prt_counter, badge_reward_list)
        user_prt_list.append(usr_obj)

    print(user_prt_list)

    # 5. award badge ✅
    for usr in user_prt_list:
        ddb.update(profile_table, usr['userid'], usr['badge'], crowdaction_id)

    # 6. delete event ✅
    crowdaction_id_e = crowdaction_id.replace('#', '_')
    e_client.delete_rule(
        Name=crowdaction_id_e,
    )

    # 7. delete permission ✅
    l_client.remove_permission(
        FunctionName=target_name,
        StatementId=crowdaction_id_e,
    )

    return {
        'statusCode': 200,
        'body': json.dumps('Crowdaction Ended!')
    }
