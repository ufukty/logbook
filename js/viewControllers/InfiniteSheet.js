import { adoption, domElementReuseCollector, pSymbol } from "../utilities.js";
import { InfiniteSheetTableCellViewController } from "./InfiniteSheetTableCellViewController.js";
import { AbstractTableViewController } from "./AbstactTableViewController.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";

export class InfiniteSheet extends AbstractTableViewController {
    constructor() {
        super();

        /** @type { DataSource } */
        this.dataSource = undefined; // should be assigned by callee

        this.regularCellId = pSymbol.get("regularCellViewContainer");

        this.config.margins = {
            pageContent: {
                before: 100,
                after: 300,
            },
            section: {
                before: 100,
                between: 100,
            },
            row: {
                before: 100,
                between: 20,
            },
        };

        this.registerCellConstructor(this.regularCellId, () => {
            return new InfiniteSheetTableCellViewController();
        });
    }

    /**
     * @param {Symbol} objectSymbol
     * @returns {number} expected height of the specified object in pixels
     */
    getDefaultHeightOfObject(objectSymbol) {
        return 25;
    }

    /** @returns { AbstractTableCellViewController } */
    getCellForObjectId(objectSymbol) {
        // TODO: variable cell type
        // const objectType = this.config.structuredDataMedium;

        /** @type { AbstractTableCellViewController } */
        const cell = this.getRecycledCell(this.regularCellId);
        cell.setContent(this.dataSource.getTextContent(objectSymbol));
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
