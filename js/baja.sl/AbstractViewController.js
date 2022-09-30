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
}
