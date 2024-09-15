import { timestampToLocalizedText } from "../utility/dateTime";

import React from "react";

class DayHeader extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            dateStamp: props.dateStamp,
            effectiveDepth: props.effectiveDepth,
        };
    }

    static getDerivedStateFromProps(props, state) {
        return {
            dateStamp: props.dateStamp,
            effectiveDepth: props.effectiveDepth,
        };
    }

    render() {
        var style = {
            transform:
                "translateX(calc(" + this.state.effectiveDepth + " * var(--infinite-sheet-pixels-for-each-shift)))",
        };
        return (
            <div className="chronological-dvm-compo day-header" style={style}>
                {timestampToLocalizedText(this.state.dateStamp)}
            </div>
        );
    }
}

export default DayHeader;
