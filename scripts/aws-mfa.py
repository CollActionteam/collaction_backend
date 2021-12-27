#!/usr/bin/python3

import json
from pathlib import Path
from subprocess import PIPE, Popen
import sys

"""
Before using this script:
 - Configure AWS cli to output JSON (https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-output-format.html)
 - Rename the default profile in your AWS credentials file to "default_no_mfa"
 - Add the key value pair "iam_user_arn" as the first line of that profile (do not put a comment in that line)
"""

base_profile_name = 'default_no_mfa'
credentials_file_path = Path.home() / '.aws' / 'credentials'
credentials_file_lines = []
base_profile_line_index = -1
with open(credentials_file_path, 'r') as file:
    line = file.readline()
    line_index = 0
    while len(line) > 0:
        if line.strip() == '[default]':
            line = file.readline()
            while not line.strip().startswith('['):
                line = file.readline()
        credentials_file_lines.append(line)
        if line.strip().startswith('[' + base_profile_name):
            base_profile_line_index = line_index
        line = file.readline()
        line_index += 1
if base_profile_line_index == -1:
    print(f'Error: Profile `{base_profile_name}\' not found', file=sys.stderr)
    exit(-1)

k, iam_user_arn = [x.strip() for x in credentials_file_lines[base_profile_line_index + 1].split('=')]
if k != 'iam_user_arn':
    print(f'Error: First key of `{base_profile_name}\' must be `iam_user_arn\'', file=sys.stderr)
    exit(-1)

while True:
    try:
        code = input('Enter MFA code: ').strip()
    except:
        exit()
    out, err = Popen(['aws', 'sts', 'get-session-token',
                      '--duration-seconds', '129600',
                      '--profile', base_profile_name,
                      '--serial-number', iam_user_arn,
                      '--token-code', code
                      ], stdout=PIPE, stderr=PIPE).communicate()
    if len(err) > 0:
        print(
            f'Error: `{err.decode(sys.stderr.encoding).strip()}\'\n', file=sys.stderr)
    else:
        try:
            credentials = json.loads(out.decode(sys.stdout.encoding).strip())[
                'Credentials']
        except:
            print('Error: Please configure aws to output JSON', file=sys.stderr)
            exit(1)
        if credentials_file_lines[-1].strip() != '':
            credentials_file_lines += '\n'
        credentials_file_lines += [l + '\n' for l in [
            '[default]',
            f'aws_access_key_id = {credentials["AccessKeyId"]}',
            f'aws_secret_access_key = {credentials["SecretAccessKey"]}',
            f'aws_session_token = {credentials["SessionToken"]}'
        ]]
        with open(credentials_file_path, 'w') as file:
            file.writelines(credentials_file_lines)
        print('Updated default profile')
        exit()
