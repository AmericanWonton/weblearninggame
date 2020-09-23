//Variables for managing the signin/create tabs
var hideSignIn = true;
var hideCreate = true;
//Variables for checking various fields
//Create Account
var formInformerCreate = document.getElementById("formInformerCreate"); //Informs User of mistake
var submitButtonCreate = document.getElementById("submitButtonCreate"); //Button to submit if values are correct
var formInfoCreateP = document.getElementById("formInfoCreateP");
var usernameSINCreate = document.getElementById("usernameSINCreate");
var passwordSINCreate = document.getElementById("passwordSINCreate");
var passwordConfSINCreate = document.getElementById("passwordConfSINCreate");
var emailSINCreate = document.getElementById("emailSINCreate");
var inputSINCreate = document.getElementById("inputSINCreate");

//Load all Username values to check
window.addEventListener('DOMContentLoaded', function(){
    console.log("DEBUG: Loading all Usernames for field validations.");
    //Send Ajax to query our MongoDB with Golang
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/loadUsernames', true);
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState == XMLHttpRequest.DONE && xhr.status === 200) {
            var item = xhr.responseText;
            console.log(item);
            if (item == 'true') {
                console.log("We successfully loaded all the data");
            } else {
                console.log("Data was NOT successfully loaded.");
            }
        }
    });
    xhr.send(userName.value);
})

function checkSignIn(signInForm){
    var username = signInForm.usernameSIN.value;
    var password = signInForm.passwordSIN.value;
    var action = signInForm.inputSIN.value;

    console.log("Username: " + username);
    console.log("password: " + password);
    console.log("action: " + action);
}

function checkCreateAccount(createForm){
    var username = createForm.usernameSIN.value;
    var password = createForm.passwordSIN.value;
    var passwordConf = createForm.passwordConfSIN.value;
    var email = createForm.emailSIN.value;
    var action = createForm.inputSIN.value;
    
    var canSendAjax = true; //A check to determine if our values are okay before sending Ajax
    if (checkUserName(username) === false){
        canSendAjax = false;
    }
    if (checkPassword(password) === false){
        canSendAjax = false;
    }
    if (checkPassConf(password, passwordConf) === false){
        canSendAjax = false;
    }
    if (checkEmail(email) === false){
        canSendAjax = false;
    }
    
    if (canSendAjax === true) {
        /* SEND INFO TO AJAX FOR OUR GOLANG SERVER TO VALIDATE */
        var Player = {
            UserID = 0,
            Email = email,
            Username = username,
            Password = password,
            DateCreated = "",
            DateUpdated = ""
        };
        var sendData = {
            ThePlayer = Player,
            Action = action,
            PassConf = passwordConf
        }
        var jsonString = JSON.stringify(sendData);
        //Send Data with Ajax
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/loginPage', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var successMSG = JSON.parse(item);
                if (successMSG.SuccessNum == 0){
                    console.log("DEBUG: Successful User submission, going to main page.");
                    window.location.replace("http://localhost:80");
                } else {
                    console.log("DEBUG: Unsuccessful User submission: " + successMSG.Message);
                    document.getElementById("informFormP").innerHTML = successMSG.Message;
                    document.getElementById("informFormDiv").style.display = "block";
                }
            }
        });
        xhr.send(jsonString);
    } else {
        console.log("Error given, cannot create User with Ajax.")
    }
}

//Reveals/Hides Sign In information based on 'hideSignIn'
function revealnhideSignIn(){
    if (hideSignIn == true){
        //Reveal form for User
        hideSignIn = false;
        document.getElementById("signInForm").reset(); //Reset form
        document.getElementById("pullOutCredDiv1").style.display = "block";  //Reveal form
    } else {
        //Hide Form for User
        hideSignIn = true;
        document.getElementById("signInForm").reset(); //Reset form
        document.getElementById("pullOutCredDiv1").style.display = "none"; //Hide Form
    }
}
//Reveals/Hides Create Account information based on 'hideCreate'
function revealnhideCreate(){
    if (hideSignIn == true){
        //Reveal form for User
        hideSignIn = false;
        document.getElementById("createAccountForm").reset(); //Reset form
        document.getElementById("pullOutCredDiv2").style.display = "block";  //Reveal form
    } else {
        //Hide Form for User
        hideSignIn = true;
        document.getElementById("createAccountForm").reset(); //Reset form
        document.getElementById("pullOutCredDiv2").style.display = "none"; //Hide Form
    }
}
//Checks to see if username is legit
function checkUserName(givenUsername){
    /* DETERMINE IF THE USERNAME HAS A BAD VALUE OR NOT...IF SO, INFORM USER */
    var goodUsername = true; // Variable used to determine if this Username is good
    //Length Check
    if ((givenUsername.length < 5) || (givenUsername.length > 20)){
        console.log("DEBUG: Username too large: " + givenUsername.length);
        formInformerCreate.innerHTML = "Username must be between 6 and 20 characters.";
        goodUsername = false;
    }
    //Language check

    //Check to see if Username is taken

    return goodUsername;
}
//Checks to see if password is legit
function checkPassword(givenPassword){
    /* DETERMINE IF THE PASSWORD HAS A BAD VALUE OR NOT...IF SO, INFORM USER */
    var goodPassword = true; // Variable used to determine if this Username is good
    //Length Check
    if ((givenPassword.length < 5) || (givenPassword.length > 20)){
        console.log("DEBUG: Username too large: " + givenPassword.length);
        formInformerCreate.innerHTML = "Password must be between 6 and 20 characters.";
        goodPassword = false;
    }

    return goodPassword;
}
//Checks to see if passConf is legit
function checkPassConf(givenPassword, givenPassConf){
    /* DETERMINE IF THE PASSWORD AND PASSCONF HAS A BAD VALUE OR NOT...IF SO, INFORM USER */
    var goodPassword = true; // Variable used to determine if this Username is good
    //Does Password Conf match password?
    if (givenPassConf === givenPassword){
        console.log("DEBUG: Password good to go");
    } else {
        console.log("Passwords do not match: " + givenPassword + " " + givenPassConf);
        formInformerCreate.innerHTML = "Passwords do not match";
        goodPassword = false;
    }

    return goodPassword;
}
//Checks to see if email is legit
function checkEmail(givenEmail){
    /* DETERMINE IF THE EMAIL HAS A BAD VALUE OR NOT...IF SO, INFORM USER */
    var goodEmail = true; // Variable used to determine if this Email is good
    //Is email long/short enough?
    if ((givenEmail.length < 5) || (givenEmail.length > 200)){
        console.log("DEBUG: Email has bad length: " + givenEmail.length);
        formInformerCreate.innerHTML = "Email must be 5-200 characters";
        goodPassword = false;
    }
    //Does email contain the correct characters?
    if (givenEmail.includes("@") == false){
        console.log("DEBUG: Email has bad value: " + givenEmail);
        formInformerCreate.innerHTML = "Not a valid email";
        goodPassword = false;
    }

    return goodEmail;
}