package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//routers instance
	router := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"*"})
	credentials := handlers.AllowCredentials() // to solve cookie not getting stored at browser issue 
	header := handlers.AllowedHeaders([]string{"*"}) // allow header so the header from client side request will be accepted
	//loading env from .env
	err := godotenv.Load() //require('dotenv').config()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT") //const { PORT } = process.env

	//db connection
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	/*this will allow us to perform an action for certain time period and
	if action is not performed suppose in this case in 15 sec then it will throw an error like connection timeout*/

	clientURL := options.Client().ApplyURI("mongodb+srv://root:root@cluster0.kpteb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority") //this how we establish the url for which later server will store data

	client, err := mongo.Connect(ctx, clientURL) //this is actual connection ie mongoose.connect(url)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil) //this will check the ping of the server and if server is not connected then it will throw an error
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database("AuthUsers") //creating a database
	collection := database.Collection("Users")
	collect = collection
	ctxs = ctx
	//routes
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/signin", siginHandler).Methods("POST")
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/refresh", refreshHandler)

	//spinning server
	log.Fatal(http.ListenAndServe(port, handlers.CORS(origins,credentials, header)(router) ))
}
