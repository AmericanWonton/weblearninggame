package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

/* INFORMATION FOR OUR EMAIL VARIABLES */
var senderAddress string
var senderPWord string

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

//Parse our templates
func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
	getCreds() //Get creds for various variables
}

//Initial Index page for warning new Users before entering our site
func indexPage(w http.ResponseWriter, r *http.Request) {

	err1 := template1.ExecuteTemplate(w, "index.gohtml", "")
	HandleError(w, err1)
}

//Login Page to handle credential checks to enter our site.
func loginPage(w http.ResponseWriter, r *http.Request) {
	loggedIn := false
	if loggedIn == true {

	} else {
		if r.Method == http.MethodPost {
			/* DETERMINE IF THIS IS CREATING AN ACCOUNT OR LOGGING IN,
			THEN DIRECT USERS TO THE APPROPRIATE SPOT IF SUCCESSFUL */
			//Collect JSON
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}
			//Declare DataType from Ajax
			type UserData struct {
				ThePlayer Player `json:"ThePlayer"`
				Action    string `json:"Action"`
				PassConf  string `json:"PassConf"`
			}
			//Marshal all neccessary data
			var dataPosted UserData
			json.Unmarshal(bs, &dataPosted)
			logWriter("User given through Ajax: " + string(bs))
			//Assign values from data
			username := dataPosted.ThePlayer.Username
			password := dataPosted.ThePlayer.Password
			email := dataPosted.ThePlayer.Email
			action := dataPosted.Action
			//Declare a message to send back to Ajax if this fails and why
			type successMSG struct {
				Message   string `json:"Message"`
				ResultNum int    `json:"ResultNum"`
			}
			/* DETERMINE IF THIS IS SIGNING IN OR CREATING ACCOUNT */
			if strings.Contains(action, "createuser") {
				/* ATTEMPT TO SEND USER AN EMAIL; IF IT FAILS, GIVE ERROR MESSGAE TO AJAX */
				emailSend, theErr := sendEmailToPlayer(email, dataPosted.ThePlayer)
				if emailSend == false {
					//Inform User of failure
					//Log error and return Ajax with failure
					message := "Error with sending Player email: " + theErr.Error()
					fmt.Printf("DEBUG: %v\n", message)
					logWriter(message)
					successSend := successMSG{
						Message:   message,
						ResultNum: -1,
					}
					theJSONMessage, err := json.Marshal(successSend)
					if err != nil {
						fmt.Println(err)
						logWriter(err.Error())
					}
					fmt.Fprint(w, string(theJSONMessage))
				} else {
					//Begin Creating Player User in database
					theTimeNow := time.Now() //Time Definition
					playerSend := Player{
						UserID:      idCreation("Player"),
						Email:       email,
						Username:    username,
						Password:    passwordEncrypt(password),
						DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
						DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
					}
					createSuccessfully := simplePlayerCreate(playerSend)
					//Relay news if created successfully or not
					if createSuccessfully == true {
						//Send success message
						message := "Player created successfully"
						fmt.Printf("DEBUG: %v\n", message)
						logWriter(message)
						successSend := successMSG{
							Message:   message,
							ResultNum: 1,
						}
						theJSONMessage, err := json.Marshal(successSend)
						if err != nil {
							fmt.Println(err)
							logWriter(err.Error())
						}
						fmt.Fprint(w, string(theJSONMessage))
					} else {
						//Log error and return Ajax with failure
						message := "Error with Player insretion"
						fmt.Printf("DEBUG: %v\n", message)
						successSend := successMSG{
							Message:   message,
							ResultNum: -2,
						}
						theJSONMessage, err := json.Marshal(successSend)
						if err != nil {
							fmt.Println(err)
							logWriter(err.Error())
						}
						fmt.Fprint(w, string(theJSONMessage))
					}
				}
			} else if strings.Contains(action, "signin") {
				/* QUERY MONGO TO RETURN A USER WITH THIS USERNAME,(also checks if password matches) */
				returnedPlayer, anErr := getAPlayer(username, password)
				switch anErr {
				case "all good":
					//Player found, load mainpage with this information
					//User found for Username/Password; update session with information
					createNewSession(w, returnedPlayer)
					http.Redirect(w, r, "/mainPage", http.StatusSeeOther)
					break
				case "all bad":
					//There was an error somewhere...go back to login page and inform Player/User
					message := "Unique error for Player while logging in: " + anErr
					fmt.Printf("DEBUG: %v\n", message)
					logWriter(message)
					successSend := successMSG{
						Message:   message,
						ResultNum: 0,
					}
					theJSONMessage, err := json.Marshal(successSend)
					if err != nil {
						fmt.Println(err)
						logWriter(err.Error())
					}
					fmt.Fprint(w, string(theJSONMessage))
					break
				case "password bad":
					//Inform User that they entered the wrong password in Ajax
					//Log error and return Ajax with failure
					message := "Incorrect Password entered: " + password
					fmt.Printf("DEBUG: %v\n", message)
					logWriter(message)
					successSend := successMSG{
						Message:   message,
						ResultNum: -3,
					}
					theJSONMessage, err := json.Marshal(successSend)
					if err != nil {
						fmt.Println(err)
						logWriter(err.Error())
					}
					fmt.Fprint(w, string(theJSONMessage))
					break
				default:
					//Unexpected switch output; return an error to Ajax
					message := "Error with User Post; wrong action for Player check: " + anErr
					fmt.Printf("DEBUG: %v\n", message)
					logWriter(message)
					successSend := successMSG{
						Message:   message,
						ResultNum: 0,
					}
					theJSONMessage, err := json.Marshal(successSend)
					if err != nil {
						fmt.Println(err)
						logWriter(err.Error())
					}
					fmt.Fprint(w, string(theJSONMessage))
					break
				}
			} else {
				//Log error and return Ajax with failure
				message := "Error with User Post; wrong action: " + action
				fmt.Printf("DEBUG: %v\n", message)
				logWriter(message)
				successSend := successMSG{
					Message:   message,
					ResultNum: 0,
				}
				theJSONMessage, err := json.Marshal(successSend)
				if err != nil {
					fmt.Println(err)
					logWriter(err.Error())
				}
				fmt.Fprint(w, string(theJSONMessage))
			}

		} else {
			//Just serve this page normally and allow User to login or create account
			err1 := template1.ExecuteTemplate(w, "loginPage.gohtml", nil)
			HandleError(w, err1)
		}
	}
}

