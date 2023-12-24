import { AbstractViewController } from "./AbstractViewController.js";

export class AbstractTwoStateButtonContainer extends AbstractViewController {
    constructor() {
        super();
        this.dom = Object.assign(this.dom, {
            container: createElement("div", ["baja-sl-two-state-button-container"]),
        });
        this.state = {
            clicked: false,
        };
        this.dom.container.addEventListener("click", this._clickEventHandler.bind(this));
    }

    /** @param {MouseEvent} e */
    _clickEventHandler(e) {
        console.log(e);
        this.state.clicked = !this.state.clicked;
        if (this.state.clicked) this.unclicked();
        else this.clicked();
    }

    clicked() {
        console.error("abstract function is not implemented");
    }

    unclicked() {
        console.error("abstract function is not implemented");
    }
}
