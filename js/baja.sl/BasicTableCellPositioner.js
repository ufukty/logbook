import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement, lerp, symbolizer, avg } from "./utilities.js";

export class BasicTableCellPositioner extends AbstractTableCellPositioner {
    constructor() {
        super();

        this.container = createElement("div", ["abstract-cell-scroller-view-cell-positioner"]);

        this.state = this._stateTemplate();

        this.config = {
            ...this.config,
            anim: {
                duration: 300,
                iterations: 1,
                fill: "forwards",
                easing: "ease-in-out",
            },
        };
    }

    _stateTemplate() {
        return {
            /** @type {Animation} */
            animation: undefined,
            isAnimationOngoing: false,
            /** @type {Symbol} */
            itemSymbolAtAnimationStart: undefined,
            ongoingAnimationParameters: {
                optimizeForEnding: false,
                positionAfterTransition: 0,
                translationStart: 0,
                translationEnd: 0,
            },
            lastCellPosition: 0,
            callback: undefined,
        };
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        if (this.state.isAnimationOngoing) {
            this.state.animation.cancel();
            delete this.animation;
            this.state.isAnimationOngoing = false;
        }
        this.cell.prepareForFree();
        this.state = this._stateTemplate();
        this.container.style.top = "0px";
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
        this.cell.prepareForUse();
    }

    // MARK: SetPosition & its handlers

    _setPositionWithAnimation(newPosition) {
        const neededVerticalTranslation = newPosition - this.state.lastCellPosition;

        const keyframes = [
            { transform: `translateY(0px)` },
            { transform: `translateY(${neededVerticalTranslation}px)` },
        ];
        this.state.ongoingAnimationParameters = {
            translationStart: 0,
            translationEnd: neededVerticalTranslation,
            positionAfterTransition: newPosition,
        };
        this._startTransition(keyframes);
    }

    _setPositionRedirectOngoingAnimationWithExpandingDuration(newPosition) {
        this.state.animation.pause();

        if (!this.state.animation.effect instanceof KeyframeEffect) {
            this._error("Unexpected type");
            return;
        }

        const currentVerticalTranslation = lerp(
            this.state.ongoingAnimationParameters.translationStart,
            this.state.ongoingAnimationParameters.translationEnd,
            this.state.animation.effect.getComputedTiming().progress
        );
        const neededVerticalTranslation = newPosition - this.state.lastCellPosition;
        this.state.animation.effect.setKeyframes([
            { transform: `translateY(${currentVerticalTranslation}px)` },
            { transform: `translateY(${neededVerticalTranslation}px)` },
        ]);

        this.state.ongoingAnimationParameters = {
            translationStart: currentVerticalTranslation,
            translationEnd: neededVerticalTranslation,
            positionAfterTransition: newPosition,
        };

        this.state.animation.currentTime = 0;
        this.state.animation.play();
    }

    _setPositionRedirectOngoingAnimation(newPosition) {
        this.state.animation.pause();

        if (!this.state.animation.effect instanceof KeyframeEffect) {
            this._error("Unexpected type");
            return;
        }

        const currentVerticalTranslation = lerp(
            this.state.ongoingAnimationParameters.translationStart,
            this.state.ongoingAnimationParameters.translationEnd,
            this.state.animation.effect.getComputedTiming().progress
        );
        const neededVerticalTranslation = newPosition - this.state.lastCellPosition;
        this.state.animation.effect.setKeyframes([
            { transform: `translateY(${currentVerticalTranslation}px)` },
            { transform: `translateY(${neededVerticalTranslation}px)` },
        ]);

        this.state.ongoingAnimationParameters = {
            translationStart: currentVerticalTranslation,
            translationEnd: neededVerticalTranslation,
            positionAfterTransition: newPosition,
        };
        this.state.animation.play();
    }

    _setPositionByInterruptingOngoingAnimation(newPosition) {
        this.state.animation.cancel();
        this.container.style.top = `${newPosition}px`;
        this.state.lastCellPosition = newPosition;
    }

    _setPositionInstantly(newPosition) {
        this.container.style.top = `${newPosition}px`;
        this.state.lastCellPosition = newPosition;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} startNewAnimation
     * @param {boolean} stopOngoingAnimation
     * @param {function():void} callback
     * This method uses combination of "transform:translateY()" and "top:" properties.
     * Animations done with translateY when they are requested and either case in the
     * end element will only have "top:" prop.
     */
    setPosition(newPosition, startNewAnimation = false, stopOngoingAnimation = false, callback = undefined) {
        // protection against object reuser change the content of cell
        if (callback) this.state.callback = callback;

        this.state.itemSymbolAtAnimationStart = this.itemSymbol;
        const symbol = symbolizer.desymbolize(this.itemSymbol);
        if (this.state.isAnimationOngoing) {
            if (startNewAnimation) {
                // console.log(
                //     `_setPositionRedirectOngoingAnimationWithExpandingDuration(${newPosition}px) for: ${symbol}`
                // );
                this._setPositionRedirectOngoingAnimationWithExpandingDuration(newPosition);
            } else if (stopOngoingAnimation) {
                this._setPositionByInterruptingOngoingAnimation(newPosition);
            } else {
                // console.log(`_setPositionRedirectOngoingAnimation(${newPosition}px) for: ${symbol}`);
                this._setPositionRedirectOngoingAnimation(newPosition);
            }
        } else {
            if (startNewAnimation) {
                // console.log(`_setPositionWithAnimation(${newPosition}px) for: ${symbol}`);
                this._setPositionWithAnimation(newPosition);
            } else {
                // console.log(`_setPositionInstantly(${newPosition}px) for: ${symbol}`);
                this._setPositionInstantly(newPosition);
            }
        }
    }

    // MARK: Animation Handlers

    _startTransition(keyframes) {
        this.state.isAnimationOngoing = true;
        this.state.animation = this.container.animate(keyframes, this.config.anim);
        this.state.animation.finished
            .then(this._animationCompletionHandler.bind(this))
            .catch(this._animationAbortHandler.bind(this));
    }

    _animationCompletionHandler() {
        // if the cell collected and reassigned to another "item" meanwhile
        if (this.state.itemSymbolAtAnimationStart !== this.itemSymbol) return;

        const targetPos = this.state.ongoingAnimationParameters.positionAfterTransition;
        this.state.animation.cancel();
        this.container.style.top = `${targetPos}px`;
        this.state.lastCellPosition = targetPos;

        delete this.state.animation;
        delete this.state.ongoingAnimationParameters;

        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
        this.state.ongoingAnimationParameters = undefined;

        if (this.state.callback) {
            const callback = this.state.callback;
            this.state.callback = undefined;
            callback();
        }
    }

    _animationAbortHandler() {
        // if the cell collected and reassigned to another "item" meanwhile
        if (this.state.itemSymbolAtAnimationStart !== this.itemSymbol) return;

        delete this.state.animation;
        delete this.state.ongoingAnimationParameters;

        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
        this.state.ongoingAnimationParameters = undefined;

        if (this.state.callback) {
            const callback = this.state.callback;
            console.log(callback);
            this.state.callback = undefined;
            callback();
        }
    }
}
