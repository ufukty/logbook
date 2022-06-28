import { createElement } from "../utilities.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";

const LockKind = {
    Folding: Symbol("folding"),
    Unfolding: Symbol("unfolding"),
    MovingY: Symbol("movingY"),
    MovingX: Symbol("movingX"),
    Deleting: Symbol("deleting"),
    Creating: Symbol("creating"),
};

export class InfiniteSheetTableCellViewController extends AbstractTableCellViewController {
    constructor() {
        super();
        // this.container.addEventListener("transitionend", this.transitionEndEventListener.bind(this));
        this.task = createElement("div", ["task", "done"]);
        this.container = createElement("div", ["task-positioner"], [this.task]);

        this.resetState();
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        this.resetState();
    }

    resetState() {
        this.transitionLock = false;
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
    }

    fold(withAnimation) {
        this.container.style.opacity = "0";
    }

    unfold(withAnimation) {
        this.container.style.opacity = "1";
    }

    /**
     * @param {number} newOpacity
     * @param {boolean} withAnimation
     */
    setOpacity(newOpacity, withAnimation) {
        console.error("Abstract class method .setOpacity() is called directly.");
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
}
