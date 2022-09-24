import { createElement } from "../baja.sl/utilities.js";

import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";

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
