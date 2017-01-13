/**
 * Application state, sectioned off into different areas of responsibility (ex., authentication, database)
 * to maintain organization.
 * Basically a masthead.
 */
class State {
    class() {
        return "State";
    }

    constructor() {
        this.db = new Database();

        var model = new Model("users");
        model.id = 1;

        model.load(this);
    }
}