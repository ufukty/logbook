import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement } from "./utilities.js";

export class BasicTableCellPositioner extends AbstractTableCellPositioner {
    constructor() {
        super();
        this.container = createElement("div", ["abstract-cell-scroller-view-cell-positioner"]);
        this.cell = undefined; // should be assigned by callee
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        if (this.animation) {
            this.animation.cancel();
            delete this.animation;
        }
        this.cell.prepareForFree();
        // this.container.style.top = `0px`;
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
        this.cell.prepareForUse();
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation = false) {
        if (withAnimation) {
            const objectSymbolAtAnimationStart = this.objectSymbol;

            const currentPos = this.container.style.top.match(/\d+/);
            const deltaPos = newPosition - currentPos;

            // prettier-ignore
            const keyframes = [
                { transform: "translateY(0px)" }, 
                { transform: `translateY(${deltaPos}px)` }
            ]
            const config = {
                duration: 200,
                iterations: 1,
                fill: "none",
                // easing: "cubic-bezier(0.5, 1.2, 0.8, 1.0)",
            };
            this.animation = this.container.animate(keyframes, config);
            this.animation.finished
                .then(() => {
                    if (objectSymbolAtAnimationStart === this.objectSymbol) {
                        this.container.style.top = `${newPosition}px`;
                    }
                })
                .catch(() => {
                    // console.log("animation is aborted");
                });
        } else {
            this.container.style.top = `${newPosition}px`;
        }
    }
}
