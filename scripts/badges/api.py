"""
    Script for calling API gateway endpoints
    if the API is not available, then try to use
    boto3 to put records in DynamoDB
"""
from datetime import datetime
from datetime import timedelta
import json
import string
import random
import boto3
import requests


l_client = boto3.client('lambda')


def id_generator(size=6, chars=string.ascii_letters + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))


"""
    CREATE USER
    put user in profile dynamodb
    randomID for the user
        if the API is not working,
        the work around is to use boto3
        or event the other work would be
        to use an event directly
    the idea is to create couple of users,
    save them in a list to use later in the
    participation portion of the script
"""

"""
    for loop
        post - multiple users
        save them into a list
            along with their commitment
"""
commitment_arr = [
    "no-diary",
    "no-cheese",
    "no-beef",
    "pescatarian",
    "vegetarian",
    "vegan"
]

randomNum = random.randrange(4)  # this should be inside for loop
# usr_list = [
#     {"id": usr_id, "commitment": commitment_arr[randomNum]}
# ]

# randomly create a user_id
usr_id = id_generator(28)
usr_payload = {
    "body": "{\"displayname\": \"Timothy\", \"city\": \"New York\", \"country\": \"USA\", \"bio\": \"Hi, I'm Timothy\"}",
    "pathParameters": {
        "userID": usr_id
    },
    "requestContext": {
        "authorizer": {
            "jwt": {
                "claims": {
                    "name": "Timothy",
                    "user_id": usr_id,
                    "sub": usr_id,
                    "phone_number": "+31612345678"
                }
            }
        },
        "http": {
            "method": "POST"
        }
    }
}

# uri = 'https://5y310ujdy1.execute-api.eu-central-1.amazonaws.com/dev/cms/crowdactions'
# user = requests.post(uri, json=usr_payload)
# print(user.json())
res = l_client.invoke(
    FunctionName='collaction-dev-edreinoso-ProfileCRUDFunction-gLFLlRBye4eA',
    InvocationType='RequestResponse',
    Payload=bytes(json.dumps(usr_payload), encoding='utf8'),
)
print(res)

"""
    CREATE CROWDACTION
"""
current_time = datetime.now()

gmt_time = datetime.now() - timedelta(hours=2)  # converting time to GMT
crowdacticon_date_end = gmt_time + timedelta(minutes=7)
crowdacticon_date_limit = gmt_time + timedelta(minutes=4)
# datetime.now().strftime("%Y-%m-%d")
crowdaction_start_time = gmt_time.replace(microsecond=0)
crowdaction_expiry_time = crowdacticon_date_end.replace(microsecond=0)
crowdaction_join_limit = crowdacticon_date_limit.replace(microsecond=0)

category = id_generator(8)
subcategory = id_generator(8)

cwr_payload = {
    "title": "querico",
    "description": "test1",
    "category": category,
    "subcategory": subcategory,
    "location": "test1",
    "date_end": str(crowdaction_expiry_time),
    "date_start": str(crowdaction_start_time),
    "date_limit_join": "2022-06-19",
    # "date_limit_join": str(crowdaction_join_limit),
    "password_join": "",
    "images": {
        "card": "hello",
        "banner": "world"
    },
    "badges": [20, 40, 60, 80],
    "participation_count": 0,
    "top_participants": [],
    "commitment_options": [
        {
            "description": "(in case you dont want to commit to 7/7 days a week)",
            "id": "working-days-only",
            "label": "5/7 days a week"
        },
        {
            "description": "",
            "id": "vegan",
            "label": "Vegan",
            "requires": [
                {
                    "description": "",
                    "id": "vegetarian",
                    "label": "Vegetarian",
                    "requires": [
                        {
                            "description": "",
                            "id": "pescatarian",
                            "label": "Pescatarian",
                            "requires": [
                                {
                                    "description": "",
                                    "id": "no-beef",
                                    "label": "No Beef"
                                }
                            ]
                        }
                    ]
                },
                {
                    "description": "",
                    "id": "no-dairy",
                    "label": "No Dairy",
                    "requires": [
                        {
                            "description": "",
                            "id": "no-cheese",
                            "label": "No Cheese"
                        }
                    ]
                }
            ]
        }
    ],
}

uri = 'https://5y310ujdy1.execute-api.eu-central-1.amazonaws.com/dev/cms/crowdactions'

# # print(cwr_payload)

res = requests.post(uri, json=cwr_payload)
crowdaction = res.json()
cid = crowdaction['data']['crowdactionID']
print(cid)

"""
    # CREATE PARTICIPATION
    take the newly created user ID
    take the newly created crowdaction ID
        I would like to be returing the crowdaction itself
        not necessary if I end up using the 
            category and 
            subcategory
        variables
    this could be a for loop that iterates
    through a stack of users (with their ID)
"""

# TODO: complete this loop
"""
    for n in list of users
        from list
            grab the id
            grab the commimtment
"""

prt_payload = {
    # it would be nice that the commitment is randomly generated based on certain
    # options
    "body": "{\"password\":\"myEvent-myCompany2021\", \"commitments\":[\"vegan\"]}",
    "pathParameters": {
        "crowdactionID": cid
    },
    "requestContext": {
        "authorizer": {
            "jwt": {
                "claims": {
                    "name": "Hello World",
                    "user_id": usr_id,
                    "sub": usr_id,
                    "phone_number": "+31612345678"
                }
            }
        },
        "http": {
            "method": "POST"
        }
    }
}

# uri = f'https://5y310ujdy1.execute-api.eu-central-1.amazonaws.com/dev/crowdactions/{crowdactionId}/participation'
# # for n in range(0,10):
# participation = requests.post(uri, json=prt_payload)
# print(participation.json())

res = l_client.invoke(
    FunctionName='collaction-dev-edreinoso-ParticipationFunction-zc9MRMJVkIjO',
    InvocationType='RequestResponse',
    Payload=bytes(json.dumps(prt_payload), encoding='utf8'),
)

print(res)
