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
        // console.log("task-positioner calculatePositionsToLeft defined");
        return amountOfShiftInPixelsForTaskDepth * max(0, depth - 1);
    } else {
        // console.log("task-positioner calculatePositionsToLeft undefined");
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
                left: calculatePositionsToLeft(props.documentViewMode, props.taskDetails.depth),
            },
            sizeChangeDelegate: props.sizeUpdateHandler,
        };
        this.div = React.createRef();
    }

    static getDerivedStateFromProps(props, state) {
        // console.log(nextProps, currentState);
        return {
            taskDetails: props.taskDetails,
            documentViewMode: props.documentViewMode,
            documentPositions: {
                top: props.posY,
                left: calculatePositionsToLeft(props.documentViewMode, props.taskDetails.depth),
            },
            sizeChangeDelegate: props.sizeUpdateHandler,
        };
    }

    reportSizeChange() {
        this.state.sizeChangeDelegate(this.state.taskDetails.task_id, this.getHeight());
    }

    componentDidMount() {
        this.mounted = true;
        console.log("task-positioner componentDidMount");
        this.reportSizeChange();
    }

    componentDidUpdate() {
        console.log("task-positioner componentDidUpdate");
        // this.reportSizeChange();
    }

    // Returns the height of element by pixels and by
    // minding the responsive design constraints
    getHeight() {
        if (this.mounted) {
            // var s = thisq.div.current.offsetHeight;
            // debugger;
            return this.div.current.offsetHeight; // Real size
        } else {
            return 32; // Estimated size
        }
    }

    render() {
        console.log("task-positioner render");
        // var objectKey = "taskPositioner" + this.state.taskDetails.id;
        var style = {
            top: this.state.documentPositions.top,
            transform: "translateX(" + this.state.documentPositions.left + "px)",
            // display: "none",
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

        // console.log(props);
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
