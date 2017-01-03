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
