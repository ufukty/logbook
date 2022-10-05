import { Size } from "./Coordinates.js";

export class AbstractViewController {
    constructor() {
        this.dom = {
            /** @type {HTMLElement} */
            container: undefined,
        };
        this.state = {};
        this.config = {};
    }

    /** @param {AbstractViewController} placeholder */
    exportSubview(placeholder) {
        console.error("abstract function is not implemented");
    }

    /** @param {Symbol} trigger */
    updateView(trigger) {
        console.error("abstract function is not implemented");
    }

    measureSize() {
        const computedStyle = getComputedStyle(this.dom.container);
        const computedHeight = parseFloat(computedStyle.getPropertyValue("height"));
        const computedWidth = parseFloat(computedStyle.getPropertyValue("width"));
        return new Size(computedWidth, computedHeight);
    }
}
