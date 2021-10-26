# Get s3 upload url


## Usage:
```bash
curl --location --request GET 'http://domain/upload-profile-picture' \
--header 'Authorization: Bearer my-jwt-here'
```

## Response
{
    "upload_url": "s3 upload url"
}

## Using upload url

curl --location --request PUT 'https://upload/url/' \
--form 'image=@"/Path/to/image"'