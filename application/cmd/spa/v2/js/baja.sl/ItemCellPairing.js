import { AbstractViewController } from "./AbstractViewController.js";
import { Size, Position } from "./Layout/Coordinates.js";
import { iota, symbolizer } from "./utilities.js";
import { reuseCollector } from "./ReuseCollector.js";
import { AbstractManagedLayoutCellViewController } from "./AbstractManagedLayoutCellViewController.js";
import { resizeObserverWrapper } from "./ResizeObserverWrapper.js";
import { itemMeasurer } from "./ItemMeasurer.js";

/** Bi-directional Map, use it when you need a map you have to search on values. */
class BiMap {
    constructor() {
        this.forwards = new Map();
        this.backwards = new Map();
    }

    set(term1, term2) {
        this.forwards.set(term1, term2);
        this.backwards.set(term2, term1);
    }

    get(key) {
        return this.forwards.get(key) ?? this.backwards.get(key);
    }

    has(key) {
        return this.forwards.has(key) ? true : this.backwards.has(key);
    }
}

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} EnvironmentSymbol
 * @typedef {Symbol} CellSymbol - a unique identifier created on cell creation and assigned its container node as dataset entry
 */

/**
 * @param {HTMLElement} left
 * @param {HTMLElement} right
 */
function findNearestCommonParentNode(left, right) {
    return document.body;
}

/**
 * @param {HTMLElement} parent
 * @param {HTMLElement} node
 * @returns {Position}
 */
function calculateRecursivePositioning(parent, node) {
    if (node === parent) {
        return new Position(0, 0);
    } else {
        const totalPosition = calculateRecursivePositioning(parent, node.parentElement);
        return totalPosition.add(node.clientLeft, node.clientTop);
    }
}

const DEFAULT_ENVIRONMENT = symbolizer.symbolize("DEFAULT_ENVIRONMENT");

