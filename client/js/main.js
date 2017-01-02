"use strict";

class Database {
    constructor() {
        // DB name
        this.name = "SquadUpDatabase";

        // Define schema
        this.db = new Dexie(this.name);
        this.opened = false;

        db.version(1).stores({
            users: "id,firstName,lastName,email,profilePictureUrl,createdAt,updatedAt"
        });

        // Open db
        // This open method can be called in the constructor because Mixie automatically holds all
        // db operations until after open is done (Relevant docs: https://github.com/dfahlander/Dexie.js/wiki/Dexie.open())
        db.open().catch(err => {
            console.error(`Failed to open database: "${this.name}", err: ${err.stack || err}`);
        });
    }

    /**
     * Returns database
     * @returns {Dexie} Database
     */
    get () {
        return this.db;
    }
}

/**
 * A base class for data in the application. Provides base methods for saving data and fetching
 * updates from the server.
 */
class Model {
    constructor(state, name) {
        this.state = state;

        // Model name, must be the same as the IDB store name
        this.name = name;

        // Lazy loaded model data
        this.data = undefined;
    }
}

/**
 * Application state, sectioned off into different areas of responsibility (ex., authentication, database)
 * to maintain organization.
 * Basically a masthead.
 */
class State {
    constructor () {
        this.db = new Database();
    }
}

Polymer({
    // Element namhttp://jinzhu.me/gorm/e
    is: "x-app",

    // Polymer Component Lifecycle
    // Ready - Web components and Polymer are a OK
    ready: function() {
        console.log("x-app is ready");

        this.listen(this.$["sign-in-btn"], "google-signin-success", "onGSignInSuccess");
    },

    // Application state
    // Holds all data and logic of the app
    state: new State(),

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