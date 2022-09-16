import { createElement } from "../bjsl/utilities.js";
import { AbstractTableCellViewController } from "../bjsl/AbstractTableCellViewController.js";

class InfiniteSheetTask extends AbstractTableCellViewController {
    constructor() {
        super();
        this.container = createElement("div", ["infinite-sheet-task", "done"]);
    }

    prepareForFree() {
        this.container.innerText = "";
    }

    prepareForUse() {}

    setContent(newContent) {
        this.container.innerText = newContent;
    }
}

export default InfiniteSheetTask;
