import React from "react";

import * as constants from "../utility/constants";

class ModeSelector extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedMode: 0,
            documentViewModeChangeDelegate: props.documentViewModeChangeDelegate,
        };
        document.addEventListener("keyup", this.keyboardEventListener.bind(this));
    }

    keyboardEventListener(e) {
        if (e.ctrlKey && (e.key === "c" || e.key === "C")) {
            this.eventSwitchModes(e);
        }
    }

    eventSwitchModes(e) {
        e.preventDefault();
        var changeModeTo = 1 - this.state.selectedMode;
        this.setState({
            selectedMode: changeModeTo,
        });
        var delegate = this.state.documentViewModeChangeDelegate;
        delegate([constants.DVM_CHRONO, constants.DVM_HIERARCH][changeModeTo]);
    }

    render() {
        var classNameForWrapper = ["left-picked", "right-picked"][this.state.selectedMode];
        var eventSwitchModes = this.eventSwitchModes.bind(this);
        return (
            <div id="settings-documentViewMode" className={classNameForWrapper} onClick={eventSwitchModes}>
                <div id="left">C</div>
                <div id="right">H</div>
                <div id="left-activated-caption">Chronological</div>
                <div id="right-activated-caption">Hierarchical</div>
            </div>
        );
    }
}

export default ModeSelector;