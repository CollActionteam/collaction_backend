# User's profile endpoint

## Create  
### Usage
```bash
curl --location --request POST 'http://domain/profiles/{userID}' \
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

### Response
```
{
    "message":"Profile Created",
    "data":"",
    "status":200
}
```

## Update
### Usage
```bash
curl --location --request PATCH 'http://domain/profiles/{userID}' \
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

### Response
```
{
    "message":"profile update successful",
    "data":"",
    "status":200
}
```



## Get
### Usage
```bash
curl --location --request GET 'http://domain/profiles/' \
--header 'Authorization: Bearer my-jwt-here'
```

or

```bash
curl --location --request GET 'http://domain/profiles/{userID}' \
--header 'Authorization: Bearer my-jwt-here'
```

### Response
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