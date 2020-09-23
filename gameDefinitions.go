package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Player struct {
	UserID      int    `json:"UserID"`
	Email       string `json:"Email"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//Send email to Player when creating account
func sendEmailToPlayer(theEmail string, thePlayer Player) (bool, error) {
	successSend := true

	theMessage := "Hello " + thePlayer.Username + ", here is your password for login: " + thePlayer.Password +
		"We will be sending you" +
		" another email with your 'login code' to play the first game. Have fun!"
	theSubject := "Welcome to weblearninggame"

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderAddress)
	mailer.SetHeader("To", theEmail)
	mailer.SetAddressHeader("Cc", senderAddress, "Joe")
	mailer.SetHeader("Subject", theSubject)
	mailer.SetBody("text/html", theMessage)
	//m.Attach("furries.jpg")

	c := gomail.NewDialer("smtp.gmail.com", 587, senderAddress, senderPWord)
	// Send to me and User
	err := c.DialAndSend(mailer)
	if err != nil {
		successSend = false
	}

	return successSend, err
}

//Creates Random IDS for every one of our objects,(each with their own unique length)
func idCreation(whichObject string) int {
	fmt.Printf("DEBUG: Creating Random ID for Object: %v\n", whichObject)
	finalID := 0        //The final, unique ID to return to the food/user
	randInt := 0        //The random integer added onto ID
	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt
	foundID := false
	/* DETERMINE THE LENGTH OF ID BASED ON OBJECT */
	idLength := 0
	switch whichObject {
	case "Player":
		fmt.Printf("DEBUG: Creating random ID for player: %v\n", whichObject)
		idLength = 8
		break
	default:
		errMessage := "DEBUG: Error determining idLength, wrong whichObject: " + whichObject
		fmt.Printf(errMessage)
		logWriter(errMessage)
		break
	}
	for foundID == false {
		randInt = 0
		randIntString = ""
		//Create the random number, convert it to string
		for i := 0; i < idLength; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		//Once we have a string of numbers, we can convert it back to an integer
		theID, err := strconv.Atoi(randIntString)
		if err != nil {
			errMessage := "We got an error converting a string back to a number: " + err.Error() +
				"\n" + "randInt: " + randInt + "   " + "randIntString: " + randIntString
			fmt.Printf(errMessage)
			logWriter(errMessage)
		}
		/* RUN MONGO QUERY BASED ON WHICHOBJECT AND IDLENGTH*/
		switch idLength {
		case 8:
			//Run Mongo query for Player IDS
			idReturned, returnedError := checkObjectID(whichObject, theID)
			if idReturned == true {
				//Found ID, start another search
				fmt.Printf("DEBUG: Returned ID Error: %v\n", returnedError)
			} else {
				//ID is unique, good to use
				finalID = theID
				foundID = true
			}
			break
		default:
			errMessage := "DEBUG: Error determining idLength, wrong idLength: " + string(idLength)
			fmt.Printf(errMessage)
			logWriter(errMessage)
			foundID = false
			break
		}
	}

	return finalID
}

//Encrypts Password for User
func passwordEncrypt(thePassword string) string {
	bsString := []byte(thePassword)               //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2

	return encodedString
}
