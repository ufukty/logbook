// TODO: Update status + CSS ==> Active tasks

import React from "react";

import "./Task.css";

class Task extends React.Component {
    constructor(props) {
        super(props);
        console.log(props.data);

        var status;
        if (props.data.completed_at != null) {
            status = "done";
        } else if (props.data.ready_to_pick_up) {
            status = "ready-to-pick-up";
        } else {
            status = "pending";
        }

        this.state = {
            data: props.data,
            status: status,
        };
    }

    render() {
        return (
            <div className="task-wrapper">
                <div
                    className={"task " + this.state.status}
                    style={{
                        "--depth": this.state.data.depth,
                    }}
                    draggable="true"
                >
                    {this.state.data.content}
                </div>
            </div>
        );
    }
}

export default Task;
