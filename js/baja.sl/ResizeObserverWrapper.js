import { iota } from "./utilities";

class ResizeObserverWrapper {
    constructor() {
        /**
         * @private
         */
        this.observer = new ResizeObserver(this._dispatch.bind(this));

        /**
         * @private
         * @type {Map.<string, Set.<function>>}
         */
        this._subscribers = new Map();
    }

    /**
     * @param {HTMLElement} element
     * @param {function} callback
     */
    subscribe(element, callback) {
        const resizeObserverWrapperId = iota().toString();
        element.dataset.resizeObserverWrapperId = resizeObserverWrapperId;

        var subscribers = this._subscribers.get(resizeObserverWrapperId);
        if (!subscribers) {
            const subscribers = new Set();
            this._subscribers.set(resizeObserverWrapperId, subscribers);
        }
        subscribers.add(callback);

        this.observer.observe(element);
    }

    /**
     * @private
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _dispatch(entries) {
        for (const entry of entries) {
            const resizeObserverWrapperId = entry.target.dataset.resizeObserverWrapperId;
            const subscribers = this._subscribers.get(resizeObserverWrapperId);
            for (const subscriber of subscribers) {
                subscriber();
            }
        }
    }
}

export const resizeObserverWrapper = new ResizeObserverWrapper();
