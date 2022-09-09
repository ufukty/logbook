import {
    adoption,
    assert,
    createElementInSVGNamespace,
    toggleAnimationWithClass,
    setStyleProperties,
    createAnObjectOfLists,
    setAttributes,
    domElementReuseCollector,
    createElement,
} from "./utilities.js";

export class AbstractViewController {
    constructor() {
        this.dom = {};
        this.state = {};
        this.config = {};
    }

    buildView() {
        console.error("Abstract function has called directly.");
    }
}
