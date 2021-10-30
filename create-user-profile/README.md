# Create user's profile endpoint

## Usage:
```bash
curl --location --request POST 'http://domain/profile' \
--header 'Authorization: Bearer my-jwt-here' \
--header 'Content-Type: application/json' \
--data-raw '
{
"displayname":"displayname",
"country":"country",
"city":"city",
"bio":"bio"
}'
```

## Response
```
{
    "message":"Profile Created",
    "data":"",
    "status":200
}
```