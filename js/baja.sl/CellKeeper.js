import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { Size, Position } from "./Coordinates.js";
import { domCollector, adoption, createElement, symbolizer } from "./utilities.js";
/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} ViewControllerSymbol
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

class CellMigrationContainer {
    constructor() {
        /** @type {AbstractViewController} */
        this.cell = undefined;
        /** @type {CellTypeSymbol} */
        this.cellTypeSymbol = undefined;
        /** @type {{container: HTMLElement}} */
        this.dom = {
            container: createElement("div", ["baja-sl-cell-migration-container"]),
        };
    }

    setPosition(x, y) {}

    /**
     * @param {Position} from
     * @param {Position} to
     * @param {boolean} optimizeTransitionForEndPosition When the value false,
     *   element will start transition from starting position. When it is true,
     *   element will be instantly moved to end position and start to transition
     *   from old position to end position. This is required because mobile
     *   browsers optimize page performance by reducing framerate of animations
     *   of out-of-sight elements.
     */
    translateFromTo(from, to, optimizeTransitionForEndPosition = false) {}
}

class CellDispatcher {
    constructor() {
        /**
         * @type {Map.<ItemSymbol, CellMigrationContainer>}
         * Those items which assigned to a cell
         */
        this._assignedItems = new Map();

        /** @type {Map.<ItemSymbol, CellTypeSymbol>} */
        this._cellTypes = new Map();

        /** @type {ResizeObserver} */
        this._resizeObserver = new ResizeObserver(this._resizeObserverNotificationHandler.bind(this));

        /** @type {Map.<ItemSymbol, Size>} */
        this._observedSizes = new Map();

        /** @type {Map.<ItemSymbol, AbstractViewController>} */
        this._currentController = new Map();

        /** @type {Map.<ViewControllerSymbol, AbstractViewController>} */
        this._registeredViewControllers;
    }

    registerViewController(vcSymbol, exportCallback, resizeNotification) {
        this._registeredViewControllers;
    }

    /**
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _resizeObserverNotificationHandler(entries) {
        // console.log(entries);
        var ignoreChanges = true;
        var changedItems = [];
        entries.forEach((entry) => {
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
    register(cellTypeSymbol, cellReturningFunction) {
        domCollector.registerItemIdentifier(cellTypeSymbol, () => {
            const userProvidedCell = cellReturningFunction();
            const cellContainer = new CellMigrationContainer();
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
    create(itemSymbol, cellTypeSymbol, toController) {
        this._cellTypes.set(itemSymbol, cellTypeSymbol);

        const cellMigrationContainer = domCollector.get(cellTypeSymbol);

        const toContainer = toController.dom.container;
        this._currentController.set(itemSymbol, toController);

        cellMigrationContainer;

        return cellMigrationContainer;
    }

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

export const cellRegistry = new CellDispatcher();

// //

// const taskType = symbolizer.symbolize("task");
// cellRegistry.register(symbolizer.symbolize(taskType), () => {
//     return new InfiniteSheetTask();
// });
// const container1 = createElement("div", ["container-1"]);
// const itemSymbol = symbolizer.symbolize("task-1");

// //

// const migrationContainer = cellRegistry.getCell(itemSymbol, taskType, container1);

// cellRegistry.transfer(itemSymbol, container2, true);
