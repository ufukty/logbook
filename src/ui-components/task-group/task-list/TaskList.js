import React from "react";

import "./TaskList.css";

import Task from "./task/Task";

class TaskList extends React.Component {
    constructor(props) {
        super();
        this.state = {
            children: props.tasks.map((task) => (
                <Task key={task.id} task={task} />
            )),
        };
    }

    render() {
        return <div className="task-list">{this.state.children}</div>;
    }
}

export default TaskList;
