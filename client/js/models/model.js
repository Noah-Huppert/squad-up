/**
 * A base class for data in the application. Provides base methods for saving data and fetching
 * updates from the server.
 */
class Model {
    class() {
        return "Model";
    }

    constructor(name) {
        // Model name, must be the same as the IDB store name
        this.name = name;

        // Lazy load model primary key
        this.id = -1;

        // Lazy loaded model data
        this.data = undefined;
    }

    /**
     * Load model from cache or API with id
     * @param state State to access
     * @throws {CheckFailedError} On argument checks fail
     */
    load(state) {
        if (this.id === -1) {
            throw new Error("Can not load model with no id");
        }

        checkAndThrow(arguments, {typ: "State"});


    }
}
