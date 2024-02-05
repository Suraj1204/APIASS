
For running the code:

**** go run main.go

Then open new terminal in same path, then run below curl command:

Sign Up:

**** curl -X POST -d '{"email": "suraj@gmail.com", "password": "Suraj123"}' http://localhost:8080/signup


Sign In:

**** curl -X POST -d '{"email": "suraj@gmail.com", "password": "Suraj123"}' http://localhost:8080/signin

Authorization of token:

**** curl -X GET -H "Authorization: Token Which is found from Sign In" http://localhost:8080/verify

Replace <TOKEN> with the JWT token obtained from the sign-in response. This request should succeed if the token is valid.



Revoke:

**** curl -X POST -H "Content-Type: application/json" -d '{"token": "Existing Token"}' http://localhost:8080/revoke


Referesh

**** curl -X POST -H "Authorization: Existing token" http://localhost:8080/refresh


