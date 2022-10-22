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
         */
        this._subscribers = new Map();
    }

    /**
     * @param {HTMLElement} element
     * @param {string} elementId - That information will be returned as argument to the declared callback function
     * @param {function} callback
     */
    subscribe(element, elementId, callback) {
        if ((typeof element !== "object" || !(element instanceof HTMLElement)) && typeof element !== Window) {
            console.error("Wrong type of argument has received.");
        }
        if (typeof elementId !== "string") {
            console.error(`Wrong type of argument has received, value="${elementId}"`);
        }
        if (typeof callback !== "function") {
            console.error("Wrong type of argument has received.");
        }
        element.dataset.resizeObserverWrapperId = elementId;
        var subscribers = this._subscribers.get(elementId);
        if (!subscribers) {
            subscribers = new Set();
            this._subscribers.set(elementId, subscribers);
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
            const elementId = entry.target.dataset.resizeObserverWrapperId;
            const subscribers = this._subscribers.get(elementId);
            for (const subscriber of subscribers) {
                subscriber(elementId);
            }
        }
    }
}

export const resizeObserverWrapper = new ResizeObserverWrapper();
