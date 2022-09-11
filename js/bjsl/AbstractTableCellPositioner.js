export class AbstractTableCellPositioner {
    constructor() {
        /** @type {HTMLElement} */
        this.container = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type { AbstractTableCellViewController }
         */
        this.cell = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type { Symbol }
         */
        this.reuseIdentifier = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type {Symbol}
         */
        this.objectSymbol = undefined;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionX(newPosition, withAnimation) {
        console.error("abstract function is called directly");
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation) {
        console.error("abstract function is called directly");
    }
}
