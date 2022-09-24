import { AbstractViewController } from "./AbstractViewController.js";

export class AbstractTableCellViewController extends AbstractViewController {
    constructor() {
        super();
    }

    prepareForFree() {
        console.error("abstract function is called directly");
    }

    prepareForUse() {
        console.error("abstract function is called directly");
    }

    setContent(newContent) {
        console.error("abstract function is called directly");
    }

    setData(kv) {
        for (const key in kv) {
            if (Object.hasOwnProperty.call(kv, key)) {
                const value = kv[key];
                this.container.dataset[`${key}`] = value;
            }
        }
    }

    fold() {
        console.error("abstract function is called directly");
    }

    unfold() {
        console.error("abstract function is called directly");
    }
}
