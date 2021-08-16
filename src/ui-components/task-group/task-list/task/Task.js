import React from "react";

import "./Task.css";

class Task extends React.Component {
    constructor(props) {
        super();
        this.state = props.task;
    }

    render() {
        return (
            <div className="task-wrapper">
                <div
                    className={"task " + (this.state.active ? "active" : "")}
                    draggable="true"
                >
                    {this.state.text}
                </div>
            </div>
        );
    }
}

export default Task;
