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

    console.log("Username: " + username);
    console.log("password: " + password);
    console.log("action: " + action);
    console.log("passwordConf: " + passwordConf);
    console.log("email: " + email);
}