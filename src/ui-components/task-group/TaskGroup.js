import React from "react";

import "./TaskGroup.css";

import TaskList from "./task-list/TaskList";

class TaskGroup extends React.Component {
    constructor(props) {
        super();
        this.state = {
            group_header: props.group_header,
            group_items: props.group_items,
            group_type: props.group_type,
        };
    }

    render() {
        return (
            <div className={"task-group " + this.state.group_type}>
                <div className="task-group-header">
                    {this.state.group_header}
                </div>
                <TaskList tasks={this.state.group_items} />
            </div>
        );
    }
}

export default TaskGroup;
