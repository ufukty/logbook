import { adoption, domCollector, pick, symbolizer } from "../baja.sl/utilities.js";
import { AbstractTableCellPositioner } from "../baja.sl/AbstractTableCellPositioner.js";
import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";
import InfiniteSheetHeader from "./InfiniteSheetHeader.js";
import { AbstractTableViewController } from "../baja.sl/AbstactTableViewController.js";
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
                preload: 0.4,
                parking: 0.5,
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
     * @param {Symbol} itemSymbol
     * @returns {number} expected height of the specified object in pixels
     */
    getDefaultHeightOfItem(itemSymbol) {
        const kind = this.getCellKindForItem();
        if (kind === HEADER_CELL_SYMBOL) return 34.8;
        else if (kind === REGULAR_CELL_SYMBOL) return 32.8;
        else return 32.8;
    }

    /** @returns {number} */
    getAverageHeightForAnItem() {
        return 32;
    }

    /** @returns {AbstractTableCellViewController} */
    getCellForItem(itemSymbol) {
        // TODO: variable cell type
        // const objectType = this.config.structuredDataMedium;

        let cellContainer;
        if (
            this.dataSource.cache.placements.chronological.headerSymbols.findIndex((symbol) => {
                return symbol === itemSymbol;
            }) != -1
        ) {
            cellContainer = this.requestReusableCellContainer(HEADER_CELL_SYMBOL);
            cellContainer.cell.setContent(this.dataSource.getTextContent(itemSymbol));
            // cellContainer.cell.container.dataset[]
        } else {
            cellContainer = this.requestReusableCellContainer(REGULAR_CELL_SYMBOL);
            cellContainer.cell.setContent(this.dataSource.getTextContent(itemSymbol));
            const details = this.dataSource.cache.tasks.get(itemSymbol);
            cellContainer.cell.setData({
                isCollaborated: details.isCollaborated,
                isTarget: details.isTarget,
                isCompleted: details.isCompleted,
            });
        }

        return cellContainer;
    }

    /** @param {Symbol} itemSymbol */
    getCellKindForItem(itemSymbol) {
        if (
            this.dataSource.cache.placements.chronological.headerSymbols.findIndex((symbol) => {
                return symbol === itemSymbol;
            }) != -1
        ) {
            return HEADER_CELL_SYMBOL;
        } else {
            return REGULAR_CELL_SYMBOL;
        }
    }

    /** @param {Symbol} itemSymbol */
    cellAppears(itemSymbol) {
        // this._debug(`${objectSymbol.toString()} has placed.`);
    }

    /**@param {Symbol} itemSymbol */
    cellAppears(itemSymbol) {
        // this._debug(`${objectSymbol.toString()} has appeared.`);
    }

    /** @param {Symbol} itemSymbol */
    cellDisappears(itemSymbol) {
        // this._debug(`${objectSymbol.toString()} has disappeared.`);
    }

    /**
     * @param { Symbol } itemSymbol
     * @param { AbstractTableCellPositioner } cellContainer
     */
    updateCellIfNecessary(itemSymbol, cellContainer) {
        const newContent = this.dataSource.getTextContent(itemSymbol);
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

    /**
     * @param {Symbol} superTaskSymbol
     * @param {Array.<Symbol>} subTaskSymbols
     */
    foldTask(superTaskSymbol, subTaskSymbols) {
        subTaskSymbols.forEach((symbol) => {
            this.config.placement.ignore.add(symbol);
        });
    }
}
