"""
  Create endpoints file
"""
from datetime import datetime
from datetime import timedelta
import json
import boto3
import requests

l_client = boto3.client('lambda')


class test():
    def profile(self, usr_id):
        self.usr_id = usr_id
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

        res = l_client.invoke(
            FunctionName='collaction-dev-edreinoso-ProfileCRUDFunction-gLFLlRBye4eA',
            InvocationType='RequestResponse',
            Payload=bytes(json.dumps(usr_payload), encoding='utf8'),
        )
        return res

    def crowdaction(self, category, subcategory):
        self.category = category
        self.subcategory = subcategory

        gmt_time = datetime.now() - timedelta(hours=2)  # converting time to GMT
        crowdacticon_date_end = gmt_time + timedelta(minutes=3)
        crowdacticon_date_limit = gmt_time + timedelta(minutes=4)

        crowdaction_start_time = gmt_time.replace(microsecond=0)
        crowdaction_expiry_time = crowdacticon_date_end.replace(microsecond=0)
        crowdaction_join_limit = crowdacticon_date_limit.replace(microsecond=0)

        print(crowdaction_expiry_time, crowdaction_join_limit)

        cwr_payload = {
            "title": "querico",
            "description": "test1",
            "category": category,
            "subcategory": subcategory,
            "location": "test1",
            "date_end": str(crowdaction_expiry_time),
            "date_start": str(crowdaction_start_time),
            "date_limit_join": "2022-06-19",
            # "date_limit_join": str(crowdaction_join_limit), # this should be tested
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
                    "label": "5/7 days a week",
                    "points": 0
                },
                {
                    "description": "",
                    "id": "vegan",
                    "label": "Vegan",
                    "points": 20,
                    "requires": [
                        {
                            "description": "",
                            "id": "vegetarian",
                            "label": "Vegetarian",
                            "points": 20,
                            "requires": [
                                {
                                    "description": "",
                                    "id": "pescatarian",
                                    "label": "Pescatarian",
                                    "points": 5,
                                    "requires": [
                                        {
                                            "description": "",
                                            "id": "no-beef",
                                            "label": "No Beef",
                                            "points": 5
                                        }
                                    ]
                                }
                            ]
                        },
                        {
                            "description": "",
                            "id": "no-dairy",
                            "label": "No Dairy",
                            "points": 10,
                            "requires": [
                                {
                                    "description": "",
                                    "id": "no-cheese",
                                    "label": "No Cheese",
                                    "points": 10
                                }
                            ]
                        }
                    ]
                }
            ],
        }

        uri = 'https://5y310ujdy1.execute-api.eu-central-1.amazonaws.com/dev/cms/crowdactions'

        res = requests.post(uri, json=cwr_payload)
        crowdaction = res.json()
        cid = crowdaction['data']['crowdactionID']
        return cid

    def participation(self, cid, usr_id, commitment):
        self.cid = cid
        self.usr_id = usr_id
        self.commitment = commitment

        # dynamic commitments
        body = {
            "password": "myEvent-myCompany2021",
            "commitments": commitment
        }

        prt_payload = {
            "body": json.dumps(body),
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

        res = l_client.invoke(
            FunctionName='collaction-dev-edreinoso-ParticipationFunction-zc9MRMJVkIjO',
            InvocationType='RequestResponse',
            Payload=bytes(json.dumps(prt_payload), encoding='utf8'),
        )
        return res
