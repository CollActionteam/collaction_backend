#!/usr/bin/python3

'''
Simple command line utility for creating DynamoDB JSON items for crowdations
Prerequisites: pip install PyInquirer
Usage: python3 generate-crowdaction.py
'''

import json
from PyInquirer import prompt, Validator, ValidationError
import re
import tkinter
from tkinter import filedialog
import uuid


class UrlValidator(Validator):
    def validate(self, document):
        if not (document.text.startswith('http://') or document.text.startswith('https://')):
            raise ValidationError(
                message='Must be a http(s) URL', cursor_position=len(document.text))


class DateValidator(Validator):
    def validate(self, document):
        if not re.match(r'\d{4}-\d{2}-\d{2}', document.text):
            raise ValidationError(
                message='Must be a date with format YYYY-mm-dd', cursor_position=len(document.text))


print('Create a new crowdaction:')
answers = prompt([
    {
        'type': 'input',
        'name': 'title',
        'message': 'Title',
    },
    {
        'type': 'input',
        'name': 'category',
        'default': 'sustainability',
        'message': 'Category',
    },
    {
        'type': 'input',
        'name': 'subcategory',
        'default': 'food',
        'message': 'Sub-category (optional)',
    },
    {
        'type': 'editor',
        'name': 'description',
        'message': 'Description',
        'default': '<Replace with crowdaction description>',
        'eargs': {
            'editor': 'default',
            'ext': '.txt'
        }
    }, {
        'type': 'input',
        'name': 'banner',
        'message': 'Banner image URL',
        'validator': UrlValidator()
    },
    {
        'type': 'input',
        'name': 'card',
        'message': 'Card image URL',
        'validator': UrlValidator()
    },
    {
        'type': 'input',
        'name': 'country',
        'default': 'NL',
        'message': 'Country code (optional)',
    }
])

category = answers['category'].strip()
subcategory = answers['subcategory'].strip()
crowdaction_id = f'{category}#{subcategory}#{str(uuid.uuid4())[:8]}'
country = answers['country'].strip()
city = ''
if len(country) > 0:
    city = prompt([{
        'type': 'input',
        'name': 'city',
        'default': 'Amsterdam',
        'message': 'City (optional)',
    }])['city'].strip()


def confirm(question, default=False):
    return prompt([{
        'type': 'confirm',
        'message': question,
        'name': 'response',
        'default': default,
    }])['response']


password = ''
if confirm('Requires password'):
    password = prompt([{
        'type': 'input',
        'message': 'Password',
        'name': 'password',
    }])['password'].strip()


def get_date_str(message):
    return prompt([{
        'type': 'input',
        'name': 'date',
        'message': message,
        'validator': DateValidator()
    }])['date'].strip()


dto = {
    'pk': {
        'S': 'act'
    },
    'sk': {
        'S': crowdaction_id
    },
    'participant_count': {
        'N': '0'
    },
    'location': {
        'S': f'{country}#{city}'
    },
    'top_participants': {
        'L': []
    },
    'password_join': {
        'S': password
    },
    'commitment_options': {
        'L': []
    },
    'crowdactionID': {
        'S': crowdaction_id
    },
    'date_start': {
        'S': get_date_str('Start date')
    },
    'date_end': {
        'S': get_date_str('End date')
    },
    'date_limit_join': {
        'S': get_date_str('Join deadline date')
    },
    'category': {
        'S': category
    },
    'images': {
        'M': {
            'banner': {
                'S': answers['banner'].strip()
            },
            'card': {
                'S': answers['card'].strip()
            }
        }
    },
    'description': {
        'S': answers['description'].strip()
    },
    'title': {
        'S': answers['title'].strip()
    }
}
if len(subcategory) > 0:
    dto['subcategory'] = {
        'S': subcategory
    }


def create_commitment_option(level=0):
    prefix = '|   ' * level + '├─'
    answers = prompt([
        {
            'type': 'input',
            'name': 'label',
            'message': f'{prefix} Label',
        },
        {
            'type': 'input',
            'name': 'id',
            'message': f'{prefix} ID',
        },
        {
            'type': 'input',
            'name': 'description',
            'message': f'{prefix} Description',
        },
    ])
    option = {
        'M': {
            'id': {
                'S': answers['id'].strip()
            },
            'label': {
                'S': answers['label'].strip()
            },
            'description': {
                'S': answers['description'].strip()
            }
        }
    }
    while confirm(f'{prefix} Add requirement'):
        if 'requires' not in option['M']:
            option['M']['requires'] = {
                'L': []
            }
        option['M']['requires']['L'].append(
            create_commitment_option(level=(level+1)))
    return option


print('Add commitment options:')
while len(dto['commitment_options']['L']) < 1 or confirm('Add commitment option'):
    dto['commitment_options']['L'].append(create_commitment_option())

print('\nGenerating JSON file...', end='')
json_dto = json.dumps(dto, indent=4, sort_keys=True)
print('Done!')

tkinter.Tk().withdraw()
filetype_options = [('DynamoDB JSON item', '.dynamodb.json')]
with filedialog.asksaveasfile(filetypes=filetype_options, defaultextension=filetype_options) as file:
    file.write(json_dto)
    print(f'\nSaved to `{file.name}\'')
