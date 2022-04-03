package database

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctxs context.Context
func DBinstance() *mongo.Collection {
	
	//loading env from .env
	err := godotenv.Load(".env") //require('dotenv').config()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientURL := options.Client().ApplyURI("mongodb+srv://root:root@cluster0.kpteb.mongodb.net/myFirstDatabase?retryWrites=true&w=majority") //this how we establish the url for which later server will store data

    client, err := mongo.Connect(Ctxs, clientURL) //this is actual connection ie mongoose.connect(url)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(Ctxs, nil) //this will check the ping of the server and if server is not connected then it will throw an error
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database("AuthUsers") //create database
	collection := database.Collection("Users")
	return collection

}