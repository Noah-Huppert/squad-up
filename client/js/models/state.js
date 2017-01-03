import Database from "/js/models/database";

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