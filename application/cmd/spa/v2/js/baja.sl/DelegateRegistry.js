export class DelegateRegistry {
    constructor() {
        /**
         * @private
         * @type {Map.<Symbol, Array<function>>}
         */
        this.delegates = new Map();
    }

    /**
     * @param {Symbol} event
     * @param {function} callback
     */
    add(event, callback) {
        if (!this.delegates.has(event)) {
            this.delegates.set(event, []);
        }
        this.delegates.get(event).push(callback);
    }

    /** @param {Symbol} event */
    notify(event, ...args) {
        if (!this.delegates.has(event)) {
            return;
        }
        this.delegates.get(event).forEach((callback) => {
            callback(...args);
        });
    }
}
