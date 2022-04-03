package main

import (
	"context"
	"log"
	"time"

	// "net/http"
	"os"

	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	routes "angular-go-jwt/routes"
    database "angular-go-jwt/database"
	handlers "angular-go-jwt/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
)

func main() {
	//routers instance
	// router := mux.NewRouter()
    gin.SetMode(gin.ReleaseMode)
	router := gin.New() //express.Routes()
	router.Use(gin.Logger()) //server log

	// router.SetTrustedProxies([]string{"localhost"}) //trust proxies

    routes.AuthRoutes(router)

	// origins := handlers.AllowedOrigins([]string{"*"})
	// credentials := handlers.AllowCredentials() // to solve cookie not getting stored at browser issue 
	// header := handlers.AllowedHeaders([]string{"*"}) // allow header so the header from client side request will be accepted

	//same as handlers above
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowHeaders: []string{"http://localhost:4200"},
		AllowMethods: []string{"POST, GET"},
		AllowCredentials: true,
	}))
	//loading env from .env
	err := godotenv.Load() //require('dotenv').config()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT") //const { PORT } = process.env
    ctx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	/*this will allow us to perform an action for certain time period and
	if action is not performed suppose in this case in 15 sec then it will throw an error like connection timeout*/
    database.Ctxs = ctx
	handlers.Ctxs = ctx
	//routes
	// router.HandleFunc("/signup", signupHandler).Methods("POST")
	// router.HandleFunc("/signin", siginHandler).Methods("POST")
	// router.HandleFunc("/home", homeHandler)
	// router.HandleFunc("/refresh", refreshHandler)

	//spinning server
	// log.Fatal(http.ListenAndServe(port, handlers.CORS(origins,credentials, header)(router) ))
	router.Run(":" + port)
}
