import { createElement } from "../bjsl/utilities.js";

import { AbstractTableCellViewController } from "../bjsl/AbstractTableCellViewController.js";

class InfiniteSheetHeader extends AbstractTableCellViewController {
    constructor() {
        super();
        this.header = createElement("div", ["header"]);
        this.container = createElement("div", ["header-positioner"], [this.header]);
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
    }

    setContent(title) {
        this.header.innerText = title;
    }
}

export default InfiniteSheetHeader;
