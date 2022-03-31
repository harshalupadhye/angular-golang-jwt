package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type Credentials struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

var jwtKey = []byte("secret_key") //create a byte array to store key which will have name secrete_key

type Claims struct {
	Username           string `bson:"username"`
	jwt.StandardClaims        //this will give you some extra info related to the users jwt token (hover over StandardClaims to see those fields)
}

var collect *mongo.Collection //we have initialize it in main.go
var credentials Credentials
var ctxs context.Context //we have initialize it in main.go
func signupHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // thats the way to let user know in network tab
		return
	}
	result, err := collect.InsertOne(ctxs, credentials)
	if err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(w).Encode(result)
	}

}
func siginHandler(w http.ResponseWriter, r *http.Request) {
	var user bson.M
	err := json.NewDecoder(r.Body).Decode(&credentials) //fetching the data from request body and passing it through credential var which is of type Credential (ie schema)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = collect.FindOne(ctxs, bson.M{"username": credentials.Username}).Decode(&user) //checking the user name comming in the body to see if user exists in database and if found the push the found data from database to user var of type bson.M
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		if user["password"] == credentials.Password { // check user["password"] (comming from database with (credentials.Password comming from body))
			expirationTime := time.Now().Add(time.Minute * 5) //expiration time for which session lasts

			claims := Claims{ //create an object of type Claims struct
				Username: credentials.Username,
				StandardClaims: jwt.StandardClaims{ //create an object of type StandardClaims struct line 25
					ExpiresAt: expirationTime.Unix(), //convert the int into string
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //provide hashing method
			tokenString, err := token.SignedString(jwtKey)             //convert hashed token into string
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//store the token in cookie
			http.SetCookie(w, //pass it in response
				&http.Cookie{ //cookie struct hover on Cookie to see
					Name:     "token",
					Value:    tokenString,
					Expires:  expirationTime,
					Path:     "/",
					SameSite: http.SameSiteNoneMode, //imp as new browser would consider that the client and server both are same and it is not cross site
					Secure:   true,                  //if you set above to none that means we have to make this true and do csrf code as well
				},
			)
		} else {
			w.WriteHeader(http.StatusUnauthorized) //if not matched say unauthorized
		}

	}

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token") // same name as per line 76 this is to get the cookie from request

	if err != nil {
		if err == http.ErrNoCookie { // if cookie is not present in the http request
			fmt.Println("no cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value //if cookie is present then get the value out of it
	claims := &Claims{}         // create a refrence to the claim object

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { //parse the token
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Write([]byte(fmt.Sprintf("welcome back %s", claims.Username)))
}
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token") // same name as per line 76 this is to get the cookie from request

	if err != nil {
		if err == http.ErrNoCookie { // if cookie is not present in the http request
			fmt.Println("no cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value //if cookie is present then get the value out of it
	claims := &Claims{}         // create a refrence to the claim object

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { //parse the token
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > time.Second*30 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5) //expiration time for which session lasts

	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //provide hashing method

	tokenString, err = token.SignedString(jwtKey) //convert hashed token into string
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//store the token in cookie
	http.SetCookie(w, //pass it in response
		&http.Cookie{ //cookie struct hover on Cookie to see
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		},
	)

}
