package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Mongo DB Declarations
var mongoClient *mongo.Client
var theContext context.Context //The context for logging out of the connections
var mongoURI string            //Connection string loaded

func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = theClient.Connect(theContext)
	if err != nil {
		fmt.Printf("Errored getting mongo client context: %v\n", err)
		log.Fatal(err)
	}
	//Double check to see if we've connected to the database
	err = theClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}

	return theClient
}

func getCreds() {
	file, err := os.Open("security/mongoConnections.txt")

	if err != nil {
		fmt.Printf("Trouble opening file for Amazon Credentials: %v\n", err.Error())
		logMessage := "Trouble opening file for Amazon Credentials: " + err.Error()
		logWriter(logMessage)
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	/*
		for _, each_ln := range text {
			fmt.Println(each_ln)
		}
	*/

	mongoURI = text[0]
}