class ItemCellPairing {
    constructor() {
        /**
         * @private
         * @type {Map.<ItemSymbol, Map.<EnvironmentSymbol, CellTypeSymbol>>}
         */
        this._cellTypes = new Map();

        /**
         * @private
         * @type {Map.<CellSymbol, AbstractManagedLayoutCellViewController>}
         */
        this._createdCells = new Map();

        /**
         * @private
         * @type {Map.<ItemSymbol, CellSymbol>}
         * Those items which assigned to a cell
         */
        this._assignedItems = new Map();

        this._itemSymbolToCellSymbol = new Map();
        this._cellSymbolToItemSymbol = new Map();

        this._itemSymbolToEnvironmentSymbol = new Map();

        /**
         * @private
         * @type {Map.<ItemSymbol, AbstractManagedLayoutCellViewController>}
         */
        // this._lastAssignedController = new Map();

        /**
         * @private
         * @type {Map.<EnvironmentSymbol, AbstractViewController>}
         */
        // this._registeredViewControllers;
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {CellTypeSymbol} cellTypeSymbol
     */
    setCellKindForItem(itemSymbol, environmentSymbol, cellTypeSymbol) {
        var cellKinds = this._cellTypes.get(itemSymbol);
        if (cellKinds === undefined) {
            cellKinds = new Map();
            this._cellTypes.set(itemSymbol, cellKinds);
        }
        cellKinds.set(DEFAULT_ENVIRONMENT, cellTypeSymbol);
        cellKinds.set(environmentSymbol, cellTypeSymbol);
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @return {CellTypeSymbol}
     */
    getCellTypeForItem(itemSymbol, environmentSymbol) {
        const cellTypes = this._cellTypes.get(itemSymbol);
        if (cellTypes) {
            return cellTypes.get(environmentSymbol) ?? cellTypes.get(DEFAULT_ENVIRONMENT) ?? undefined;
        }
        return undefined;
    }

    registerEnvironment(vcSymbol, exportCallback, resizeNotification) {
        this._registeredViewControllers;
    }

    resizeEventHandler() {}

    /**
     * @private
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _resizeObserverNotificationHandler(entries) {
        // console.log(entries);
        var ignoreChanges = true;
        var changedItems = [];
        entries.forEach((entry) => {
            entry.borderBoxSize;
            const newDimensions = new Size(entry.contentRect.width, entry.contentRect.height);

            const cellMigrationContainer = entry.target;

            const itemId = cellMigrationContainer.dataset["itemId"];
            const itemSymbol = symbolizer.symbolize(itemId);

            const isStillAllocated = this.computedValues.cellMigrationContainers.has(itemSymbol);
            if (!isStillAllocated) return;

            const isSameHeight =
                this.computedValues.lastRecordedCellHeightOfItem.has(itemSymbol) &&
                this.computedValues.lastRecordedCellHeightOfItem.get(itemSymbol) === height;

            if (isSameHeight) return;

            this.computedValues.lastRecordedCellHeightOfItem.set(itemSymbol, height);
            // changedItems.push(itemId);
            ignoreChanges = false;

            // const isInSight = this.computedValues.current.zoneCollusions.inViewport.has(itemSymbol);
            // if (isInSight) {
            //     ignoreChanges = false;
            // }
        });
        // console.log("resizeObserver detected height change for:", changedItems);
        if (ignoreChanges) return;
        // to avoid infinite resize loops
        if (this.computedValues.timeoutToResizeCallback === undefined) {
            this.computedValues.timeoutToResizeCallback = setTimeout(() => {
                requestAnimationFrame(() => {
                    this.computedValues.timeoutToResizeCallback = undefined;
                    this.updateView(TRIGGER_RESIZE_OBSERVER);
                });
            }, 500);
        }
    }

    // /**
    //  * @param {CellTypeSymbol} cellTypeSymbol
    //  * @param {function():AbstractManagedLayoutCellViewController} cellReturningFunction
    //  */
    // registerViewControllerConstructor(cellTypeSymbol, cellReturningFunction) {
    //     reuseCollector.registerViewControllerConstructor(cellTypeSymbol, cellReturningFunction);
    // }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {AbstractViewController} toController
     * @returns {AbstractViewController}
     */
    // createCell(itemSymbol, cellTypeSymbol, toController) {
    //     this._cellTypes.set(itemSymbol, cellTypeSymbol);

    //     const cellMigrationContainer = reuseCollector.get(cellTypeSymbol);

    //     const toContainer = toController.dom.container;
    //     this._currentController.set(itemSymbol, toController);

    //     cellMigrationContainer;

    //     return cellMigrationContainer;
    // }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @returns {AbstractManagedLayoutCellViewController}
     */
    assign(itemSymbol, environmentSymbol) {
        const cellTypeSymbol = this.getCellTypeForItem(itemSymbol, environmentSymbol);
        const managedLayoutCellViewController = reuseCollector.get(cellTypeSymbol, environmentSymbol);

        // save item symbol to cell
        managedLayoutCellViewController.config.itemSymbol = itemSymbol;
        managedLayoutCellViewController.dom.managedLayoutPositioner.dataset.assignedItemSymbol =
            symbolizer.desymbolize(itemSymbol);

        // save assignment information
        const cellSymbol = managedLayoutCellViewController.config.cellSymbol;
        this._itemSymbolToCellSymbol.set(itemSymbol, cellSymbol);
        this._cellSymbolToItemSymbol.set(cellSymbol, itemSymbol);
        // this._assignedItems.set(itemSymbol, managedLayoutCellViewController);

        this._itemSymbolToEnvironmentSymbol.set(itemSymbol, environmentSymbol);

        return managedLayoutCellViewController;
    }

    /**
     * @param {ItemSymbol} itemSymbol
     */
    unassign(itemSymbol) {
        const managedLayoutCellViewController = this.getAssignedCellForItem(itemSymbol);
        const cellSymbol = this._itemSymbolToCellSymbol.get(itemSymbol)
        const cellTypeSymbol = this.getCellTypeForItem(itemSymbol);
        reuseCollector.free(cellTypeSymbol, managedLayoutCellViewController);

        managedLayoutCellViewController.config.itemSymbol = undefined;
        managedLayoutCellViewController.dom.managedLayoutPositioner.dataset.assignedItemSymbol = "";

        this._itemSymbolToCellSymbol.delete(itemSymbol);
        this._cellSymbolToItemSymbol.delete(cellSymbol);

        this._itemSymbolToEnvironmentSymbol.delete(itemSymbol);
        // this._assignedItems.delete(itemSymbol);
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {AbstractViewController} nextContainer
     * @param {Position} positionInNextContainer
     * @param {bool} preserveSpaceOnExporter
     * @returns {AbstractViewController}
     */
    // transfer(itemSymbol, nextContainer, positionInNextContainer, preserveSpaceOnExporter) {
    //     // TODO: call create() if the requested item is not currently assigned to a cell
    //     const currentContainer = this._currentController.get(itemSymbol).dom.container;
    //     const cellMigrationContainerForRequestedItem = this._assignedItems.get(itemSymbol);

    //     const commonParent = findNearestCommonParentNode(currentContainer, nextContainer);
    //     const currentCellTranslationFromCommonParent = calculateRecursivePositioning(
    //         commonParent,
    //         cellMigrationContainerForRequestedItem
    //     );
    //     // prettier-ignore
    //     const nextCellTranslationFromCommonParent = calculateRecursivePositioning(
    //         commonParent,
    //         nextContainer
    //     ).addFrom(positionInNextContainer)

    //     // TODO: compare with nextCellTranslation.deltaCompFrom(currentCellTranslationFromCommonParent)
    //     const neededTranslation = currentCellTranslationFromCommonParent.deltaCompFrom(nextCellTranslation);
    //     adoption(nextContainer.dom.container, [cellMigrationContainerForRequestedItem]);
    //     cellMigrationContainerForRequestedItem.translateFromTo();
    // }

    /**
     * @param {ItemSymbol} itemSymbol
     * @returns {AbstractManagedLayoutCellViewController}
     */
    getAssignedCellForItem(itemSymbol) {
        const cellSymbol = this._itemSymbolToCellSymbol.get(itemSymbol);
        return this._createdCells.get(cellSymbol);
    }

    isItemAssignedToACell(itemSymbol) {
        return this._itemSymbolToCellSymbol.has(itemSymbol);
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {function():AbstractManagedLayoutCellViewController} constructorFunction
     */
    registerCellViewControllerConstructor(cellTypeSymbol, environmentSymbol, constructorFunction) {
        reuseCollector.registerCellViewControllerConstructor(cellTypeSymbol, environmentSymbol, () => {
            const cell = constructorFunction();

            const cellSymbolString = `cellSymbol-${iota()}`;
            const cellSymbol = symbolizer.symbolize(cellSymbolString);
            cell.config.cellSymbol = cellSymbol;
            cell.dom.container.dataset.cellSymbolStringified = cellSymbolString;
            this._createdCells.set(cellSymbol, cell);

            resizeObserverWrapper.subscribe(
                cell.dom.container,
                this._resizeObserverEventHandler.bind(this, cellSymbol)
            );
            return cell;
        });
    }

    /**
     * @private
     * @param {CellSymbol} cellSymbol
     */
    _resizeObserverEventHandler(cellSymbol) {
        // const cellSymbol = symbolizer.symbolize(cellSymbolStringified);
        const cell = this._createdCells.get(cellSymbol);
        const itemSymbol = this._cellSymbolToItemSymbol.get(cellSymbol);
        const environmentSymbol = this._itemSymbolToEnvironmentSymbol.get(itemSymbol);

        if (itemSymbol) {
            itemMeasurer.setSize(itemSymbol, environmentSymbol, cell.measureSize());
        }
    }
}

export const itemCellPairing = new ItemCellPairing();
