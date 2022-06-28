import { createElement } from "../utilities.js";
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
        console.error("abstract function is called directly");
    }

    fold() {
        console.error("abstract function is called directly");}

    unfold() {
        console.error("abstract function is called directly");}
}
