# Upload profile picture endpoint
Generate a temporary URL for uploading the profile picture.

## Usage:
```bash
curl --location --request GET 'http://domain/upload-profile-picture' \
--header 'Authorization: Bearer my-jwt-here'
```

## Response
```
{
    "upload_url": "s3 upload url"
}
```

## Using the upload URL
```bash
curl --location --request PUT 'https://upload/url/' \
--form 'image=@"/Path/to/image"'
```