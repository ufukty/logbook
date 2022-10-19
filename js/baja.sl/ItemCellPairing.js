import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { Size, Position } from "./Layout/Coordinates.js";
import { domCollector, adoption, createElement, symbolizer } from "./utilities.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} EnvironmentSymbol
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
         * @type {Map.<ItemSymbol, CellPositioner>}
         * Those items which assigned to a cell
         */
        this._assignedItems = new Map();

        /**
         * @private
         * @type {Map.<ItemSymbol, Map.<EnvironmentSymbol, CellTypeSymbol>>}
         */
        this._cellTypes = new Map();

        /**
         * @private
         * @type {ResizeObserver}
         */
        this._resizeObserver = new ResizeObserver(this._resizeObserverNotificationHandler.bind(this));

        /**
         * @private
         * @type {Map.<ItemSymbol, AbstractViewController>}
         */
        this._currentController = new Map();

        /**
         * @private
         * @type {Map.<EnvironmentSymbol, AbstractViewController>}
         */
        this._registeredViewControllers;
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {CellTypeSymbol} cellTypeSymbol
     */
    setCellKindForItem(itemSymbol, environmentSymbol, cellTypeSymbol) {
        var cellKinds = this._defaultSizes.get(itemSymbol);
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

    registerViewController(vcSymbol, exportCallback, resizeNotification) {
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

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {function():AbstractViewController} cellReturningFunction
     */
    registerConstructor(cellTypeSymbol, cellReturningFunction) {
        domCollector.registerItemIdentifier(cellTypeSymbol, () => {
            const userProvidedCell = cellReturningFunction();
            const cellContainer = new CellPositioner();
            cellContainer.cell = userProvidedCell;
            cellContainer.cellTypeSymbol = cellTypeSymbol;
            adoption(cellContainer.dom.container, [userProvidedCell.dom.container]);
            this._resizeObserver.observe(cellContainer.container);
            return cellContainer;
        });
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {AbstractViewController} toController
     * @returns {AbstractViewController}
     */
    createCell(itemSymbol, cellTypeSymbol, toController) {
        this._cellTypes.set(itemSymbol, cellTypeSymbol);

        const cellMigrationContainer = domCollector.get(cellTypeSymbol);

        const toContainer = toController.dom.container;
        this._currentController.set(itemSymbol, toController);

        cellMigrationContainer;

        return cellMigrationContainer;
    }

    /**
     * @param {ItemSymbol} itemSymbol
     * @returns {HTMLElement}
     */
    assign(itemSymbol) {
        const itemTypeIdentifier = this._cellTypes.get(itemSymbol);
        const element = domCollector.get(itemTypeIdentifier);
        element.data
        this._assignedItems.set(itemSymbol, element);
        return element;
    }

    unassign(itemSymbol) {}

    /**
     * @param {ItemSymbol} itemSymbol
     * @param {AbstractViewController} nextContainer
     * @param {Position} positionInNextContainer
     * @param {bool} preserveSpaceOnExporter
     * @returns {AbstractViewController}
     */
    transfer(itemSymbol, nextContainer, positionInNextContainer, preserveSpaceOnExporter) {
        // TODO: call create() if the requested item is not currently assigned to a cell
        const currentContainer = this._currentController.get(itemSymbol).dom.container;
        const cellMigrationContainerForRequestedItem = this._assignedItems.get(itemSymbol);

        const commonParent = findNearestCommonParentNode(currentContainer, nextContainer);
        const currentCellTranslationFromCommonParent = calculateRecursivePositioning(
            commonParent,
            cellMigrationContainerForRequestedItem
        );
        // prettier-ignore
        const nextCellTranslationFromCommonParent = calculateRecursivePositioning(
            commonParent, 
            nextContainer
        ).addFrom(positionInNextContainer)

        // TODO: compare with nextCellTranslation.deltaCompFrom(currentCellTranslationFromCommonParent)
        const neededTranslation = currentCellTranslationFromCommonParent.deltaCompFrom(nextCellTranslation);
        adoption(nextContainer.dom.container, [cellMigrationContainerForRequestedItem]);
        cellMigrationContainerForRequestedItem.translateFromTo();
    }
}

export const itemCellPairing = new ItemCellPairing();
