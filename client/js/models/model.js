/**
 * A base class for data in the application. Provides base methods for saving data and fetching
 * updates from the server.
 */
class Model {
    class() {
        return "Model";
    }

    constructor(dbName) {
        // Model name, must be the same as the IDB store name
        this.dbName = dbName;

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
        if (state.db.get()[this.dbName] === undefined) {
            throw new Error(`No IndexedDB collection with name of "${this.dbName}"`);
        }

        if (this.id === -1) {
            throw new Error("Can not load model with no id");
        }

        checkAndThrow(arguments, {typ: "State"});

        var collection = state.db.get()[this.dbName];

        collection.get(this.id)
            .then((model) => {
                console.log("ok", model);
            })
            .catch((err) => {
                console.log("err", err);
            });
    }
}
