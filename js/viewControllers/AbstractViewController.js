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
} from "../utilities.js";

export class AbstractViewController {
    constructor() {
        this.container = undefined; // createElement("div");
    }
}
