import { adoption, createElement } from "../baja.sl/utilities.js";
import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";

export class InfiniteSheetTask extends AbstractTableCellViewController {
    constructor() {
        super();
        this.dom = {
            container: createElement("div", ["infinite-sheet-task", "done"]),
            editableArea: createElement("span", ["infinite-sheet-task-editable-area"]),
        };
        // prettier-ignore
        adoption(this.dom.container, [
            adoption(this.dom.editableArea)
        ])
        // this.dom.editableArea.contentEditable = true;

        this.config = {
            ...this.config,
            translationForDepth: 20,
        };
    }

    prepareForFree() {
        this.dom.editableArea.innerText = "";
    }

    prepareForUse() {}

    setContent(newContent) {
        this.dom.editableArea.innerText = newContent;
    }

    setDegree(degree) {}

    setDepth(depth) {
        this.dom.container.style.transform = `translateX(${this.config.translationForDepth * depth}px)`;
    }

    enableEditMode() {
        this.dom.editableArea.contentEditable = true;
        this.dom.editableArea.focus();
    }
}
