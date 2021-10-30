# Get user's profile data endpoint

## Usage:
```bash
curl --location --request GET 'http://domain/profile' \
--header 'Authorization: Bearer my-jwt-here'
```

## Response
```
{
    "message": "Successfully Retrieved Profile",
    "data": {
        "userid": "id",
        "displayname": "name",
        "country": "country",
        "city": "city",
        "bio": "bio",
        "phone": "phone number"
    },
    "status": 200
}
```
or 

```
{
    "message": "no user Profile found",
    "data": "",
    "status": 404
}
```