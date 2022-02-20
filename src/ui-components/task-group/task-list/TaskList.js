import React from "react";

import "./TaskList.css";

import TaskPositioner from "./task/Task";

class TaskList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }
    static getDerivedStateFromProps(props, state) {
        return {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }

    render() {
        var total_height = 0;
        var children = this.state.tasks.map((task) => {
            var child = (
                <TaskPositioner
                    key={task.task_id}
                    taskDetails={task}
                    documentViewMode={this.state.documentViewMode}
                    posY={total_height}
                />
            );
            total_height += 60;
            return child;
        });
        return <div className="task-list">{children}</div>;
    }
}

export default TaskList;
