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