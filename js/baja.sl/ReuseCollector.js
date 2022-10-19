import { AbstractManagedLayoutCellViewController } from "./AbstractManagedLayoutCellViewController.js";
import { AbstractViewController } from "./AbstractViewController.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} EnvironmentSymbol
 */

/**
 * Keeps constructor functions for superclasses of AbstractViewController,
 *   calls prepareForFree methods, waits
 */
class ReuseCollector {
    constructor() {
        /**
         * @private
         * @type {Map.<CellTypeSymbol, Array.<AbstractViewController>>}
         * ViewControllers that are currently unassigned to any item, grouped with CellTypeSymbol
         */
        this._reusableViewControllers = new Map();
        /**
         * @private
         * @type {Map.<CellTypeSymbol, function>}
         */
        this._constructors = new Map();
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {function():AbstractViewController} constructor
     */
    registerConstructor(cellTypeSymbol, constructor) {
        if (this._constructors.has(cellTypeSymbol)) return;
        this._reusableViewControllers.set(cellTypeSymbol, []);
        this._constructors.set(cellTypeSymbol, constructor);
    }

    /**
     * @private
     * @param {CellTypeSymbol} cellTypeSymbol
     */
    _createViewController(cellTypeSymbol) {
        const _constructor = this._constructors[cellTypeSymbol];
        const item = _constructor();
        return item;
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @returns {AbstractViewController}
     */
    get(cellTypeSymbol) {
        let viewController;
        if (this._reusableViewControllers.get(cellTypeSymbol).length === 0) {
            viewController = this._createViewController(cellTypeSymbol);
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
     * @param {AbstractViewController} viewController
     */
    free(cellTypeSymbol, viewController) {
        if (viewController instanceof AbstractManagedLayoutCellViewController) {
            viewController.prepareForFree().then(() => {
                this._reusableViewControllers.get(cellTypeSymbol).push(viewController);
            });
        } else {
            if (
                typeof viewController.prepareForFree !== "undefined" &&
                typeof viewController.prepareForFree === "function"
            )
                viewController.prepareForFree();
            this._reusableViewControllers.get(cellTypeSymbol).push(viewController);
        }
    }
}

export const reuseCollector = new ReuseCollector();
