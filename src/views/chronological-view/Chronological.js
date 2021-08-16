import React from "react";

import TaskGroup from "./../../ui-components/task-group/TaskGroup";

class ChronologicalView extends React.Component {
    constructor(props) {
        super();
        this.state = {
            dataset: props.dataset,
        };
    }

    render() {
        return (
            <div>
                {this.state.dataset.days.map((data) => (
                    <TaskGroup
                        key={data.day}
                        group_header={data.day}
                        group_items={data.tasks}
                        group_type="regular"
                    />
                ))}
                <TaskGroup
                    key="Active Tasks"
                    group_header="Active Tasks"
                    group_items={this.state.dataset.active_tasks}
                    group_type="active-tasks"
                />
                <TaskGroup
                    key="To-do Drawer"
                    group_header="To-do Drawer"
                    group_items={this.state.dataset.todo_drawer}
                    group_type="to-do-drawer"
                />
            </div>
        );
    }
}

export default ChronologicalView;
