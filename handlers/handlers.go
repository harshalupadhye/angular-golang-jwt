package handlers

import (
	database "angular-go-jwt/database"
	models "angular-go-jwt/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// type Credentials struct {
// 	Username string `bson:"username" validate:"required, min=3, max=6"`
// 	Password string `bson:"password" validate:"required, min=6"`
// }

var jwtKey = []byte("secret_key") //create a byte array to store key which will have name secrete_key

// type Claims struct {
// 	Username           string `bson:"username"`
// 	jwt.StandardClaims        //this will give you some extra info related to the users jwt token (hover over StandardClaims to see those fields)
// }

var collect *mongo.Collection = database.DBinstance() //we have initialize it in main.go

var credentials models.Credentials
var validate = validator.New() //validator has default method to validate

var Ctxs context.Context //we have initialize it in main.go

func SignupHandler(c *gin.Context) {
	r := c.Request                               // r *http.Request
	w := c.Writer                                // w http.ResponseWriter
	validatorErr := validate.Struct(credentials) //this is to check validation given in the struct in models
	if validatorErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": validatorErr.Error()})
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // thats the way to let user know in network tab
		return
	}
	result, err := collect.InsertOne(Ctxs, credentials)
	if err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(w).Encode(result)
	}

}
func SiginHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	var user bson.M
	r := c.Request
	w := c.Writer
	err := json.NewDecoder(r.Body).Decode(&credentials) //fetching the data from request body and passing it through credential var which is of type Credential (ie schema)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = collect.FindOne(Ctxs, bson.M{"username": credentials.Username}).Decode(&user) //checking the user name comming in the body to see if user exists in database and if found the push the found data from database to user var of type bson.M
	fmt.Println("no user", user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		if user["password"] == credentials.Password { // check user["password"] (comming from database with (credentials.Password comming from body))
			expirationTime := time.Now().Add(time.Minute * 5) //expiration time for which session lasts

			claims := models.Claims{ //create an object of type Claims struct
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
			// http.SetCookie(w, //pass it in response
			// 	&http.Cookie{ //cookie struct hover on Cookie to see
			// 		Name:     "token",
			// 		Value:    tokenString,
			// 		Expires:  expirationTime,
			// 		Path:     "/",
			// 		SameSite: http.SameSiteNoneMode, //imp as new browser would consider that the client and server both are same and it is not cross site
			// 		Secure:   true,                  //if you set above to none that means we have to make this true and do csrf code as well
			// 	},
			// )
			c.SetSameSite(http.SameSiteNoneMode)                                   //tells program that client is on diffrent site and server is on diff
			c.SetCookie("token", tokenString, 5*60, "/", "localhost", true, false) //secure: true for cross site
		} else {
			w.WriteHeader(http.StatusUnauthorized) //if not matched say unauthorized
		}

	}

}
func HomeHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") //to read the cookie from request header otherwise wont be able to
	r := c.Request
	w := c.Writer

	cookie, err := r.Cookie("token") // same name as per line 88 this is to get the cookie from request

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
	claims := &models.Claims{}  // create a refrence to the claim object

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
	json.NewEncoder(w).Encode(credentials.Username)
	// w.Write([]byte(fmt.Sprintf("welcome back %s", claims.Username)))
}
func RefreshHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer
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
	claims := &models.Claims{}  // create a refrence to the claim object

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

func LogoutHandler(c *gin.Context){
	fmt.Println("...loging out")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	r := c.Request
	w := c.Writer

	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie { // if cookie is not present in the http request
			fmt.Println("no cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	fmt.Println("string",tokenString)
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims , func(t *jwt.Token) (interface{}, error) {
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

	expirationTime := time.Now() //expiration time for which session lasts

	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //provide hashing method

	tokenString, err = token.SignedString(jwtKey) //convert hashed token into string
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", tokenString, 2, "/", "localhost", true, false)

}