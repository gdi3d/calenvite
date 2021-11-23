<meta property="og:title" content="Invite Users and Send Calendar Events" />
# Calenvite

A simple microservice designed in [GO](https://golang.org/) using [Echo Microframework](https://echo.labstack.com/) for sending emails and/or calendar invitations to users.


# Features

- Send emails using your Mailgun API credentials.
- Send using a standar SMTP server.
- Support for HTML and Plain Text emails.
- Calendar invitation with RSVP.
- Docker image is built using multistage and alpine image to keep it as small and secure as possible.

# How to Use

### Build the Docker image

Download the repo and build the docker image:  
  
  ```
  $ git clone https://github.com/gdi3d/calenvite
  $ docker build -t calenvite_svc:latest .
  ```
 and use provided [docker-compose files](#Sample-docker-compose-files-included) for more info.

### Or Build the binary

```
$ git clone https://github.com/gdi3d/calenvite
$ go get -d -v
$ go mod download
$ go mod verify
$ go build -a -o calenvite

# run de service
$ ./calenvite
```

### Set the Env vars

There's a few env vars that you need to set when you launch the container in order to work:

```
# If you want use Mailgun API:
CALENVITE_SVC_MAILGUN_DOMAIN: The domain from which the email are going to be sent
CALENVITE_SVC_MAILGUN_KEY: The Mailgun API secret key

# If you want to use SMTP:
CALENVITE_SVC_SMTP_HOST: The host/ip of the SMTP server
CALENVITE_SVC_SMTP_PORT: The port of the SMTP server
CALENVITE_SVC_SMTP_USER: The username to authenticate to the SMTP server
CALENVITE_SVC_SMTP_PASSWORD: The password to authenticate to the SMTP server

# common to both options
CALENVITE_SVC_EMAIL_SENDER_ADDRESS: The email address that would be used to send the email (this value will be used in the FROM part of the email)
CALENVITE_SVC_SEND_USING: MAILGUN or SMTP
CALENVITE_SVC_PORT: Port to expose (optional, default: 8000)
```

### Sample docker-compose files included

```
# mailgun-docker-compose.yml

version: "3.9"
services:
  app_backend:
    image: calenvite_svc:latest
    ports:
      - "8080:8000"
    environment:
      - CALENVITE_SVC_MAILGUN_DOMAIN=mycooldomain.com
      - CALENVITE_SVC_MAILGUN_KEY=abcd1234
      - CALENVITE_SVC_EMAIL_SENDER_ADDRESS=no-reply@mycooldomain.com
      - CALENVITE_SVC_SEND_USING=MAILGUN
```


```
# smtp-docker-compose.yml

version: "3.9"
services:
  app_backend:
    image: calenvite_svc:latest
    ports:
      - "8080:8000"
    environment:
      - CALENVITE_SVC_SMTP_HOST=smtp.mailprovider.com
      - CALENVITE_SVC_SMTP_PORT=587
      - CALENVITE_SVC_SMTP_USER=mysmtpuser
      - CALENVITE_SVC_SMTP_PASSWORD=shhhh.is.secret
      - CALENVITE_SVC_EMAIL_SENDER_ADDRESS=no-reply@mycooldomain.com
      - CALENVITE_SVC_SEND_USING=SMTP
```

# API Docs

## Healtcheck Endpoint

A healthcheck endpoint to test if the service is up and running and a valid configuration is present.

*Note: This healtcheck is only going to check that the environment vars are defined. It will not check if the credentials are valid or not*

### URL

`/healthcheck/`

### Request Example

`curl --location --request GET 'http://127.0.0.1:8080/healthcheck'`

### Responses

**Service ok**

```
HTTP/1.1 200 
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding

null
```

**Service not working**

```
HTTP/1.1 500 Internal Server Error 
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding

{
    "message": "Internal Server Error"
}
```



## Invite Endpoint

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

# Questions, complains, death threats?

You can [Contact me 🙋🏻‍♂️](https://www.linkedin.com/in/adrianogalello/) on LinkedIn if you have any questions. Otherwise you can open a ticket 😉



