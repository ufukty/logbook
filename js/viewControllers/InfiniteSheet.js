import { adoption, domElementReuseCollector, pSymbol } from "../utilities.js";
import { InfiniteSheetTableCellViewController } from "./InfiniteSheetTableCellViewController.js";
import {
    AbstractTableViewCellContainerViewController,
    AbstractTableViewController,
} from "./AbstactTableViewController.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";
import InfiniteSheetHeader from "./InfiniteSheetHeader.js";

export class InfiniteSheet extends AbstractTableViewController {
    constructor() {
        super();

        /** @type {DataSource} */
        this.dataSource = undefined; // should be assigned by callee

        this.config.margins = {
            pageContent: {
                before: 100,
                after: 300,
            },
            section: {
                before: 100,
                between: 0,
            },
            row: {
                before: 50,
                between: 30,
            },
        };

        this.regularCellId = pSymbol.get("regularCellViewContainer");
        this.headerCellId = pSymbol.get("headerCellViewContainer");

        this.registerCellIdentifier(this.regularCellId, () => {
            return new InfiniteSheetTableCellViewController();
        });
        this.registerCellIdentifier(this.headerCellId, () => {
            return new InfiniteSheetHeader();
        });
    }

    /**
     * @param {Symbol} objectSymbol
     * @returns {number} expected height of the specified object in pixels
     */
    getDefaultHeightOfObject(objectSymbol) {
        return 25;
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
