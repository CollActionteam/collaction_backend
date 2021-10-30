# Update user's profile data endpoint

## Usage:
```bash
curl --location --request PATCH 'http://domain/profile' \
--header 'Authorization: Bearer my-jwt-here' \
--header 'Content-Type: application/json' \
--data-raw '
{
"displayname":"displayname",
"country":"country",
"city":"city",
"phone":"phone",
}'
```

## Response
```
{
    "message":"profile update successful",
    "data":"",
    "status":200
}
```