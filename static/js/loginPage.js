var hideSignIn = true;
var hideCreate = true;

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
    /* DEBUG PRINTING */
    console.log("Username: " + username);
    console.log("password: " + password);
    console.log("action: " + action);
    console.log("passwordConf: " + passwordConf);
    console.log("email: " + email);
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