package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

//Mongo DB Declarations
var mongoClient *mongo.Client
var theContext context.Context //The context for logging out of the connections
var mongoURI string            //Connection string loaded

//Connect to the DB upon program entry
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

//Get our credentials when starting the program
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

	mongoURI = text[0]
	senderAddress = text[1]
	senderPWord = text[2]
}

//Check to see if an ID exists for a certain object
func checkObjectID(whichObject string, theID int) (bool, string) {
	idReturned := false
	returnedError := ""
	switch whichObject {
	case "Player":
		//Search Mongo to see if ID appears
		playerCollection := mongoClient.Database("learningdb").Collection("players") //Here's our collection
		var thePlayer Player
		//Give 0 values to determine if these IDs are found
		theFilter := bson.M{
			"$or": []interface{}{
				bson.M{"UserID": theID},
			},
		}
		theErr := playerCollection.FindOne(theContext, theFilter).Decode(&thePlayer)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				idReturned = false //No ID found, it's good to use
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				errMessage := "Another error for finding unique player ID: " + theErr.Error()
				returnedError = errMessage
				logWriter(errMessage)
				idReturned = true
			}
		}
		break
	default:
		errMessage := "DEBUG: Error determining whichObject in checkObjectID, wrong whichObject: " + whichObject
		returnedError = errMessage
		fmt.Printf(errMessage)
		logWriter(errMessage)
		idReturned = true
	}

	return idReturned, returnedError
}

//simple creation of Player, no API
func simplePlayerCreate() {

}
