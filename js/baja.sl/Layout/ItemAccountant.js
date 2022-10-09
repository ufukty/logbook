import { Size } from "./Coordinates.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} ViewControllerSymbol
 */

class ItemAccountant {
    constructor() {
        /** @type {Map.<Symbol, Size>} */
        this._defaultSizes = new Map();
        /** @type {Map.<Symbol, {lastEnvironment: Size}>} */
        this._measuredSizes = new Map();
        /** @type {Map.<ItemSymbol, CellTypeSymbol>} */
        this.cellKindForItem = new Map();
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

    dsd() {
        const cellHeight = this.computedValues.lastRecordedCellHeightOfItem.has(itemSymbol)
            ? this.computedValues.lastRecordedCellHeightOfItem.get(itemSymbol)
            : this.getDefaultHeightOfItem(itemSymbol);
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

    getCellKindForItem(itemSymbol) {
        return this._cellKindForItem.get(itemSymbol);
    }
}

export const itemAccountant = new ItemAccountant();
