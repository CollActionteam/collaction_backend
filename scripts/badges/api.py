"""
    Script for calling API gateway endpoints
    if the API is not available, then try to use
    boto3 to put records in DynamoDB
"""
from create import test
import string
import random
import json

create = test()  # init test class


def id_generator(size=6, chars=string.ascii_letters + string.digits):
    return ''.join(random.choice(chars) for _ in range(size))


# Goal!
# 5 users
# 1 crowdaction
# 5 different participations
def main():
    commitment_arr = [
        ["no-cheese"],
        ["no-dairy", "no-cheese"],
        ["no-beef"],
        ["pescatarian", "no-beef"],
        ["vegetarian", "pescatarian", "no-beef"],
        ["vegan", "vegetarian", "pescatarian", "no-beef"],
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

    print()

    """create crowdaction"""
    category = id_generator(8)
    subcategory = id_generator(8)
    cid = create.crowdaction(category, subcategory)

    print()

    """create participation"""
    for n in range(0, len(usr_list)):
        res = create.participation(
            cid, usr_list[n]['id'], usr_list[n]['commitment'])

        print('user id:', usr_list[n]['id'], 'commitment:',
              usr_list[n]['commitment'], 'res:', res)


if __name__ == '__main__':
    main()
