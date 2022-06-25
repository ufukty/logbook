import { createElement } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";

export class AbstractTableCellViewController extends AbstractViewController {
    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation) {
        console.error("Abstract class method .setPositionY() is called directly.");
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionX(newPosition, withAnimation) {
        console.error("Abstract class method .setPositionX() is called directly.");
    }

    /**
     * @param {number} newOpacity
     * @param {boolean} withAnimation
     */
    setOpacity(newOpacity, withAnimation) {
        console.error("Abstract class method .setOpacity() is called directly.");
    }

    fold() {}

    unfold() {}
}
