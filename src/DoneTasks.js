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

    updateTaskGroups() {}

    initialTaskCategorizationByDays() {
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
        this.setState((state, props) => ({
            childrenDivs: preparedChildrenDivs,
        }));
    }

    attachDelegates() {
        this.state.childrenDivs.forEach((childDiv) => {
            childDiv.setDelegate(this);
        });
    }

    componentDidMount() {
        this.initialTaskCategorizationByDays();
        this.attachDelegates();
    }

    componentDidUpdate() {
        console.log("vcdonetasks componentDidUpdate");
    }

    render() {
        console.log("vcdonetasks render");
        return <div>{this.state.childrenDivs}</div>;
    }
}

export default VCDoneTasks;
