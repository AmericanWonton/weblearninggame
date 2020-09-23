package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// WARNING! THIS CODE CONTAINS DEROGATORY TERMS, RACIAL/ETHNIC/SEXUAL SLURS,
// AND OTHER OFFENSIVE CONTENT. THE PURPOSE IS TO REMOVE THIS CONTENT OFF OF
// MY PLATFORM. IF ANY OF THIS CONTENT OFFENDS YOU, I APOLOGIZE; PLEASE STAY OFF
// OF THIS PAGE!!!

var allUsernames []string
var usernameMap map[string]bool

/* DEFINED SLURS */
var slurs []string = []string{"penis", "vagina", "dick", "cunt", "asshole", "fag", "faggot",
	"nigglet", "nigger", "beaner", "wetback", "wet back", "chink", "tranny", "bitch", "slut",
	"whore", "fuck", "damn",
	"shit", "piss", "cum", "jizz"}

func containsLanguage(theText string) bool {
	hasLanguage := false
	textLower := strings.ToLower(theText)
	for i := 0; i < len(slurs); i++ {
		if strings.Contains(textLower, slurs[i]) {
			hasLanguage = true
			return hasLanguage
		}
	}
	return hasLanguage
}

//Checks the Usernames after every keystroke
func checkUsername(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 4 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 20 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, usernameMap[sbs])
	}
}

//Loads all our Usernames when the document loads.
func loadUsernames(w http.ResponseWriter, req *http.Request) {
	/* DEBUG NOTE: I SHOULD RE-WRITE THIS TO USE CHANNELS AT SOME POINT */
	//Query the database for all names
	usernameMap, goodCheck := getAllPlayers() //Clear Map for future use on page load

	fmt.Printf("DEBUG: Need to print a value of usernameMap: %v\n", usernameMap["me"])

	if goodCheck == false {
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, "true")
	}
}
