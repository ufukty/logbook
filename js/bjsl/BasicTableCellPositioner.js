import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement, symbolizer } from "./utilities.js";

export class BasicTableCellPositioner extends AbstractTableCellPositioner {
    constructor() {
        super();
        this.container = createElement("div", ["abstract-cell-scroller-view-cell-positioner"]);
        this.cell = undefined; // should be assigned by callee

        // TODO: Reset "from" keyframe to initial values
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
        console.log("setPosition is called", symbolizer.desymbolize(this.objectSymbol), newPosition, withAnimation);
        if (withAnimation) {
            // TODO: if there is animation currently unfinished => abort

            const objectSymbolAtAnimationStart = this.objectSymbol;

            const currentPos = this.container.style.top.match(/(\d+)/)[1];
            const deltaPos = newPosition - currentPos;
            // console.log(currentPos);

            // prettier-ignore
            const keyframes = [
                { transform: `translateY(0px)` }, // TODO: get starting keyframe from "this"
                { transform: `translateY(${deltaPos}px)` }
            ]
            const config = {
                duration: 2000,
                iterations: 1,
                fill: "none", // TODO: forwards
                easing: "cubic-bezier(0.3, 0.1, 0.7, 0.9)",
            };
            this.animation = this.container.animate(keyframes, config);
            // this.animation.commitStyles();
            this.animation.finished
                .then(() => {
                    console.log(symbolizer.desymbolize(this.objectSymbol), "finished");
                    // this.animation.commitStyles();
                    // TODO: Reset "from" keyframe to initial values
                    // TODO: Remove fill mode, unset "transform", set "top"

                    // this.animation.effect.updateTiming({ fill: none });
                    if (objectSymbolAtAnimationStart === this.objectSymbol) {
                        this.container.style.top = `${newPosition}px`;
                    }
                })
                .catch(() => {
                    // console.log(this.animation.effect.getTiming());
                    // console.log(this.animation.effect.getComputedTiming());
                    console.log(symbolizer.desymbolize(this.objectSymbol), "aborted");

                    // set "from" keyframe to last values
                });
        } else {
            if (this.animation) {
                this.animation.cancel();
                this.animation = undefined;
            }
            console.log(symbolizer.desymbolize(this.objectSymbol), "w/o animation");
            this.container.style.top = `${newPosition}px`;
        }
    }
}
