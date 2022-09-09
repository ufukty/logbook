import {
    adoption,
    assert,
    createElementInSVGNamespace,
    toggleAnimationWithClass,
    setStyleProperties,
    createAnObjectOfLists,
    setAttributes,
    domElementReuseCollector,
    executeWhenDocumentIsReady,
    createElement,
} from "../bjsl/utilities.js";

import AbstractViewController from "./AbstractViewController.js";

class DocumentSheet extends AbstractViewController {
    build() {
        console.log("hello world");
        return createElement("div", ["document-sheet"]);
    }
}

export default DocumentSheet;
