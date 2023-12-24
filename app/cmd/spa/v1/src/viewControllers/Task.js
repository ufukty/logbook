// TODO: Update status + CSS ==> Active tasks

import React from "react";

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
        var horizontalShift = this.state.task.effectiveDepth;
        var style = {
            transform: "translateX(calc(" + horizontalShift + " * var(--infinite-sheet-pixels-for-each-shift)))",
        };
        return (
            <div ref={this.div} className="task-positioner" style={style}>
                <TaskViewController key={this.state.task.taskId} task={this.state.task}></TaskViewController>
            </div>
        );
    }
}

class TaskViewController extends React.Component {
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
            // contentEditable="true"
            // suppressContentEditableWarning={true}
            >
                {this.state.task.content}
            </div>
        );
    }
}

export default TaskPositioner;
