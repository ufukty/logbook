import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
import { createElement, lerp, symbolizer, avg } from "./utilities.js";

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
            anim: {
                duration: 2000,
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
            console.error("Unexpected type");
            return;
        }

        const currentVerticalTranslation = lerp(
            this.state.ongoingAnimationParameters.translationStart,
            this.state.ongoingAnimationParameters.translationEnd,
            this.state.animation.effect.getComputedTiming().progress
        );
        const neededVerticalTranslation = newPosition - this.state.lastObjectPosition;
        this.state.animation.effect.setKeyframes([
            { transform: `translateY(${currentVerticalTranslation}px) scale(1)`, opacity: 1 },
            {
                transform: `translateY(${avg(currentVerticalTranslation, neededVerticalTranslation)}px) scale(0.96)`,
                opacity: 0.6,
            },
            { transform: `translateY(${neededVerticalTranslation}px) scale(1)`, opacity: 1 },
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
    setPositionY(newPosition, withAnimation = false) {
        // protection against object reuser change the content of cell
        this.state.objectSymbolAtAnimationStart = this.objectSymbol;

        const isAnimationOngoing = this.state.isAnimationOngoing;
        if (isAnimationOngoing && withAnimation) {
            this._setPositionRedirectOngoingAnimation(newPosition);
        } else if (isAnimationOngoing && !withAnimation) {
            this._setPositionByInterruptingOngoingAnimation(newPosition);
        } else if (!isAnimationOngoing && withAnimation) {
            this._setPositionWithAnimation(newPosition);
        } else if (!isAnimationOngoing && !withAnimation) {
            this._setPositionInstantly(newPosition);
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

        this.state.animation.cancel();
        this.container.style.top = `${this.state.ongoingAnimationParameters.positionAfterTransition}px`;
        this.state.lastObjectPosition = this.state.ongoingAnimationParameters.positionAfterTransition;

        delete this.state.animation;
        delete this.state.ongoingAnimationParameters;

        this.state.animation = undefined;
        this.state.isAnimationOngoing = false;
        this.state.ongoingAnimationParameters = undefined;
    }

    _animationAbortHandler() {
        this.state.isAnimationOngoing = false;
    }
}
