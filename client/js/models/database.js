import Dexie from "/lib/dexie";

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