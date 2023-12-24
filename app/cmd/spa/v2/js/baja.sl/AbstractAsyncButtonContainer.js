import { createElement, symbolizer } from "./utilities.js";

const STATE_ACTIVE = symbolizer.symbolize("STATE_ACTIVE");
const STATE_PASSIVE = symbolizer.symbolize("STATE_PASSIVE");

const STATE_READY = symbolizer.symbolize("STATE_READY");
const STATE_PROCESSING = symbolizer.symbolize("STATE_PROCESSING");
const STATE_FAILURE = symbolizer.symbolize("STATE_FAILURE");
const STATE_SUCCESS = symbolizer.symbolize("STATE_SUCCESS");

export class AbstractAsyncButtonContainer extends AbstractViewController {
    constructor() {
        super();
        this.dom = Object.assign(this.dom, {
            container: createElement("div", ["baja-sl-two-state-button-container"]),
        });
        this.state = {
            active: STATE_ACTIVE,
            phase: STATE_READY,
        };
        this.dom.container.addEventListener("click", this._clickEventHandler.bind(this));

        this._classes = new Map([
            [STATE_ACTIVE, "baja-sl-button-phase-active"],
            [STATE_PASSIVE, "baja-sl-button-phase-passive"],
            [STATE_READY, "baja-sl-button-phase-ready"],
            [STATE_PROCESSING, "baja-sl-button-phase-processing"],
            [STATE_FAILURE, "baja-sl-button-phase-failure"],
            [STATE_SUCCESS, "baja-sl-button-phase-success"],
        ]);
    }

    /** @param {MouseEvent} e */
    _clickEventHandler(e) {
        if (!this.state.active) return;
    }

    clicked() {
        console.error("abstract function is not implemented");
    }

    unclicked() {
        console.error("abstract function is not implemented");
    }

    showFailure() {
        this.toggleContainerClassesAccordinglyToNextPhase(this.state.phase, STATE_FAILURE);
    }

    showSuccess() {
        this.toggleContainerClassesAccordinglyToNextPhase(this.state.phase, STATE_SUCCESS);
    }

    toggleContainerClassesAccordinglyToNextPhase(currentPhase, nextPhase) {
        this.dom.container.classList.remove(this._classes[currentPhase]);
        this.dom.container.classList.add(this._classes[nextPhase]);
    }
}
