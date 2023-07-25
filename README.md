# grpc-gateway

## Build the project
`make run`

## API
```sh
1. User Registration:
Endpoint:
curl --location 'localhost:8080/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"phongnt",
    "password":"phong",
    "name":"Nguyen Thanh Phong",
    "email":"thanhphong9718@gmail.com",
    "gender":"name"
}'

2. Login
Endpoint:
curl --location 'localhost:8080/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"phongnt",
    "password":"phong"
}'

3. Get user profile
Endpoint:
curl --location 'localhost:8080/user/2'
```
