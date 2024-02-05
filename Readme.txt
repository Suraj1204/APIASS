
For running the code:

**** go run main.go

Then open new terminal in same path, then run below curl command:

Sign Up:
**** curl -X POST -d '{"email": "suraj@gmail.com", "password": "Suraj123"}' http://localhost:8080/signup


Sign In:
**** curl -X POST -d '{"email": "suraj@gmail.com", "password": "Suraj123"}' http://localhost:8080/signin

Authorization of token:
**** curl -X GET -H "Authorization: <Token Which is found from Sign In>" http://localhost:8080/protected
Replace <TOKEN> with the JWT token obtained from the sign-in response. This request should succeed if the token is valid.

Ex:
curl -X POST -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MDcxOTk0ODh9.6dLKYZ8tmV1EDss0TsOVimumc7vHD363vnSiAgpeUos" http://localhost:8080/refresh


Revoke:
**** curl -X POST -H "Content-Type: application/json" -d '{"token": "Token"}' http://localhost:8080/revoke

Ex:
curl -X POST -H "Content-Type: application/json" -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MDcxMTkzNjl9.73Y8TnnUvAH8xExyFA8M2XfTyh8oIDSdB9CyRxulbGg"}' http://localhost:8080/revoke


Referesh
**** curl -X POST -H "Authorization: <Existing token>" http://localhost:8080/refresh

Ex
curl -X POST -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE3MDcxOTk0ODh9.6dLKYZ8tmV1EDss0TsOVimumc7vHD363vnSiAgpeUos" http://localhost:8080/refresh



