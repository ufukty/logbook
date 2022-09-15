import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement, lerp, symbolizer } from "./utilities.js";

export class BasicTableCellPositioner extends AbstractTableCellPositioner {
    constructor() {
        super();

        this.container = createElement("div", ["abstract-cell-scroller-view-cell-positioner"]);

        this.state = {
            /** @type {Animation} */
            animation: undefined,
            isAnimationOngoing: false,
            /** @type {Symbol} */
            objectSymbolAtAnimationStart: undefined,
            animationProps: {
                start: 0,
                end: 0,
            },
        };

        this.config = {
            debug: true,
            anim: {
                duration: 200,
                iterations: 1,
                fill: "both", // TODO: forwards
                easing: "cubic-bezier(0.3, 0.1, 0.7, 0.9)",
            },
        };
    }

    _debug(...data) {
        if (this.config.debug) console.log(...data);
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        if (this.state.isAnimationOngoing) {
            this.state.animation.cancel();
            delete this.animation;
            this.state.isAnimationOngoing = false;
        }
        this.cell.prepareForFree();
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
        this.state = {
            /** @type {Animation} */
            animation: undefined,
            isAnimationOngoing: false,
            /** @type {Symbol} */
            objectSymbolAtAnimationStart: undefined,
            animationProps: {
                start: 0,
                end: 0,
            },
        };
        this.cell.prepareForUse();
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     * This method will translate the dom element to the specified position,
     * with animation or instantly as specified, if there is
     */
    setPositionY(newPosition, withAnimation = false) {
        this._debug("setPosition is called", newPosition, withAnimation, symbolizer.desymbolize(this.objectSymbol));

        // protection against object reuser change the content of cell
        this.state.objectSymbolAtAnimationStart = this.objectSymbol;

        const isAnimationOngoing = this.state.isAnimationOngoing;
        if (isAnimationOngoing && withAnimation) {
            this._debug("setPosition dispatch", symbolizer.desymbolize(this.objectSymbol), "1");
            //
            this.state.animation.pause();

            if (this.state.animation.effect instanceof KeyframeEffect) {
                this._debug(
                    this.state.animationProps.start,
                    this.state.animationProps.end,
                    this.state.animation.effect.getComputedTiming().progress
                );

                const currentPosition = lerp(
                    this.state.animationProps.start,
                    this.state.animationProps.end,
                    this.state.animation.effect.getComputedTiming().progress
                );

                this.state.animation.effect.setKeyframes([
                    { transform: `translateY(${currentPosition}px)` },
                    { transform: `translateY(${newPosition}px)` },
                ]);

                this.state.animationProps = {
                    start: currentPosition,
                    end: newPosition,
                };
            }
            this.state.animation.currentTime = 0;
            this.state.animation.play();
        } else if (isAnimationOngoing && !withAnimation) {
            console.error("setPosition dispatch", symbolizer.desymbolize(this.objectSymbol), "2");
            this.state.animation.cancel();
            this.container.style.transform = `translateY(${newPosition}px)`;
        } else if (!isAnimationOngoing && withAnimation) {
            this._debug("setPosition dispatch", symbolizer.desymbolize(this.objectSymbol), "3");
            const currentPosition = this.state.animationProps.end;
            const keyframes = [
                { transform: `translateY(${currentPosition}px)` },
                { transform: `translateY(${newPosition}px)` },
            ];
            this.state.animationProps = {
                start: currentPosition,
                end: newPosition,
            };
            this._startTransition(keyframes);
            //
        } else if (!isAnimationOngoing && !withAnimation) {
            this._debug("setPosition dispatch", symbolizer.desymbolize(this.objectSymbol), "4");
            this.container.style.transform = `translateY(${newPosition}px)`;
            this.state.animationProps.end = newPosition;
            //
        }
    }

    _startTransition(keyframes) {
        this.state.isAnimationOngoing = true;
        this.state.animation = this.container.animate(keyframes, this.config.anim);
        this.state.animation.finished
            .then(this._animationCompletionHandler.bind(this))
            .catch(this._animationAbortHandler.bind(this));
    }

    _animationCompletionHandler() {
        // if the cell collected and reassigned to another "item" meanwhile
        if (this.state.objectSymbolAtAnimationStart !== this.objectSymbol) return;

        this._debug(symbolizer.desymbolize(this.objectSymbol), "finished");

        this.state.animation.commitStyles();

        this._debug("finished completion:", this.state.animation.effect.getComputedTiming().progress);

        // TODO: Reset "from" keyframe to initial values
        // TODO: Remove fill mode, unset "transform", set "top"

        // this.state.animation.effect.updateTiming({ fill: "none" });

        // this.container.style.top = `${this.state.targetPosition}px`;

        delete this.state.animation;
        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
    }

    _animationAbortHandler() {
        this._debug(symbolizer.desymbolize(this.objectSymbol), "aborted");
        this.state.isAnimationOngoing = false;
    }
}
