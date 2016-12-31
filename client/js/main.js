Polymer({
    // Element namhttp://jinzhu.me/gorm/e
    is: "x-app",

    // Polymer Component Lifecycle
    // Ready - Web components and Polymer are a OK
    ready: function() {
        console.log("x-app is ready");

        this.listen(this.$["sign-in-btn"], "google-signin-success", "onGSignInSuccess");
    },

    // Event handlers
    onGSignInSuccess: function(a, b, c) {
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