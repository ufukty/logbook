// TODO: Update status + CSS ==> Active tasks

import React from "react";
// import { timestampToLocalizedText } from "../../../../utility/dateTime";

import "./Task.css";
import * as constants from "../../../../constants";
import { calculateShiftForItem } from "../../../../AutoFocusManager";

function max(left, right) {
    if (left > right) {
        return left;
    } else {
        return right;
    }
}

function calculatePositionsToLeft(documentViewMode, adjustmentForDepth, depth) {
    if (documentViewMode === constants.DVM_HIERARCH) {
        // console.log(calculateShiftForItem(adjustmentForDepth, depth));
        return calculateShiftForItem(adjustmentForDepth, depth);
    } else {
        return 0;
    }
}

class TaskPositioner extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            taskDetails: props.taskDetails,
            documentViewMode: props.documentViewMode,
            adjustmentForDepth: props.adjustmentForDepth,
            horizontalShift: calculatePositionsToLeft(
                props.documentViewMode,
                props.adjustmentForDepth,
                props.taskDetails.depth
            ),
        };
        this.div = React.createRef();
    }

    static getDerivedStateFromProps(props, state) {
        // console.log(nextProps, currentState);
        return {
            taskDetails: props.taskDetails,
            documentViewMode: props.documentViewMode,
            adjustmentForDepth: props.adjustmentForDepth,
            horizontalShift: calculatePositionsToLeft(
                props.documentViewMode,
                props.adjustmentForDepth,
                props.taskDetails.depth
            ),
        };
    }

    render() {
        // console.log("task-positioner render");
        var style = {
            transform: "translateX(" + this.state.horizontalShift + "px)",
        };
        var taskID = this.state.taskDetails.task_id;
        var taskDetails = this.state.taskDetails;
        return (
            <div ref={this.div} className="task-positioner" style={style}>
                <Task key={taskID} taskDetails={taskDetails} documentViewMode={this.state.documentViewMode}></Task>
            </div>
        );
    }
}

class Task extends React.Component {
    constructor(props) {
        super(props);
        var status;
        if (props.taskDetails.completed_at != null) {
            status = "done";
        } else {
            status = "pending";
        }
        this.state = {
            taskDetails: props.taskDetails,
            taskStatus: status,
            documentViewMode: props.documentViewMode,
        };
    }

    render() {
        // console.log("task render");
        return (
            <div
                task_id={this.state.taskDetails.task_id}
                task_depth={this.state.taskDetails.depth}
                className={"task " + this.state.taskStatus}
                // draggable="true"
                contentEditable="true"
                suppressContentEditableWarning={true}
            >
                {this.state.taskDetails.content}
            </div>
        );
    }
}

export default TaskPositioner;
