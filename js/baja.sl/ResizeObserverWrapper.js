import { iota, symbolizer } from "./utilities.js";

class ResizeObserverWrapper {
    constructor() {
        /**
         * @private
         */
        this.observer = new ResizeObserver(this._dispatch.bind(this));

        /**
         * @private
         * @type {Map.<string, Set.<function>>}
         * Each element can have multiple subscribers
         */
        this._subscribers = new Map();
    }

    /**
     * @param {HTMLElement} element
     * @param {function} callback
     */
    subscribe(element, callback) {
        if ((typeof element !== "object" || !(element instanceof HTMLElement)) && typeof element !== Window) {
            console.error("Wrong type of argument has received.");
        }
        if (typeof callback !== "function") {
            console.error("Wrong type of argument has received.");
        }
        if (element.dataset.resizeObserverWrapperId) {
            const resizeObserverWrapperId = element.dataset.resizeObserverWrapperId;
            this._subscribers.get(resizeObserverWrapperId).add(callback);
        } else {
            const resizeObserverWrapperId = `${iota()}`;
            var subscribersForCell = new Set();
            subscribersForCell.add(callback);
            this._subscribers.set(resizeObserverWrapperId, subscribersForCell);
            element.dataset.resizeObserverWrapperId = resizeObserverWrapperId;
            this.observer.observe(element);
        }
    }

    /**
     * @private
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _dispatch(entries) {
        for (const entry of entries) {
            const resizeObserverWrapperId = entry.target.dataset.resizeObserverWrapperId;
            const subscribers = this._subscribers.get(entry.target.dataset.resizeObserverWrapperId);
            for (const subscriber of subscribers) {
                subscriber(entry.target);
            }
        }
    }
}

export const resizeObserverWrapper = new ResizeObserverWrapper();
