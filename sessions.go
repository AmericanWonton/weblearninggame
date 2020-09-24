package main

import (
	"net/http"
	"time"

	"github.com/akyoto/uuid"
)

const sessionLength int = 180 //Length of sessions
//Session Database info
var dbPlayers = map[string]Player{}      // user ID, user
var dbSessions = map[string]theSession{} // session ID, session
var dbSessionsCleaned time.Time

//Here's our session struct
type theSession struct {
	username     string
	gamearea     string
	lastActivity time.Time
}

//Creates a new session after player logs in from login screen
func createNewSession(w http.ResponseWriter, thePlayer Player) {
	uuidWithHyphen := uuid.New().String()

	dbPlayers[thePlayer.Username] = thePlayer
	cookie := &http.Cookie{
		Name:  "playersession",
		Value: uuidWithHyphen,
	}
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
	dbSessions[cookie.Value] = theSession{thePlayer.Username, "mainPage", time.Now()}
}

//Returns a User after a certain page loads
func getUser(w http.ResponseWriter, r *http.Request) (Player, string) {
	errMsg := "all good"
	// get cookie
	cookie, err := r.Cookie("playersession")
	//If there is no session cookie, create a new session cookie
	if err != nil {
		uuidWithHyphen := uuid.New().String()
		cookie = &http.Cookie{
			Name:  "playersession",
			Value: uuidWithHyphen,
		}
	}
	//Set the cookie age to the max length again.
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie) //Set the cookie to our grabbed cookie,(or new cookie)

	// if the user exists already, get user
	var thePlayer Player
	if session, ok := dbSessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
		thePlayer = dbPlayers[session.username]
	}
	//Check to see if Player is REALLY given for session
	if len(thePlayer.Username) <= 0 {
		errMsg = "missing username"
	} else if len(thePlayer.Password) <= 0 {
		errMsg = "missing password"
	} else {
		//Player is good
	}
	return thePlayer, errMsg
}

//Checks to see if the User had the cookie credentials for logging in
func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("playersession")
	if err != nil {
		return false //If there is an error getting the cookie, return false
	}
	//if session is found, we update the session with the newest time since activity!
	session, ok := dbSessions[cookie.Value]
	if ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
	}
	/* Check to see if the Username exists from this Session Username. If not, we return false. */
	_, ok = dbPlayers[session.username]
	// refresh session
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
	return ok
}
