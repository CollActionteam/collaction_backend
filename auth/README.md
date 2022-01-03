# JWT claims extraction
Extracts user info struct from jwt claims given an APIGateway request with a valid [JWT](https://jwt.io/).

## JWT claims
Refer to [this iana document](https://www.iana.org/assignments/jwt/jwt.xhtml#claims) for a list of standard claims to provide in a JWT.

The following claims are required:
* `user_id` (same as `sub`)
* `name`
* `phone_number`