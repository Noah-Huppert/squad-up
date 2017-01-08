Polymer({
    // Element namhttp://jinzhu.me/gorm/e
    is: "x-app",

    // Polymer Component Lifecycle
    // Ready - Web components and Polymer are a OK
    ready: function() {
        console.log("x-app is ready");
    },

    // Application state
    // Holds all data and logic of the app
    state: new State(),

    // Event handlers
    onGSignInSuccess: function() {
        var idToken = gapi.auth2.getAuthInstance()["currentUser"].get().getAuthResponse().id_token;

        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/v1/auth/token/google");
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.onload = function() {
            console.log(xhr.responseText);
        };
        xhr.send("id_token=" + idToken);
    }
});