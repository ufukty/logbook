import { adoption, domCollector, pick, symbolizer } from "../bjsl/utilities.js";
import { AbstractTableCellPositioner } from "../bjsl/AbstractTableCellPositioner.js";
import { AbstractTableCellViewController } from "../bjsl/AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";
import InfiniteSheetHeader from "./InfiniteSheetHeader.js";
import { AbstractTableViewController } from "../bjsl/AbstactTableViewController.js";
import InfiniteSheetTask from "./InfiniteSheetTask.js";

const REGULAR_CELL_SYMBOL = symbolizer.symbolize("regularCellViewContainer");
const HEADER_CELL_SYMBOL = symbolizer.symbolize("headerCellViewContainer");

export class InfiniteSheet extends AbstractTableViewController {
    constructor() {
        super();

        /** @type {DataSource} */
        this.dataSource = undefined; // should be assigned by callee

        this.regularCellId = REGULAR_CELL_SYMBOL;
        this.headerCellId = HEADER_CELL_SYMBOL;

        Object.assign(this.config, {
            zoneOffsets: {
                preload: 0.3,
                parking: 0.4,
            },
            margins: {
                pageContent: {
                    before: 100,
                    after: 300,
                },
                [HEADER_CELL_SYMBOL]: {
                    before: 100,
                    between: 0,
                },
                [REGULAR_CELL_SYMBOL]: {
                    before: 10,
                    between: 5,
                },
            },
        });

        this.registerCellIdentifier(REGULAR_CELL_SYMBOL, () => {
            return new InfiniteSheetTask();
        });
        this.registerCellIdentifier(HEADER_CELL_SYMBOL, () => {
            return new InfiniteSheetHeader();
        });
    }

    addSectionHeadersToPlacement() {
        let lastSectionTimestamp = -1;
        for (let i = 0; i < this.placement; i++) {
            const taskSymbol = this.placement[i];
            const task = this.dataSource.cache.tasks.get(symbolizer.desymbolize(taskSymbol));
            const timestamp = task.createdAt;
            if (day(timestamp) !== lastTimestamp) {
                console.log("addSectionHeadersToPlacement");
            }
        }
    }

    /**
     * @param {Symbol} objectSymbol
     * @returns {number} expected height of the specified object in pixels
     */
    getDefaultHeightOfObject(objectSymbol) {
        return 32;
    }

    /** @returns {number} */
    getAverageHeightForAnObject() {
        return 32;
    }

    /** @returns {AbstractTableCellViewController} */
    getCellForObject(objectSymbol) {
        // TODO: variable cell type
        // const objectType = this.config.structuredDataMedium;

        let cellContainer;
        if (
            this.dataSource.cache.placements.chronological.headerSymbols.findIndex((symbol) => {
                return symbol === objectSymbol;
            }) != -1
        ) {
            cellContainer = this.requestReusableCellContainer(HEADER_CELL_SYMBOL);
            cellContainer.cell.setContent(this.dataSource.getTextContent(objectSymbol));
            // cellContainer.cell.container.dataset[]
        } else {
            cellContainer = this.requestReusableCellContainer(REGULAR_CELL_SYMBOL);
            cellContainer.cell.setContent(this.dataSource.getTextContent(objectSymbol));
            const details = this.dataSource.cache.tasks.get(objectSymbol);
            cellContainer.cell.setData({
                isCollaborated: details.isCollaborated,
                isTarget: details.isTarget,
                isCompleted: details.isCompleted,
            });
        }

        return cellContainer;
    }

    /** @param {Symbol} objectSymbol */
    getCellKindForObject(objectSymbol) {
        if (
            this.dataSource.cache.placements.chronological.headerSymbols.findIndex((symbol) => {
                return symbol === objectSymbol;
            }) != -1
        ) {
            return HEADER_CELL_SYMBOL;
        } else {
            return REGULAR_CELL_SYMBOL;
        }
    }

    /**
     * This function will be called for each cell that enters into the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellAppears(objectSymbol, cellPositioner) {
        // console.log(`${objectSymbol.toString()} has appeared.`);
    }

    /**
     * This function will be called for each cell that exits from the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellDisappears(objectSymbol, cellPositioner) {
        // console.log(`${objectSymbol.toString()} has disappeared.`);
    }

    /**
     * @param { Symbol } objectSymbol
     * @param { AbstractTableCellPositioner } cellContainer
     */
    updateCellIfNecessary(objectSymbol, cellContainer) {
        const newContent = this.dataSource.getTextContent(objectSymbol);
        cellContainer.cell.setContent(newContent);
    }

    deleteTask(taskId) {
        const ref = this.getReferenceOfAllocatedRowElement(0, 0);
        this.hideRowOnce(0, 1);
        this.state.effectiveOrdering[0].splice(1);
        this.calculateElementBounds();
        this.rePosition();
        // debugger;
    }
}
