# Invite Endpoint

Send email, and optionally, an calendar invitation with RSVP to the users.

### URL

`/invite/`

### Payload

```
{
    "users": [
        {
            "full_name": "Eric Cartman",
            "email": "ihateyouguys@southpark.cc"
        },
        {
            "full_name": "Tina belcher",
            "email": "aaaooo@bobsburger.com"
        }
    ],
    "invitation": {
        "start_at": "2030-10-12T07:20:50.52Z",
        "end_at": "2030-10-12T08:20:50.52Z",
        "organizer_email": "meetingorganizer@meeting.com",
        "organizer_full_name": "Mr. Mojo Rising",
        "summary": "This meeting will be about...",
        "location": "https://zoom.us/332324342",
        "description": "Voluptatum ut quis ut. Voluptas qui pariatur quo. Omnis enim rerum dolorum. Qui aut est sed qui voluptatem harum. Consequuntur et accusantium culpa est fuga molestiae in ut. Numquam harum"
    },
    "email_subject": "You've just been invited!",
    "email_body": "<html><body><h1>email body about the invitation/event</h1></body></html>",
    "email_is_html": true
}
```
*Notes about fields:*

- If you don't need to send a calendar invitation you can omit the field `invitation`
- If you want to send plain text messages set the key `email_is_html` to `false`

### Request example

```
curl --location --request POST 'http://127.0.0.1:8080/invite/' \
--header 'Content-Type: application/json' \
--data-raw '{"users":[{"full_name":"Eric Cartman","email":"ihateyouguys@southpark.cc"},{"full_name":"Tina belcher","email":"aaaooo@bobsburger.com"}],"invitation":{"start_at":"2030-10-12T07:20:50.52Z","end_at":"2030-10-12T08:20:50.52Z","organizer_email":"meetingorganizer@meeting.com","organizer_full_name":"Mr. Mojo Rising","summary":"This meeting will be about...","location":"https://zoom.us/332324342","description":"Voluptatum ut quis ut. Voluptas qui pariatur quo. Omnis enim rerum dolorum. Qui aut est sed qui voluptatem harum. Consequuntur et accusantium culpa est fuga molestiae in ut. Numquam harum"},"email_subject":"You'\''ve just been invited!","email_body":"<html><body><h1>email body about the invitation/event</h1></body></html>","email_is_html":true}'
```

### Responses

**Successful**

```
HTTP/1.1 200 
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding

{
    "message": "SENT_OK",
    "status_code": 200,
    "error_fields": null
}
```

**Field missing/invalid**

```
HTTP/1.1 400 BAD REQUEST
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding

{
    "message": "INVALID_PAYLOAD",
    "status_code": 400,
    "error_fields": [
        {
            "field": "email_body",
            "message": "",
            "code": "required"
        }
    ]
}
```

**Error**

```
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding

{
    "message": "ERROR",
    "status_code": 500,
    "error_fields": null
}
```