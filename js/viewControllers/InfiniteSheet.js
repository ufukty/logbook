import { adoption, domCollector, symbolizer } from "../bjsl/utilities.js";
import { InfiniteSheetTableCellViewController } from "./InfiniteSheetTableCellViewController.js";
import {
    AbstractTableViewCellContainerViewController,
    AbstractTableViewController,
} from "../bjsl/AbstactTableViewController.js";
import { AbstractTableCellViewController } from "../bjsl/AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";
import InfiniteSheetHeader from "./InfiniteSheetHeader.js";

export class InfiniteSheet extends AbstractTableViewController {
    constructor() {
        super();

        /** @type {DataSource} */
        this.dataSource = undefined; // should be assigned by callee

        const regularCellId = symbolizer.symbolize("regularCellViewContainer");
        const headerCellId = symbolizer.symbolize("headerCellViewContainer");

        this.regularCellId = regularCellId;
        this.headerCellId = headerCellId;

        this.config.margins = {
            pageContent: {
                before: 100,
                after: 300,
            },
            headerCellId: {
                before: 100,
                between: 0,
            },
            regularCellId: {
                before: 50,
                between: 30,
            },
        };

        this.registerCellIdentifier(this.regularCellId, () => {
            return new InfiniteSheetTableCellViewController();
        });
        this.registerCellIdentifier(this.headerCellId, () => {
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
        return 31;
    }

    /** @returns {number} */
    getAverageHeightForAnObject() {
        return 31;
    }

    /** @returns {AbstractTableCellViewController} */
    getCellForObject(objectSymbol) {
        // TODO: variable cell type
        // const objectType = this.config.structuredDataMedium;

        let cellContainer;
        if (this.dataSource.medium.data.rows.has(objectSymbol)) {
            cellContainer = this.requestReusableCellContainer(this.headerCellId);
            cellContainer.cell.setContent(this.dataSource.getTextContent(objectSymbol));
        } else {
            cellContainer = this.requestReusableCellContainer(this.regularCellId);
            cellContainer.cell.setContent(this.dataSource.getTextContent(objectSymbol));
            // cellContainer.cell.container.dataset[]
        }

        return cellContainer;
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
     * @param { AbstractTableViewCellContainerViewController } cellContainer
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
