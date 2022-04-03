package models

import "github.com/dgrijalva/jwt-go" 

type Credentials struct {
	Username string `bson:"username" validator:"required, min=3, max=6"`
	Password string `bson:"password" validator:"required, min=6"`
}

type Claims struct {
	Username           string `bson:"username"`
	jwt.StandardClaims        //this will give you some extra info related to the users jwt token (hover over StandardClaims to see those fields)
}
