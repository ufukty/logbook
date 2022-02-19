// TODO: Update status + CSS ==> Active tasks

import React from "react";

import "./Task.css";

class Task extends React.Component {
    constructor(props) {
        super(props);

        // console.log(props);
        var status;
        if (props.task.completed_at != null) {
            status = "done";
        } else if (props.task.ready_to_pick_up) {
            status = "ready-to-pick-up";
        } else {
            status = "pending";
        }

        this.state = {
            task: props.task,
            status: status,
            documentViewMode: props.documentViewMode,
        };
    }

    render() {
        return (
            <div
                className="task-wrapper"
                style={{
                    "--depth": this.state.task.depth,
                }}
                documentViewMode={this.state.documentViewMode}
            >
                <div
                    task_id={this.state.task.task_id}
                    task_depth={this.state.task.depth}
                    className={"task " + this.state.status}
                    // draggable="true"
                    contentEditable="true"
                    suppressContentEditableWarning={true}
                >
                    {this.state.task.content}
                </div>
            </div>
        );
    }
}

export default Task;
