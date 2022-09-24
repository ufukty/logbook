export class DelegateRegistry {
    /** @param {Array.<string>} eventList */
    constructor(eventList) {
        /** @type {Map.<String, Array<function>>} */
        this._delegates = new Map();
        if (eventList && eventList.length > 0) {
            eventList.forEach((event) => {
                this._delegates.set(event, []);
            });
        }
    }

    /**
     * @param {string} event
     * @param {function} callback
     */
    add(event, callback) {
        if (!this._delegates.has(event)) {
            console.error("Invalid event for DelegateRegistry.add()");
            return;
        }
        this._delegates.get(event).push(callback);
    }

    /** @param {string} event */
    nofify(event, ...args) {
        if (!this._delegates.has(event)) {
            console.error("Invalid event for DelegateRegistry.notify()");
            return;
        }
        this._delegates.get(event).forEach((callback) => {
            callback(...args);
        });
    }
}
