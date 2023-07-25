# grpc-gateway

## Build the project
`make run`

# REST API
## User Registration
```sh
curl --location 'localhost:8080/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"phongnt",
    "password":"phong",
    "name":"Nguyen Thanh Phong",
    "email":"thanhphong9718@gmail.com",
    "gender":"name"
}'

```
## Login
```sh
curl --location 'localhost:8080/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"phongnt",
    "password":"phong"
}'
```
## Get user profile
```sh
curl --location 'localhost:8080/user/2'
```
