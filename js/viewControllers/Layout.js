import { Size } from "../baja.sl/Layout/Coordinates.js";
import { AbstractViewController } from "../baja.sl/AbstractViewController.js";
import { AbstractLayoutDecorator } from "./Layout/AbstractLayoutCalculator.js";

export class PlaceholderViewController extends AbstractViewController {
    constructor() {
        super();

        this.dom = {
            ...this.dom,
            container: createElement("div", ["placeholder"]),
        };

        this.config = {
            ...this.config,
            /** @type {Size} */
            size: undefined,
        };
    }

    updateView() {
        this.dom.container.style.width = this.config.size.width;
        this.dom.container.style.height = this.config.size.height;
    }
}

export class AvatarLayout extends AbstractLayoutDecorator {
    constructor() {
        super();
    }
}
