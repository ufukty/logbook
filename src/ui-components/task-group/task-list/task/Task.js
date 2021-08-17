import React from "react";

import "./Task.css";

class Task extends React.Component {
    constructor(props) {
        super();
        this.state = {
            info: props.info,
        };
    }

    render() {
        return (
            <div className="task-wrapper">
                <div
                    className={"task " + this.state.info.task.task_status}
                    style={{
                        "--depth": this.state.info.depth,
                    }}
                    draggable="true"
                >
                    {this.state.info.task.text}
                </div>
            </div>
        );
    }
}

export default Task;