//Main Page
func mainPage(w http.ResponseWriter, r *http.Request) {
	/* CHECK TO SEE IF USER HAS THE CORRECT COOKIE CREDENTIALS TO STAY LOGGED IN */
	if !alreadyLoggedIn(w, r) {
		//No login information for given session username. Returning to login screen
		errMsg := "No given information for username session. Going back to login screen."
		fmt.Println(errMsg)
		logWriter(errMsg)
		http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
	} else {
		//Player session found, get player information to pass onto the mainpage
		//Get Player for information gathering
		thePlayer, theErr := getUser(w, r)
		//Put User back to login screen if there's an error
		switch theErr {
		case "all good":
			//All good, display mainpage with user stats
			mpViewData := MPViewData{
				thePlayer, thePlayer.Username, thePlayer.UserID,
			}
			err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", mpViewData)
			HandleError(w, err1)
			break
		case "missing username":
			//No Username for session/Player...return to login screen
			errMsg := "Missing Username for session."
			fmt.Println(errMsg)
			logWriter(errMsg)
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			break
		case "missing password":
			//No Password for session/Player...return to login screen
			errMsg := "Missing Password for session."
			fmt.Println(errMsg)
			logWriter(errMsg)
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			break
		default:
			//Wrong error message given, go back to login page.
			errMsg := "Wrong switch statement given for session check in mainPage: " + theErr
			fmt.Println(errMsg)
			logWriter(errMsg)
			http.Redirect(w, r, "/loginPage", http.StatusSeeOther)
			break
		}

	}

}

//Handles all incoming requests
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	//Serving Webpages
	myRouter.HandleFunc("/", indexPage)
	myRouter.HandleFunc("/loginPage", loginPage)
	myRouter.HandleFunc("/mainPage", mainPage)
	//Validation checking
	myRouter.HandleFunc("/loadUsernames", loadUsernames) //Loads Usernames
	//Middleware logging
	myRouter.Handle("/", loggingMiddleware(http.HandlerFunc(logHandler)))
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Connect to MongoDB
	mongoClient = connectDB()
	defer mongoClient.Disconnect(theContext) //Disconnect in 10 seconds if you can't connect

	handleRequests() //Handle incoming webrequests
}

//Hanling Web Errors
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		logWriter("Error handling WebPage: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Main Log Writer
func logWriter(logMessage string) {
	//Logging info
	theTimeNow := time.Now() //Set time for logging
	messageWrite := theTimeNow.Format("2006-01-02 15:04:05") + ":  "
	fmt.Println("Writing log files.")
	logFile, err := os.OpenFile("logging/superDBAppLog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	defer logFile.Close()

	if err != nil {
		//log.Fatalln("Failed opening file")
		fmt.Println("Failed opening file")
	}

	log.SetOutput(logFile)

	messageWrite = messageWrite + logMessage
	log.Println(messageWrite)
}

//Some stuff for logging
func logHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Package main, son")
	fmt.Fprint(w, "package main, son.")
}

//Some other Stuff for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infof("uri: %v\n", req.RequestURI)
		next.ServeHTTP(w, req)
	})
}
