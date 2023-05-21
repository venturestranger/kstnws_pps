## Internal API Documentation

###### REST Application Conventions and TLDR:

```
Abbreviations:
1. PPS - Posts Pull Server
2. API - Application Programming Interface (internal)

HTTP Request Consists Of:
1. Destination (host:port/resource)
2. Method (GET, POST, PUT, DELETE)
3. Headers 
4. Request Body (JSON)

HTTP Response Consists of:
1. Status Code
2. Response Body (JSON)

Methods:
1. GET - read (returns status code and response body)
2. POST - create (requires request body, returns status code)
3. PUT - update | push to API (returns status code)
4. DELETE - delete (returns status code)

Status Codes:
1. 200 OK (successful request)
2. 400 Bad Request (wrong query format or non-existent fields)
3. 401 Unathorized (no access due to invalid token)
5. 405 Method Not Allowed (no method specified for a resource)
```

###### API Routing:

```
- /validate/auth [GET] | returns a token if valid key 
* /validate [GET] | returns posts from PPS 
* /validate [POST] | creates a post on PPS
* /validate [PUT] | pushes a post from PPS to API
* /validate [DELETE] | deletes a post from PPS
```

*- endpoints requiring authorization

###### Getting Access:

```
1. Access "/validate/auth" with "key" for fetching a token
2. Set a header "Authorization" to "Bearer @token"
```

###### Resources:


