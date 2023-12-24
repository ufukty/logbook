import { createElement } from "../baja.sl/utilities.js";
import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";

export class InfiniteSheetHeader extends AbstractTableCellViewController {
    constructor() {
        super();
        this.dom = {
            container: createElement("div", ["infinite-sheet-header"]),
        };
    }

    prepareForUse() {}

    prepareForFree() {
        this.dom.container.innerText = "";
    }

    setContent(title) {
        this.dom.container.innerText = title;
    }
}
