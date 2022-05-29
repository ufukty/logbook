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
} from "../utilities.js";

import AbstractViewController from "./AbstractViewController.js"

function createVC(viewControllerClassName) {
    const vc = new viewControllerClassName();
    return vc.container;
}



class RootViewController extends AbstractViewController {

    build() {

        const homeButton = createElement("a", {
            id: "home-button",
            classList: ["floating-corner", "left", "top"],
            href: "index.html"
        })

        // const modeSelectorView = createVC(DocumentViewModeSelector, {})


        const container = createElement("div", {
            classList: ["root"],
            children: [
                homeButton,
                // modeSelectorView,
            ]
        })

        return container;
    }
}

executeWhenDocumentIsReady(function () {
    const body = document.getElementsByTagName("body")[0];
    const rootVC = new RootViewController()
    body.appendChild(rootVC.container)
})