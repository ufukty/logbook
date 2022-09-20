import { AbstractViewController } from "./AbstractViewController.js";

export const POSITION_INSTANT = "POSITION_INSTANT";
export const POSITION_ANIMATE = "POSITION_ANIMATE";
export const POSITION_REDIRECT_IF_PLAYING = "POSITION_REDIRECT_IF_PLAYING";

export class AbstractTableCellPositioner extends AbstractViewController {
    constructor() {
        super();
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
        this.itemSymbol = undefined;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPosition(newPosition, withAnimation) {
        console.error("abstract function is called directly");
    }
}
