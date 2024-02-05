package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User represents a user in the system
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User
var TokenKey = []byte("super-secret")

// JWTClaims represents JWT claims
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var revokedTokens = struct {
	tokens map[string]struct{}
	sync.Mutex
}{tokens: make(map[string]struct{})}

var signedInUsers = make(map[string]bool)

// GenerateToken generates a JWT token
func generateToken(email string) (string, error) {
	fmt.Println("Inside Generate Token function")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token expiration time
		},
	})

	return token.SignedString(TokenKey)
}

// RefreshToken refreshes a JWT token
func refreshToken(tokenString string) (string, error) {
	fmt.Println("Inside Refresh Token function")
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return TokenKey, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Generate a new token with extended expiration time
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Email: claims.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // New token expiration time
		},
	})

	return newToken.SignedString(TokenKey)
}

// signUpFunc handles sign up requests
func signUpFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside SignUp function")

	// Parse JSON request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	for _, u := range users {
		if u.Email == user.Email {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
	}

	// Add the user to the list of users
	users = append(users, user)

	// Send response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully\n")
}

// signInFunc handles sign in requests
func signInFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Signin function")
	// Parse JSON request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is already signed in
	// if _, ok := signedInUsers[user.Email]; ok {
	// 	http.Error(w, "User is already signed in", http.StatusForbidden)
	// 	return
	// }

	// Check if user exists and password matches
	for _, u := range users {
		if u.Email == user.Email && u.Password == user.Password {
			// Generate JWT token
			token, err := generateToken(user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Mark user as signed in
			signedInUsers[user.Email] = true

			// Send response with token
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `"token is": "%s"`, token)
			return
		}
	}

	// Send error response if user not found or password incorrect
	http.Error(w, "Invalid Id and Pass", http.StatusUnauthorized)
}

// refreshHandler handles token refresh requests
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Refresh the token
	newToken, err := refreshToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Send response with new token
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `"token": "%s"`, newToken)
}

// authMiddleware is middleware to authenticate requests using JWT token
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if token is revoked
		revokedTokens.Lock()
		defer revokedTokens.Unlock()
		if _, ok := revokedTokens.tokens[tokenString]; ok {
			http.Error(w, "Token is revoked", http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return TokenKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, call the next handler
		next.ServeHTTP(w, r)
	}
}

// VerifyFunc is a protected route
func VerifyFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token Verified")
}

// revoketokenFunc handles token revocation requests
func revoketokenFunc(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add token to the list of revoked tokens
	revokedTokens.Lock()
	revokedTokens.tokens[body.Token] = struct{}{}
	revokedTokens.Unlock()

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Revocation of token completed\n")
}

func main() {
	// Define API endpoints
	http.HandleFunc("/signup", signUpFunc)
	http.HandleFunc("/signin", signInFunc)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/revoke", revoketokenFunc)
	http.HandleFunc("/verify", authMiddleware(VerifyFunc))

	// Starting the code
	fmt.Println("-------------------Code Started--------------------------------")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
