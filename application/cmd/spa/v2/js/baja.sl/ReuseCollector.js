import { AbstractManagedLayoutCellViewController } from "./AbstractManagedLayoutCellViewController.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} EnvironmentSymbol
 */

/**
 * Keeps constructor functions for superclasses of AbstractManagedLayoutCellViewController,
 *   calls prepareForFree methods, waits
 */
class ReuseCollector {
    constructor() {
        /**
         * @private
         * @type {Map.<CellTypeSymbol, Array.<AbstractManagedLayoutCellViewController>>}
         * ViewControllers that are currently unassigned to any item, grouped with CellTypeSymbol
         */
        this._reusableViewControllers = new Map();
        /**
         * @private
         * @type {Map.<CellTypeSymbol, Map.<EnvironmentSymbol, function():AbstractManagedLayoutCellViewController>>}
         */
        this._constructors = new Map();
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {function():AbstractManagedLayoutCellViewController} constructorFunction
     */
    registerCellViewControllerConstructor(cellTypeSymbol, environmentSymbol, constructorFunction) {
        if (typeof cellTypeSymbol !== "symbol") {
            console.error("Wrong type of argument has received.");
        }
        if (typeof environmentSymbol !== "symbol") {
            console.error("Wrong type of argument has received.");
        }
        if (typeof constructorFunction !== "function") {
            console.error("Wrong type of argument has received.");
        }

        if (!this._constructors.has(cellTypeSymbol)) {
            this._constructors.set(cellTypeSymbol, new Map());
        }
        const constructorsForCellType = this._constructors.get(cellTypeSymbol);
        if (!constructorsForCellType.has(environmentSymbol)) {
            constructorsForCellType.set(environmentSymbol, constructorFunction);
        }

        if (!this._reusableViewControllers.has(cellTypeSymbol)) {
            this._reusableViewControllers.set(cellTypeSymbol, []);
        }
    }

    /**
     * @private
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     */
    _createViewController(cellTypeSymbol, environmentSymbol) {
        const _constructor = this._constructors.get(cellTypeSymbol).get(environmentSymbol);
        const item = _constructor();
        return item;
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @returns {AbstractManagedLayoutCellViewController}
     */
    get(cellTypeSymbol, environmentSymbol) {
        let viewController;
        if (this._reusableViewControllers.get(cellTypeSymbol).length === 0) {
            viewController = this._createViewController(cellTypeSymbol, environmentSymbol);
        } else {
            viewController = this._reusableViewControllers.get(cellTypeSymbol).pop();
        }
        if (typeof viewController.prepareForUse !== "undefined" && typeof viewController.prepareForUse === "function") {
            viewController.prepareForUse();
        }
        return viewController;
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {AbstractManagedLayoutCellViewController} viewController
     */
    free(cellTypeSymbol, viewController) {
        viewController.prepareForFree();
        (async () => {
            await viewController.prepareForFreeAsync();
            this._reusableViewControllers.get(cellTypeSymbol).push(viewController);
        })();
    }
}

export const reuseCollector = new ReuseCollector();
