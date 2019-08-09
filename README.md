# OKPOCK

[![Build Status](https://travis-ci.org/danikarik/okpock.svg?branch=master)](https://travis-ci.org/danikarik/okpock)

## PassKit Endpoints

### GET `/v1/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}`

Response Codes

- `200`
- `204`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "lastUpdated": "2019-08-09 18:34:17",
  "serialNumbers": [
    "02f9ce28-96f5-4e8f-bcb8-d37e7d1e956f"
  ]
}
```

### POST `/v1/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}/{serialNumber}`

Response Codes

- `200`
- `201`
- `400`
- `500`

### DELETE `/v1/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}/{serialNumber}`

Response Codes

- `200`
- `404`
- `500`

### GET `/v1/passes/{passTypeIdentifier}/{serialNumber}`

Response Codes

- `200`
- `304`
- `500`

Response Headers

- `Content-Type - application/vnd.apple.pkpass`

### POST `/v1/log`

Response Codes

- `200`
- `400`

## Public API Endpoints

### GET `/health`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "status": "200"
}
```

### GET `/version`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "version": "6423692"
}
```

### POST `/login`

Request Body

```json
{
  "username": "danikarik",
  "password": "qwerty123"
}
```

Response Codes

- `200`
- `400`
- `404`
- `500`

Response Headers

- `Content-Type - application/json`
- `X-CSRF-Token - <token>`

Response Body

```json
{
  "lastSignInAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### DELETE `/logout`

Response Codes

- `200`

### POST `/register`

Request Body

```json
{
  "username": "danikarik",
  "email": "baitursynov92@gmail.com",
  "password": "qwerty123"
}
```

Response Codes

- `200`
- `400`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "email": "baitursynov92@gmail.com",
  "messageId": "0100016c53867688-e2c27f76-8bed-4b6b-99be-824fbf6cbc20-000000",
  "sentAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### POST `/recover`

Request Body

```json
{
  "email": "baitursynov92@gmail.com"
}
```

Response Codes

- `200`
- `400`
- `404`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "email": "baitursynov92@gmail.com",
  "messageId": "0100016c53867688-e2c27f76-8bed-4b6b-99be-824fbf6cbc20-000000",
  "sentAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### POST `/reset`

Types

- `invite`
- `recovery`

Request Body

```json
{
  "type": "recovery",
  "token": "<token>",
  "password": "qwertt123"
}
```

Response Codes

- `202`
- `400`
- `404`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "confirmedAt": "2019-08-05T23:27:28.981648+06:00"
}
```

```json
{
  "updatedAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### GET `/verify`

Types

- `register`
- `invite`
- `recovery`
- `email_change`

Request query parameters

- `type`
- `token`
- `redirect_url`

Response Codes

- `301`

## Protected API Endpoints

Requirements

- `Cookie`
- `X-CSRF-Token`

### POST `/invite`

Request Body

```json
{
  "email": "baitursynov92@gmail.com"
}
```

Response Codes

- `201`
- `400`
- `401`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "email": "baitursynov92@gmail.com",
  "messageId": "0100016c53867688-e2c27f76-8bed-4b6b-99be-824fbf6cbc20-000000",
  "sentAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### GET `/account`

Response Codes

- `200`
- `401`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "id": 1,
  "role": "client",
  "username": "danikarik",
  "email": "baitursynov92@gmail.com",
  "confirmedAt": "2019-08-06T11:57:10Z",
  "lastSignInAt": "2019-08-09T11:42:57Z",
  "userMetadata": {},
  "createdAt": "2019-08-06T11:57:10Z",
  "updatedAt": "2019-08-06T11:57:10Z"
}
```

### PUT `/account/email`

Request Body

```json
{
  "email": "baitursynov.ds@gmail.com"
}
```

Response Codes

- `200`
- `400`
- `401`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "email": "baitursynov92@gmail.com",
  "messageId": "0100016c53867688-e2c27f76-8bed-4b6b-99be-824fbf6cbc20-000000",
  "sentAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### PUT `/account/username`

Request Body

```json
{
  "username": "daniyar"
}
```

Response Codes

- `200`
- `400`
- `401`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "username": "daniyar",
  "updatedAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### PUT `/account/password`

Request Body

```json
{
  "password": "asdzxc456"
}
```

Response Codes

- `200`
- `400`
- `401`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "updatedAt": "2019-08-05T23:27:28.981648+06:00"
}
```

### PUT `/account/metadata`

```json
{
  "data": {}
}
```

Response Codes

- `200`
- `400`
- `401`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "userMetaData": {},
  "updatedAt": "2019-08-05T23:27:28.981648+06:00"
}
```

## Author

[@danikarik](https://github.com/danikarik)
