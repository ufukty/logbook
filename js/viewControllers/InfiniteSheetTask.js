import { createElement } from "../utilities.js";
import AbstractViewController from "./AbstractViewController.js";

class InfiniteSheetTask extends AbstractViewController {
    constructor() {
        super();
        this.task = createElement("div", ["task", "done"]);
        this.container = createElement("div", ["task-positioner"], [this.task]);
    }

    setContent(newContent) {
        this.task.innerText = newContent;
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

export default InfiniteSheetTask;
