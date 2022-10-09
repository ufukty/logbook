import { Size } from "./Coordinates.js";

class ItemAccountant {
    constructor() {
        /** @type {Map.<Symbol, Size>} */
        this._defaultSizes = new Map();
        /** @type {Map.<Symbol, {lastEnvironment: Size}>} */
        this._measuredSizes = new Map();
    }

    setAverageSize(environmentSymbol, size) {
        this._averageSizes;
    }

    /**
     * @param {Symbol} itemSymbol
     * @param {Symbol} environmentSymbol
     * @param {Size} size
     */
    setDefaultSize(itemSymbol, environmentSymbol, size) {
        const current = this._defaultSizes.get(itemSymbol) ?? {};
        this._defaultSizes.set(itemSymbol, {
            ...current,
            lastEnvironment: size,
            [environmentSymbol]: size,
        });
    }

    /**
     * @param {Symbol} itemSymbol
     * @param {Symbol} environmentSymbol
     * @param {Size} size
     */
    setSize(itemSymbol, environmentSymbol, size) {
        const current = this._measuredSizes.get(itemSymbol) ?? {};
        this._measuredSizes.set(itemSymbol, {
            ...current,
            lastEnvironment: size,
            [environmentSymbol]: size,
        });
    }

    /**
     * @param {Symbol} itemSymbol
     * @param {Symbol} environmentSymbol
     * @returns {Size}
     */
    getSize(itemSymbol, environmentSymbol) {
        const measuredSizes = this._measuredSizes.get(itemSymbol);
        const defaultSizes = this._defaultSizes.get(itemSymbol);
        if (measuredSizes && measuredSizes.hasOwn(environmentSymbol)) return measuredSizes[environmentSymbol];
        if (defaultSizes && defaultSizes.hasOwn(environmentSymbol)) return defaultSizes[environmentSymbol];
        if (measuredSizes && measuredSizes.hasOwn("lastEnvironment")) return measuredSizes["lastEnvironment"];
        if (defaultSizes && defaultSizes.hasOwn("lastEnvironment")) return defaultSizes["lastEnvironment"];
        console.error("asked for an item size unavailable currently");
        return Size(0, 0);
    }

    /**
     * @param {Symbol} environmentSymbol
     */
    subscribeForSizeChanges(environmentSymbol) {
        //
    }
}

export const itemAccountant = new ItemAccountant();
