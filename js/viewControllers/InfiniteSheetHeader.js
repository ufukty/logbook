import { createElement } from "../bjsl/utilities.js";

import { AbstractTableCellViewController } from "../bjsl/AbstractTableCellViewController.js";

class InfiniteSheetHeader extends AbstractTableCellViewController {
    constructor() {
        super();
        this.container = createElement("div", ["infinite-sheet-header"]);
    }

    prepareForUse() {}

    prepareForFree() {
        this.container.innerText = "";
    }

    setContent(title) {
        this.container.innerText = title;
    }
}

export default InfiniteSheetHeader;
