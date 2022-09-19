import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement, lerp, symbolizer, avg } from "./utilities.js";

export const POSITION_ANIMATE = "POSITION_ANIMATE";
export const POSITION_INSTANT = "POSITION_INSTANT";
export const POSITION_REDIRECT_IF_PLAYING = "POSITION_REDIRECT_IF_PLAYING";

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
            ongoingAnimationParameters: {
                positionAfterTransition: 0,
                translationStart: 0,
                translationEnd: 0,
            },
            lastObjectPosition: 0,
        };

        this.config = {
            ...this.config,
            anim: {
                duration: 300,
                iterations: 1,
                fill: "forwards",
                easing: "ease-out",
            },
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
        this.container.style.top = "0px";
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
        this.state = {
            /** @type {Animation} */
            animation: undefined,
            isAnimationOngoing: false,
            /** @type {Symbol} */
            objectSymbolAtAnimationStart: undefined,
        };
        this.cell.prepareForUse();
    }

    // MARK: SetPosition & its handlers

    _setPositionWithAnimation(newPosition) {
        const neededVerticalTranslation = newPosition - this.state.lastObjectPosition;

        // prettier-ignore
        const keyframes = [
            { transform: `translateY(0px) scale(1)`, opacity: 1, }, 
            { transform: `translateY(${neededVerticalTranslation / 2}px) scale(0.96)` , opacity: 0.6},
            { transform: `translateY(${neededVerticalTranslation}px) scale(1)` , opacity: 1}
        ];
        this.state.ongoingAnimationParameters = {
            translationStart: 0,
            translationEnd: neededVerticalTranslation,
            positionAfterTransition: newPosition,
        };
        this._startTransition(keyframes);
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
        const neededVerticalTranslation = newPosition - this.state.lastObjectPosition;
        // prettier-ignore
        this.state.animation.effect.setKeyframes([
            { 
                transform: `translateY(${currentVerticalTranslation}px) scale(1)`, 
                opacity: 1 
            },
            {
                transform: `translateY(${avg(currentVerticalTranslation, neededVerticalTranslation)}px) scale(0.96)`,
                opacity: 0.6,
            },
            { 
                transform: `translateY(${neededVerticalTranslation}px) scale(1)`, 
                opacity: 1 
            },
        ]);

        this.state.ongoingAnimationParameters = {
            translationStart: currentVerticalTranslation,
            translationEnd: neededVerticalTranslation,
            positionAfterTransition: newPosition,
        };

        this.state.animation.currentTime = 0;
        this.state.animation.play();
    }

    _setPositionByInterruptingOngoingAnimation(newPosition) {
        this.state.animation.cancel();
        this.container.style.top = `${newPosition}px`;
        this.state.lastObjectPosition = newPosition;
    }

    _setPositionInstantly(newPosition) {
        this.container.style.top = `${newPosition}px`;
        this.state.lastObjectPosition = newPosition;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     * This method uses combination of "transform:translateY()" and "top:" properties.
     * Animations done with translateY when they are requested and either case in the
     * end element will only have "top:" prop.
     */
    setPositionY(newPosition, animation) {
        // protection against object reuser change the content of cell
        this.state.objectSymbolAtAnimationStart = this.objectSymbol;
        const symbol = symbolizer.desymbolize(this.objectSymbol);
        if (this.state.isAnimationOngoing) {
            if (animation === POSITION_ANIMATE || animation === POSITION_REDIRECT_IF_PLAYING) {
                console.log(`_setPositionRedirectOngoingAnimation(${newPosition}px) for: ${symbol}`);
                this._setPositionRedirectOngoingAnimation(newPosition);
            } else if (animation === POSITION_INSTANT) {
                console.log(`_setPositionByInterruptingOngoingAnimation(${newPosition}px) for: ${symbol}`);
                this._setPositionByInterruptingOngoingAnimation(newPosition);
            }
        } else {
            if (animation === POSITION_ANIMATE) {
                console.log(`_setPositionWithAnimation(${newPosition}px) for: ${symbol}`);
                this._setPositionWithAnimation(newPosition);
            } else if (animation === POSITION_REDIRECT_IF_PLAYING || animation === POSITION_INSTANT) {
                console.log(`_setPositionInstantly(${newPosition}px) for: ${symbol}`);
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
        if (this.state.objectSymbolAtAnimationStart !== this.objectSymbol) return;

        const targetPos = this.state.ongoingAnimationParameters.positionAfterTransition;
        this.state.animation.cancel();
        this.container.style.top = `${targetPos}px`;
        this.state.lastObjectPosition = targetPos;

        delete this.state.animation;
        delete this.state.ongoingAnimationParameters;

        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
        this.state.ongoingAnimationParameters = undefined;
    }

    _animationAbortHandler() {
        delete this.state.animation;
        delete this.state.ongoingAnimationParameters;

        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
        this.state.ongoingAnimationParameters = undefined;
    }
}
