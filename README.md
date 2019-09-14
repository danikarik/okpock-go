# OKPOCK

## PassKit Endpoints

### GET `/v1/devices/{deviceID}/registrations/{passTypeID}`

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

### POST `/v1/devices/{deviceID}/registrations/{passTypeID}/{serialNumber}`

Response Codes

- `200`
- `201`
- `400`
- `500`

### DELETE `/v1/devices/{deviceID}/registrations/{passTypeID}/{serialNumber}`

Response Codes

- `200`
- `404`
- `500`

### GET `/v1/passes/{passTypeID}/{serialNumber}`

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
- `403`
- `404`
- `423`
- `500`

Response Headers

- `Content-Type - application/json`
- `X-XSRF-Token - <token>`

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

### POST `/check/email`

Request Body

```json
{
  "email": "baitursynov92@gmail.com"
}
```

Response Codes

- `200`
- `403`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "email": "baitursynov92@gmail.com"
}
```

### POST `/check/username`

Request Body

```json
{
  "username": "danikarik"
}
```

Response Codes

- `200`
- `403`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "username": "danikarik"
}
```

## Protected API Endpoints

Requirements

- `Cookie`
- `X-XSRF-Token`

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

### POST `/uploads`

`file` multipart field

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
  "id": 1,
  "uuid": "1/4553abc6-64ba-47e9-80e9-51b214faed4b",
  "filename": "gopher.jpg",
  "hash": "i1hNKTFWBI18JOWn9VcSENiteao3aexiPHCjax4OtZg=",
  "createdAt": "2019-09-01T03:23:14.162087+06:00"
}
```

### GET `/uploads`

Response Codes

- `200`
- `401`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
[
  {
    "id": 1,
    "uuid": "1/4553abc6-64ba-47e9-80e9-51b214faed4b",
    "filename": "gopher.jpg",
    "hash": "i1hNKTFWBI18JOWn9VcSENiteao3aexiPHCjax4OtZg=",
    "createdAt": "2019-09-01T03:23:14.162087+06:00"
  }
]
```

### GET `/uploads/{id}`

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
  "id": 1,
  "uuid": "1/4553abc6-64ba-47e9-80e9-51b214faed4b",
  "filename": "gopher.jpg",
  "hash": "i1hNKTFWBI18JOWn9VcSENiteao3aexiPHCjax4OtZg=",
  "createdAt": "2019-09-01T03:23:14.162087+06:00"
}
```

### GET `/uploads/{id}/file`

Response Codes

- `200`
- `400`
- `401`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

- `binary` stream

### GET `/account/info`

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

### POST `/projects/check`

Request Body

```json
{
  "title": "Friday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon"
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
  "title": "Friday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon"
}
```

### POST `/projects/`

Request Body

```json
{
  "title": "Friday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon"
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
  "id": 27
}
```

### GET `/projects/`

Response Codes

- `200`
- `401`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
[
  {
    "id": 27,
    "title": "Friday Deal",
    "organizationName": "Okpock",
    "description": "Free Coupon",
    "passType": "coupon",
    "backgroundImage": "background.png",
    "footerImage": "footer.png",
    "iconImage": "icon.png",
    "stripImage": "strip.png",
    "createdAt": "2019-08-29T22:37:57+06:00",
    "updatedAt": "2019-08-29T22:37:57+06:00"
  }
]
```

### GET `/projects/{id}`

Response Codes

- `200`
- `400`
- `401`
- `404`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "id": 27,
  "title": "Friday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon",
  "backgroundImage": "background.png",
  "footerImage": "footer.png",
  "iconImage": "icon.png",
  "stripImage": "strip.png",
  "createdAt": "2019-08-29T22:37:57+06:00",
  "updatedAt": "2019-08-29T22:37:57+06:00"
}
```

### PUT `/projects/{id}`

Request Body

```json
{
  "title": "Saturday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
}
```

Response Codes

- `200`
- `400`
- `401`
- `404`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "id": 27,
  "title": "Saturday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon",
  "backgroundImage": "background.png",
  "footerImage": "footer.png",
  "iconImage": "icon.png",
  "stripImage": "strip.png",
  "createdAt": "2019-08-29T22:37:57+06:00",
  "updatedAt": "2019-08-29T23:43:24+06:00"
}
```

### POST `/projects/{id}/upload`

Request Body

```json
{
  "uuid": "1/3165f717-0086-40de-aa01-eab5104c8e0f",
  "type": "background"
}
```

Available types

- `background`
- `footer`
- `icon`
- `logo`
- `strip`

Response Codes

- `200`
- `400`
- `401`
- `404`
- `406`
- `500`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "id": 27,
  "title": "Saturday Deal",
  "organizationName": "Okpock",
  "description": "Free Coupon",
  "passType": "coupon",
  "backgroundImage": "1/3165f717-0086-40de-aa01-eab5104c8e0f",
  "footerImage": "",
  "iconImage": "",
  "stripImage": "",
  "createdAt": "2019-08-29T22:37:57+06:00",
  "updatedAt": "2019-08-29T23:43:24+06:00"
}
```

### GET `/dictionary/passtypes`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "boardingPass",
    "coupon",
    "eventTicket",
    "generic"
    "storeCard"
  ]
}
```

### GET `/dictionary/detectortypes`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKDataDetectorTypePhoneNumber",
    "PKDataDetectorTypeLink",
    "PKDataDetectorTypeAddress",
    "PKDataDetectorTypeCalendarEvent"
  ]
}
```

### GET `/dictionary/textalignment`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKTextAlignmentLeft",
    "PKTextAlignmentCenter",
    "PKTextAlignmentRight",
    "PKTextAlignmentNatural"
  ]
}
```

### GET `/dictionary/datestyle`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKDateStyleNone",
    "PKDateStyleShort",
    "PKDateStyleMedium",
    "PKDateStyleLong",
    "PKDateStyleFull"
  ]
}
```

### GET `/dictionary/numberstyle`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKNumberStyleDecimal",
    "PKNumberStylePercent",
    "PKNumberStyleScientific",
    "PKNumberStyleSpellOut"
  ]
}
```

### GET `/dictionary/transittype`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKTransitTypeAir",
    "PKTransitTypeBoat",
    "PKTransitTypeBus",
    "PKTransitTypeGeneric",
    "PKTransitTypeTrain"
  ]
}
```

### GET `/dictionary/barcodeformat`

Response Codes

- `200`

Response Headers

- `Content-Type - application/json`

Response Body

```json
{
  "data": [
    "PKBarcodeFormatQR",
    "PKBarcodeFormatPDF417",
    "PKBarcodeFormatAztec",
    "PKBarcodeFormatCode128"
  ]
}
```

## Author

[@danikarik](https://github.com/danikarik)
