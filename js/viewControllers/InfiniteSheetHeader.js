import { createElement } from "../utilities.js";

import AbstractViewController from "./AbstractViewController.js";

class InfiniteSheetHeader extends AbstractViewController {
    constructor() {
        super();
        this.header = createElement("div", ["header"]);
        this.container = createElement("div", ["header-positioner"], [this.header]);
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
    }

    setContent(title) {
        this.header.innerText = title;
    }

    setPosition(posY) {
        this.container.style.top = `${posY}px`;
    }

    setData(kv) {
        for (const key in kv) {
            if (Object.hasOwnProperty.call(kv, key)) {
                const value = kv[key];
                this.container.dataset[`${key}`] = value;
            }
        }
    }
}

export default InfiniteSheetHeader;
