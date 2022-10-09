class ResizeObserverWrapper {
    constructor() {
        /**
         * @private
         */
        this.observer = new ResizeObserver(this._dispatch.bind(this));

        /**
         * @private
         * @type {Map.<HTMLElement,Set.<function>>}
         */
        this._subscribers = new Map();
    }

    /**
     * @param {HTMLElement} element
     * @param {function} callback
     */
    subscribe(element, callback) {
        this.observer.observe(element);
        if (!this._subscribers.get(element)) this._subscribers.set(element, new Set());
        this._subscribers.get(element).add(callback);
    }

    /**
     * @private
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _dispatch(entries) {
        for (const entry of entries) {
            const subscribers = this._subscribers.get(entry.target);
            if (subscribers) for (const subscriber of subscribers) subscriber();
        }
    }
}

export const resizeObserverWrapper = new ResizeObserverWrapper();
