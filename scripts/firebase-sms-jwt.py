#!/usr/bin/python3

import argparse, requests

url_request_code = (
    "https://identitytoolkit.googleapis.com/v1/accounts:sendVerificationCode"
)
url_sign_in_with_phone_number = (
    "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPhoneNumber"
)

argparser = argparse.ArgumentParser()
argparser.add_argument(
    "-p --phone-number",
    dest="phone_number",
    required=True,
    type=str,
    default=None,
    help="International phone number",
)
argparser.add_argument(
    "-c --code", dest="code", required=False, type=str, default=None, help="SMS OTP"
)
argparser.add_argument(
    "-k --key",
    dest="key",
    required=True,
    type=str,
    default=None,
    help="Firebase Web API key",
)

args = argparser.parse_args()

if args.code is None:
    print("Enter the SMS OTP: ", end="")
    args.code = input()

res = requests.post(
    f"{url_request_code}?key={args.key}", data={"phone_number": args.phone_number}
).json()
if "error" in res:
    print(res["error"])
    exit(1)

session_info = res["sessionInfo"]

res = requests.post(
    f"{url_sign_in_with_phone_number}?key={args.key}",
    data={
        "phone_number": args.phone_number,
        "sessionInfo": session_info,
        "code": args.code,
    },
).json()
if "error" in res:
    print(res["error"])
    exit(1)

print(res["idToken"])
