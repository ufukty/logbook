// TODO: Update status + CSS ==> Active tasks

import React from "react";
// import { timestampToLocalizedText } from "../../../../utility/dateTime";

import "./Task.css";
import * as constants from "../../../../constants";

class TaskPositioner extends React.Component {
    constructor(props) {
        super(props);
        this.state = { task: props.task };
        this.div = React.createRef();
    }

    static getDerivedStateFromProps(props, state) {
        return { task: props.task };
    }

    render() {
        // console.log("task-positioner render");
        var horizontalShift = this.state.task.effectiveDepth * constants.AUTO_FOCUS_SHIFT_IN_PIXELS;
        var style = {
            transform: "translateX(" + horizontalShift + "px)",
        };
        return (
            <div ref={this.div} className="task-positioner" style={style}>
                <Task key={this.state.task.taskId} task={this.state.task}></Task>
            </div>
        );
    }
}

class Task extends React.Component {
    constructor(props) {
        super(props);
        var status;
        if (props.task.completedAt != null) {
            status = "done";
        } else {
            status = "pending";
        }
        this.state = {
            task: props.task,
            taskStatus: status,
        };
    }

    render() {
        // console.log("task render");
        return (
            <div
                task_id={this.state.task.taskId} // FIXME: task_id -> taskId
                task_depth={this.state.task.depth}
                className={"task " + this.state.taskStatus}
                // draggable="true"
                contentEditable="true"
                suppressContentEditableWarning={true}
            >
                {this.state.task.content}
            </div>
        );
    }
}

export default TaskPositioner;
