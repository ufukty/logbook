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

class AbstractViewController {
    constructor() {
        this.container = createElement("div");
    }
}

export default AbstractViewController;