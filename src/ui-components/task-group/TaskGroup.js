import React from "react";

import "./TaskGroup.css";

import TaskList from "./task-list/TaskList";
import * as DateTime from "./../../utility/dateTime";

class TaskGroup extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            group_header: DateTime.timestampToLocalizedText(props.timestampOfDay),
            group_items: props.tasks,
            group_type: props.groupType,
            documentViewMode: props.documentViewMode,
        };
    }
    static getDerivedStateFromProps(props, state) {
        return {
            group_header: DateTime.timestampToLocalizedText(props.timestampOfDay),
            group_items: props.tasks,
            group_type: props.groupType,
            documentViewMode: props.documentViewMode,
        };
    }

    render() {
        return (
            <div className={"task-group " + this.state.group_type}>
                <div className="task-group-header">{this.state.group_header}</div>
                <TaskList tasks={this.state.group_items} documentViewMode={this.state.documentViewMode} />
            </div>
        );
    }
}

export default TaskGroup;
