import { AbstractViewController } from "./AbstractViewController.js";
import { Position } from "./Layout/Coordinates.js";
import { adoption, createElement, iota } from "./utilities.js";

/** first level of presentation */
const PRESENTATION_STATE_PLACEHOLDER = iota();
/** second level of presentation */
const PRESENTATION_STATE_SUMMARY = iota();
/** third level of presentation */
const PRESENTATION_STATE_DETAILED = iota();

export class AbstractManagedLayoutCellViewController extends AbstractViewController {
    constructor() {
        super();

        this.dom = {
            ...this.dom,
            /**
             * @private
             * @type {HTMLElement}
             * This will be used to position user-provided HTMLElement without
             *   utilizing its transform/position props.
             */
            managedLayoutPositioner: createElement("div", ["baja-sl-managed-layout-cell-view-positioner"]),
            /**
             * @type {HTMLElement}
             * Concrete classes of this AbstractManagedLayoutCellViewController
             *   should append their HTMLElements to this (directly or
             *   indirectly) element. They also can edit content and style of
             *   this element as they wish.
             */
            container: createElement("div", ["baja-sl-managed-layout-cell-view-container"]),
        };

        this.dom.managedLayoutPositioner.style.position = "absolute";

        // prettier-ignore
        adoption(this.dom.managedLayoutPositioner, [
            adoption(this.dom.container)
        ])

        this.state = this._getStateTemplate();

        this.config = {
            ...this.config,
            /** @type {Symbol} */
            assignedItemSymbol: undefined,
            animation: {
                translocation: {
                    duration: 300,
                    iterations: 1,
                    fill: "forwards",
                    easing: "cubic-bezier(0.2, 0.1, 0.4, 0.90)",
                    speed: 0.02, // pixels per millisecond
                },
            },
            leveledPresentation: {
                timeoutDuration: {
                    secondLevelOfPresentation: 15 * (1000 / 60),
                    thirdLevelOfPresentation: 15 * (1000 / 60),
                },
            },
        };

        /** @private */
        this.abstract = {
            timeouts: {
                /** @type {number} */
                secondLevelOfPresentation: undefined,
                /** @type {number} */
                thirdLevelOfPresentation: undefined,
            },
        };
    }

    /** @private */
    _getStateTemplate() {
        return {
            currentPresentationLevel: PRESENTATION_STATE_PLACEHOLDER,
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
        this.dom.managedLayoutPositioner.style.display = "none";
    }

    async prepareForFreeAsync() {
        if (this.state.isAnimationOngoing) {
            this.state.animation.cancel();
            delete this.animation;
            this.state.isAnimationOngoing = false;
        }
        // this.cell.prepareForFree();
        this.state = this._getStateTemplate();
        this.dom.managedLayoutPositioner.style.top = "0px";

        this._cancelLeveledPresentation();
    }

    _cancelLeveledPresentation() {
        for (const [timeoutName, timeoutNumber] of Object.keys(this.abstract.timeouts).entries()) {
            if (timeoutNumber) {
                this.abstract.timeouts[timeoutName] = undefined;
                clearTimeout(timeoutNumber);
            }
        }
    }

    prepareForUse() {
        this.dom.managedLayoutPositioner.style.display = "block";
        this.leveledPresentation(PRESENTATION_STATE_PLACEHOLDER);
    }

    /** @private */
    leveledPresentation(level) {
        const setupTimeItemSymbol = this.config.assignedItemSymbol;

        if (level === PRESENTATION_STATE_PLACEHOLDER) {
            this.firstLevelOfPresentation();
            this.abstract.timeouts.secondLevelOfPresentation = setTimeout(() => {
                if (setupTimeItemSymbol === this.config.assignedItemSymbol) {
                    this.leveledPresentation(PRESENTATION_STATE_SUMMARY);
                }
            }, this.config.leveledPresentation.timeoutDuration.secondLevelOfPresentation);
        } else if (level === PRESENTATION_STATE_SUMMARY) {
            this.secondLevelOfPresentation();
            this.abstract.timeouts.thirdLevelOfPresentation = setTimeout(() => {
                if (setupTimeItemSymbol === this.config.assignedItemSymbol) {
                    this.leveledPresentation(PRESENTATION_STATE_DETAILED);
                }
            }, this.config.leveledPresentation.timeoutDuration.thirdLevelOfPresentation);
        } else if (level === PRESENTATION_STATE_DETAILED) {
            this.thirdLevelOfPresentation();
        }
    }

    /** @abstract */
    firstLevelOfPresentation() {}

    /** @abstract */
    secondLevelOfPresentation() {}

    /** @abstract */
    thirdLevelOfPresentation() {}

    /** @param {Position} newPosition */
    setPosition(newPosition) {
        this.dom.managedLayoutPositioner.style.top = `${newPosition.y}px`;
        this.dom.managedLayoutPositioner.style.left = `${newPosition.x}px`;
    }
}
