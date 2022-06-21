"""
    Script for calling API gateway endpoints
    if the API is not available, then try to use
    boto3 to put records in DynamoDB
"""
from create import test
import string
import random


create = test()  # init test class


def id_generator(size=6, chars=string.ascii_letters + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))


# Goal!
# 5 users
# 1 crowdaction
# 5 different participations

commitment_arr = [
    ["no-cheese"],
    ["no-cheese", "no-diary"],
    ["no-beef"],
    ["no-beef", "pescatarian"],
    ["no-beef", "pescatarian", "vegeterian"],
    ["no-beef", "pescatarian", "vegeterian", "vegan"]
]

usr_list = []

"""create users"""
for i in range(0, 5):
    usr_id = id_generator(28)
    randomNum = random.randrange(6)

    usr_obj = {
        "id": usr_id,
        "commitment": commitment_arr[randomNum]
    }

    usr_list.append(usr_obj)

    res = create.profile(usr_id)

    print('user id:', usr_id, 'res:', res)

"""create crowdaction"""
category = id_generator(8)
subcategory = id_generator(8)
cid = create.crowdaction(category, subcategory)

"""create participation"""
for n in range(0, len(usr_list)):
    res = create.participation(cid, usr_list[n]['id'])

    print('res:', res)
