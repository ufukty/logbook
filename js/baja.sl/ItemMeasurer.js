import { symbolizer } from "./utilities.js";
import { Size } from "./Layout/Coordinates.js";

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
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @returns {Size}
     */
    getSize(itemSymbol, environmentSymbol) {
        const measuredSizes = this._measuredSizes.get(itemSymbol);
        const defaultSizes = this._defaultSizes.get(itemSymbol);
        if (measuredSizes && measuredSizes.has(environmentSymbol)) return measuredSizes.has(environmentSymbol);
        if (defaultSizes && measuredSizes.has(environmentSymbol)) return defaultSizes.has(environmentSymbol);
        if (measuredSizes && measuredSizes.has(LAST_ENVIRONMENT)) return measuredSizes.has(LAST_ENVIRONMENT);
        if (defaultSizes && defaultSizes.has(LAST_ENVIRONMENT)) return defaultSizes.has(LAST_ENVIRONMENT);
        console.error("asked for an item size unavailable currently");
        return new Size(0, 0);
    }

    /**
     * @param {Symbol} environmentSymbol
     */
    subscribeForSizeChanges(environmentSymbol) {
        //
    }
}

export const itemMeasurer = new ItemMeasurer();
