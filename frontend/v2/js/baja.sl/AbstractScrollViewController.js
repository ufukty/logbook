import { AbstractViewController } from "./AbstractViewController.js";
import { createElement, iota } from "./utilities.js";

export const SCROLL_MODE_SCROLL = "scroll";
export const SCROLL_MODE_HIDE = "hidden";
export const SCROLL_MODE_SHOW = "visible";
export const SCROLL_MODE_AUTO = "auto";

export const VERTICAL = iota();
export const HORIZONTAL = iota();

export class AbstractScrollViewController extends AbstractViewController {
    constructor() {
        super();

        this.dom = {
            ...this.dom,
            container: createElement("div", ["baja-sl-scroll-view-controller"]),
        };

        this.config = {
            ...this.config,
            overflowMode: {
                vertical: SCROLL_MODE_AUTO,
                horizontal: SCROLL_MODE_AUTO,
            },
        };
    }

    /**
     * @param {Symbol} direction
     * @returns {number}
     */
    currentScroll(direction) {
        if (direction === VERTICAL) {
            return this.dom.container.scrollTop;
        } else {
            return this.dom.container.scrollLeft;
        }
    }

    updateView() {
        this.dom.container.style.overflowY = this.config.overflowMode.vertical;
        this.dom.container.style.overflowX = this.config.overflowMode.horizontal;
    }
}
