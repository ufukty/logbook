import * as constants from "../constants.js";

export class ModeSelector {
    constructor(delegate) {
        this.container = document.getElementById("mode-selector");
        this.delegate = delegate;

        this.state = {
            selectedMode: 0,
        };

        this.setListener();
    }

    setListener() {
        document.addEventListener("keyup", this.keyboardEventListener.bind(this));
        this.container.addEventListener("click", this.eventSwitchModes.bind(this));
    }

    keyboardEventListener(e) {
        if (e.ctrlKey && (e.key === "c" || e.key === "C")) {
            this.eventSwitchModes(e);
        }
    }

    eventSwitchModes(e) {
        e.preventDefault();

        this.state.selectedMode = 1 - this.state.selectedMode;

        this.updateView();
        this.notifyDelegate();
    }

    notifyDelegate() {
        if (this.delegate !== undefined && typeof this.delegate === "function") {
            this.delegate([constants.DVM_CHRONO, constants.DVM_HIERARCH][this.state.selectedMode]);
        }
    }

    updateView() {
        const selectedModeStringified = ["chronological", "hierarchical"][this.state.selectedMode];
        this.container.dataset.selectedMode = selectedModeStringified;
    }
}
