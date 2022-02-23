import "./App.css";

// import ChronologicalView from "./views/chronological-view/Chronological";
// import TreeView from "./views/tree-view/TreeView";

// import Task from "./ui-components/task-group/task-list/task/Task";

// import VCDoneTasks from "./DoneTasks";

import React from "react";
import TaskPositioner from "./ui-components/task-group/task-list/task/Task";

var endpoint_address = "http://192.168.1.44:8080";
// var endpoint_document_overview_hierarchical = "/document/overview/hierarchical";
var endpoint_document_overview_chronological = "/document/overview/chronological/";

// function leastOfGreaterBinarySearch(list, item) {
//     var lo = 0,
//         hi = list.length - 1,
//         mid = -1;

//     while (hi - lo > 1) {
//         mid = Math.floor((lo + hi) / 2);
//         if (list[mid] <= item) {
//             lo = mid + 1;
//         } else {
//             hi = mid - 1;
//         }
//     }
//     return hi;
// }

// function searchIndex(items, item) {
//     for (const i in items) {
//         if (items[i] === item) {
//             return i;
//         }
//     }
//     return -1;
// }

class InfiniteSheet extends React.Component {
    constructor(props) {
        super(props);

        console.log("infinite-sheet constructor");

        // That property of component will hold the cell heights. Each component
        // will call a method of this class when its height changed and update
        // its record in this array. Initial values are just mock values because
        // real values needs each component to be mounted.
        this.cellHeights = this.getMockCellHeights(props.tasks);

        // Each item is an task_id of a task. Ordering of items represents the
        // rendering order of task components. Toggling between Chronological
        // and Hierarchical document view modes, directly changes the ordering
        // of items of this array and eventually will be represented in page
        // by render method.
        this.domOrdering = props.tasks.map((task) => {
            return task.task_id;
        });

        this.initializedCellsForTasks = {};
        this.initializedCellsForTasksRefs = {};

        this.state = {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
            childrenInitialized: false,

            // Based on mock values for initialization time
            verticalCellPositions: Array.from(Array(props.tasks.length).keys()).map((x) => {
                return x * 64 + 200;
            }),
        };

        // Task cells will be initialized with mock positions for calculating
        // real heights and positions.
        this.initializeTaskCells();
    }

    static getDerivedStateFromProps(props, state) {
        console.log("infinite-sheet getDerivedStateFromProps");
        return {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }

    initializeTaskCells() {
        var updateHeightRecordOfComponent = this.updateHeightRecordOfComponent.bind(this);
        for (var rowIndex = 0; rowIndex < this.state.tasksRawData.length; rowIndex++) {
            var task = this.state.tasksRawData[rowIndex];
            this.initializedCellsForTasksRefs[task.task_id] = React.createRef();
            this.initializedCellsForTasks[task.task_id] = (
                <TaskPositioner
                    ref={this.initializedCellsForTasksRefs[task.task_id]}
                    key={task.task_id}
                    taskDetails={task}
                    documentViewMode={this.state.documentViewMode}
                    posY={this.state.verticalCellPositions[rowIndex]}
                    sizeUpdateHandler={updateHeightRecordOfComponent}
                />
            );
        }
    }

    getMockCellHeights(tasks) {
        var heights = {};
        for (const task of tasks) {
            heights[task.task_id] = 32;
        }
        return heights;
    }
    getNumberOfRows() {
        return this.state.tasksRawData.length;
    }

    getVerticalCellPositions(childrenMounted) {
        if (this.state.childrenInitialized) {
            var totalNumberOfCells = this.initializedCellsForTasks.length;
        } else {
            var totalNumberOfCells = this.state.tasksRawData.length;
        }
        var padding = 32;
        if (childrenMounted) {
            console.log("initialized");
            var cumulativeHeights = [200];
            for (var cellIndex = 1; cellIndex < totalNumberOfCells + 1; cellIndex++) {
                var prevCellPosition = cumulativeHeights[cumulativeHeights.length - 1];
                var currentCellHeight = this.getCellHeight(cellIndex - 1);
                cumulativeHeights.push(prevCellPosition + currentCellHeight + padding);
            }
            return cumulativeHeights.slice(0);
        } else {
            console.log("not initialized");
            return Array.from(Array(totalNumberOfCells).keys()).map((x) => {
                return x * 64 + 200;
            });
        }
    }

    getCellHeight(cellIndex) {
        var cellTaskId = this.domOrdering[cellIndex];
        return this.cellHeights[cellTaskId];
    }

    componentDidMount() {
        console.log("infinite-sheet componentDidMount");
        var positions = this.getVerticalCellPositions(true);
        // this.setState({
        //     childrenInitialized: true,
        //     verticalCellPositions: positions,
        // });

        this.domOrdering.forEach((taskId, rowIndex) => {
            this.initializedCellsForTasksRefs[taskId].current.div.current.style.top = positions[rowIndex] + "px";
        });
    }

    updateHeightRecordOfComponent(taskPositionerTaskId, newHeight) {
        console.log(taskPositionerTaskId, newHeight);
        this.cellHeights[taskPositionerTaskId] = newHeight;
        // debugger;
        this.setState({
            verticalCellPositions: this.getVerticalCellPositions(true),
        });
        // debugger;
        // this.initializedCellsForTasksRefs[taskPositionerTaskId].current.forceUpdate();
    }

    render() {
        console.log("infinite-sheet render");
        var content = [];
        for (const taskId of this.domOrdering) {
            content.push(this.initializedCellsForTasks[taskId]);
        }

        console.log("positions: ", this.state.verticalCellPositions);

        return (
            <div
                id="infinite-sheet"
                style={{
                    "--focus-depth": "1",
                    background: "url('img/dot-background.png')",
                }}
            >
                {content}
            </div>
        );
    }
}

class ModeSelector extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedMode: 0,
            documentViewModeChangeDelegate: props.documentViewModeChangeDelegate,
        };
        this.eventSwitchModes = this.eventSwitchModes.bind(this);
    }

    eventSwitchModes(e) {
        e.preventDefault();
        var changeModeTo = 1 - this.state.selectedMode;
        this.setState((state, props) => ({
            selectedMode: changeModeTo,
        }));
        this.state.documentViewModeChangeDelegate(["chro", "hier"][changeModeTo]);
    }

    render() {
        var classNameForWrapper = ["left-picked", "right-picked"][this.state.selectedMode];
        return (
            <div id="settings-documentViewMode" className={classNameForWrapper} onClick={this.eventSwitchModes}>
                <div id="left">C</div>
                <div id="right">H</div>
                <div id="left-activated-caption">Chronological</div>
                <div id="right-activated-caption">Hierarchical</div>
            </div>
        );
    }
}

