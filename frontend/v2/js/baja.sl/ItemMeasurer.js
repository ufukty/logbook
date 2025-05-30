import { symbolizer } from "./utilities.js";
import { Size } from "./Layout/Coordinates.js";
import { DelegateRegistry } from "./DelegateRegistry.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} EnvironmentSymbol
 */

const LAST_ENVIRONMENT = symbolizer.symbolize("LAST_ENVIRONMENT");

class ItemMeasurer {
    constructor() {
        /**
         * @private
         * @type {Map.<EnvironmentSymbol, Size>}
         * Average size of any kind of element in a specific environment
         */
        this._averageSizes = new Map();
        /**
         * @private
         * @type {Map.<ItemSymbol, Map.<EnvironmentSymbol, Size>>}
         * Default size of specifc item in a specific environment before item assigned to a cell
         */
        this._defaultSizes = new Map();
        /**
         * @private
         * @type {Map.<ItemSymbol, Map.<EnvironmentSymbol, Size>>}
         * Default size of specifc item in a specific environment after item assigned to a cell
         */
        this._measuredSizes = new Map();

        /** @private */
        this._subscribers = new DelegateRegistry();
    }

    /**
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {Size} size
     */
    setAverageSize(environmentSymbol, size) {
        this._averageSizes.set(environmentSymbol, size);
    }

    /**
     * @param {EnvironmentSymbol} environmentSymbol
     * @returns {Size}
     */
    getAverageSize(environmentSymbol) {
        return this._averageSizes.get(environmentSymbol);
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {Size} size
     * Default size is the size which assumed to fit a specific unassigned-item.
     */
    setDefaultSize(itemSymbol, environmentSymbol, size) {
        var sizes = this._defaultSizes.get(itemSymbol);
        if (sizes === undefined) {
            sizes = new Map();
            this._defaultSizes.set(itemSymbol, sizes);
        }
        sizes.set(environmentSymbol, size);
        this._subscribers.notify(environmentSymbol);
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {Size} size
     */
    setSize(itemSymbol, environmentSymbol, size) {
        var sizes = this._defaultSizes.get(itemSymbol);
        if (sizes === undefined) {
            sizes = new Map();
            this._defaultSizes.set(itemSymbol, sizes);
        }
        sizes.set(LAST_ENVIRONMENT, size);
        sizes.set(environmentSymbol, size);
        this._subscribers.notify(environmentSymbol);
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @returns {Size}
     */
    getSize(itemSymbol, environmentSymbol) {
        const measuredSizes = this._measuredSizes.get(itemSymbol);
        const defaultSizes = this._defaultSizes.get(itemSymbol);
        if (measuredSizes && measuredSizes.has(environmentSymbol)) return measuredSizes.get(environmentSymbol);
        if (defaultSizes && defaultSizes.has(environmentSymbol)) return defaultSizes.get(environmentSymbol);
        if (measuredSizes && measuredSizes.has(LAST_ENVIRONMENT)) return measuredSizes.get(LAST_ENVIRONMENT);
        if (defaultSizes && defaultSizes.has(LAST_ENVIRONMENT)) return defaultSizes.get(LAST_ENVIRONMENT);
        console.error("asked for an item size unavailable currently");
        return new Size(0, 0);
    }

    /**
     * @param {Symbol} environmentSymbol
     * @param {function} callback - This function will be called after an item size updated while it is in specified environment
     */
    subscribe(environmentSymbol, callback) {
        this._subscribers.add(environmentSymbol, callback);
    }
}

export const itemMeasurer = new ItemMeasurer();
