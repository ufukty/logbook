import { createElement } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";

export class AbstractTableCellViewController extends AbstractViewController {
    constructor() {
        super();
        this.task = createElement("div", ["task", "done"]);
        this.container = createElement("div", ["task-positioner"], [this.task]);
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        // TODO: clear data that are set by .setData()
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
    }

    setContent(newContent) {
        this.task.innerText = newContent;
    }

    setData(kv) {
        for (const key in kv) {
            if (Object.hasOwnProperty.call(kv, key)) {
                const value = kv[key];
                this.container.dataset[`${key}`] = value;
            }
        }
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
