# JWT test endpoint
Echo claims.

The response will contain the claims from the valid [JWT](https://jwt.io/) that was used on the request.

## Usage:
```bash
curl GET http://path/to/my/endpoint \
   -H 'Authorization: Bearer my-jwt-here'
```

## JWT claims
Refer to [this iana document](https://www.iana.org/assignments/jwt/jwt.xhtml#claims) for a list of standard claims to provide in a JWT.
