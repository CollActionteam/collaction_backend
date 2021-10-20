# JWT test endpoint
Echo claims.

The response will contain the claims from the valid [JWT](https://jwt.io/) that was used on the request using the local `auth` package.

## Usage:
```bash
curl GET http://path/to/my/endpoint \
   -H 'Authorization: Bearer my-jwt-here'
```
or using an event
```bash
cat <<EOF >> my-event.json
{
  "requestContext": {
    "authorizer": {
      "jwt": {
        "claims": {
          "name": "Timothy",
          "user_id": "S0m3Rand0mStr1nG",
          "sub": "S0m3Rand0mStr1nG",
          "phone_number": "+31612345678"
        }
      }
    }
  }
}
EOF
sam local invoke WhoamiFunction --event my-event.json
```
or
```bash

```sam local invoke WhoamiFunction --event events/example_whoami.json