class App extends React.Component {
    constructor() {
        super();
        this.documentViewModeChangeHandler = this.documentViewModeChangeHandler.bind(this);

        this.documentViewModeSelector = (
            <ModeSelector documentViewModeChangeDelegate={this.documentViewModeChangeHandler}></ModeSelector>
        );
        this.state = {
            overviewIsLoaded: false,
            response: null,
            error: null,
            documentViewMode: "chro", // hier | chro
        };
    }

    documentViewModeChangeHandler(newMode) {
        this.setState({
            documentViewMode: newMode,
        });
    }

    componentDidMount() {
        var documentId = "61bbc44a-c61c-4d49-8804-486181081fa7";
        fetch(endpoint_address + endpoint_document_overview_chronological + documentId + "?limit=1000&offset=0")
            .then((result) => result.json())
            .then(
                (result) => {
                    this.setState((state, props) => ({
                        overviewIsLoaded: true,
                        response: result,
                    }));
                },
                (error) => {
                    this.setState((state, props) => ({
                        overviewIsLoaded: true,
                        error: error,
                    }));
                }
            );
    }

    render() {
        console.log("app render");
        var content;

        if (this.state.error) {
            content = <div>{this.state.error.message}</div>;
        } else if (!this.state.overviewIsLoaded) {
            content = <div>Loading...</div>;
        } else {
            content = (
                <InfiniteSheet
                    tasks={this.state.response.resource}
                    documentViewMode={this.state.documentViewMode}
                ></InfiniteSheet>
            );
        }

        return (
            <div className="document-sheet">
                <a id="home-button" className="floating-corner left top" href="index.html">
                    Logbook
                </a>

                <div id="sheet-settings" className="floating-corner right top dark">
                    <div>Share</div>

                    <div>Sync</div>
                </div>

                {/* <div id="active-task-details" className="floating-corner left bottom light">
                    History for active task
                    <div className="task">PAM for SSH</div>
                    <div className="task">ACL - Redis</div>
                    <div className="task">TOTP for SSH</div>
                </div> */}

                <div id="date-anchors" className="floating-corner right bottom light">
                    <a href="#august-10-2021">10th August</a>
                    <a href="#august-12-2021">12th August</a>
                    <a href="#august-13-2021">13th August</a>
                    <a href="#august-14-2021">Active Tasks</a>
                    <a href="#august-14-2021">To-do Drawer</a>
                </div>

                {/* <div id="debug" className="floating-corner left bottom light">
                    Welcome back.
                </div> */}

                <div id="settings" className="floating-corner left bottom">
                    {this.documentViewModeSelector}
                </div>

                {content}
            </div>
        );
    }
}

export default App;
