# Healtcheck Endpoint

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