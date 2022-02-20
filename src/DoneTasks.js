import TaskGroup from "./ui-components/task-group/TaskGroup";
import * as DateTime from "./utility/dateTime";

import React from "react";

class VCDoneTasks extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tasksRawData: props.tasks,
            childrenDivs: [],
            documentViewMode: props.documentViewMode,
        };
    }
    static getDerivedStateFromProps(props, state) {
        return {
            tasksRawData: props.tasks,
            childrenDivs: [],
            documentViewMode: props.documentViewMode,
        };
    }

    updateTaskGroups() {}

    createChildren() {
        var classifiedTasks = DateTime.classifyTasksByDays(this.state.tasksRawData);
        var sortedDays = Object.keys(classifiedTasks).sort();
        var preparedChildrenDivs = [];
        sortedDays.forEach((timestampOfDay) => {
            preparedChildrenDivs.push(
                <TaskGroup
                    key={"task-group" + timestampOfDay}
                    timestampOfDay={timestampOfDay}
                    tasks={classifiedTasks[timestampOfDay]}
                    groupType="done"
                    documentViewMode={this.state.documentViewMode}
                ></TaskGroup>
            );
        });
        return preparedChildrenDivs;
    }

    attachDelegates() {
        this.state.childrenDivs.forEach((childDiv) => {
            childDiv.setDelegate(this);
        });
    }

    componentDidMount() {
        this.createChildren();
        this.attachDelegates();
    }

    render() {
        var childrenDivs = this.createChildren();
        return <div>{childrenDivs}</div>;
    }
}

export default VCDoneTasks;
