// TODO: Update status + CSS ==> Active tasks

import React from "react";
// import { timestampToLocalizedText } from "../../../../utility/dateTime";

import "./Task.css";

function max(left, right) {
    if (left > right) {
        return left;
    } else {
        return right;
    }
}

const amountOfShiftInPixelsForTaskDepth = 40;
// const documentViewModeChronological = "chro";
const documentViewModeHierarchical = "hier";

function calculatePositionsToLeft(documentViewMode, depth) {
    if (documentViewMode === documentViewModeHierarchical) {
        return amountOfShiftInPixelsForTaskDepth * max(0, depth - 1);
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
            documentPositions: {
                top: props.posY,
                transform:
                    "translateX(" + calculatePositionsToLeft(props.documentViewMode, props.taskDetails.depth) + "px)",
            },
        };
    }
    static getDerivedStateFromProps(props, state) {
        // console.log(nextProps, currentState);
        return {
            taskDetails: props.taskDetails,
            documentViewMode: props.documentViewMode,
            documentPositions: {
                top: props.posY,
                transform:
                    "translateX(" + calculatePositionsToLeft(props.documentViewMode, props.taskDetails.depth) + "px)",
            },
        };
    }

    // Returns the height of element by pixels and by minding the responsive design constraints
    returnHeight() {
        // check if the component rendered. if so, return the actual height instead hard-coded
        return 40;
    }

    render() {
        // console.log("task-positioner render", this.state.documentViewMode);
        // var objectKey = "taskPositioner" + this.state.taskDetails.id;
        var style = {
            top: this.state.documentPositions.top,
            transform: this.state.documentPositions.transform,
        };
        var taskID = this.state.taskDetails.task_id;
        var taskDetails = this.state.taskDetails;
        return (
            <div className="task-positioner" style={style}>
                <Task key={taskID} taskDetails={taskDetails} documentViewMode={this.state.documentViewMode}></Task>
            </div>
        );
    }
}

class Task extends React.Component {
    constructor(props) {
        super(props);

        // console.log(props);
        var status;
        if (props.taskDetails.completed_at != null) {
            status = "done";
        } else if (props.taskDetails.ready_to_pick_up) {
            status = "ready-to-pick-up";
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
