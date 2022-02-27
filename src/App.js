import "./css/app.css";
import "./css/document-view-mode-selector.css";
import "./css/infinite-sheet.css";

import React from "react";
import * as constants from "./constants";
import TaskPositioner from "./ui-components/task-group/task-list/task/Task";
import { autoFocusManager } from "./AutoFocusManager";
// ("./AutoFocusManager");

var endpoint_address = "http://192.168.1.44:8080";
// var endpoint_document_overview_hierarchical = "/document/overview/hierarchical";
var endpoint_document_overview_chronological = "/document/overview/chronological/";

class InfiniteSheet extends React.Component {
    constructor(props) {
        super(props);
        // console.log("infinite-sheet constructor");
        this.state = {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
            adjustmentForDepth: 1,
            domOrdering: props.tasks.map((taskDetails) => taskDetails.task_id),
        };
        autoFocusManager.registerDelegateForAdjustmentDepthChange(
            this.delegateAutoFocusAdjustmentDepthChange.bind(this)
        );
    }

    static getDerivedStateFromProps(props, state) {
        // console.log("infinite-sheet getDerivedStateFromProps");
        return {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }
    delegateAutoFocusAdjustmentDepthChange(newAdjustmentDepth) {
        this.setState({
            adjustmentForDepth: newAdjustmentDepth,
        });
    }

    render() {
        // console.log("infinite-sheet render");

        var content = [];
        for (const taskId of this.state.domOrdering) {
            var task = this.state.tasksRawData.filter((taskDetails) => taskDetails.task_id === taskId)[0];
            var taskPositioner = (
                <TaskPositioner
                    key={task.task_id}
                    taskDetails={task}
                    adjustmentForDepth={this.state.adjustmentForDepth}
                    documentViewMode={this.state.documentViewMode}
                />
            );
            content.push(taskPositioner);
        }

        var className = "auto-focus-shift-area";
        if (this.state.documentViewMode === constants.DVM_HIERARCH) {
            className += " auto-focus-enabled";
        }

        autoFocusManager.registerNextDVM(this.state.documentViewMode);

        return (
            <div
                id="infinite-sheet"
                className={className}
                style={{
                    // "--focus-depth": "4",
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
        document.addEventListener("keyup", this.keyboardEventListener.bind(this));
    }

    keyboardEventListener(e) {
        if (e.ctrlKey && (e.key === "c" || e.key === "C")) {
            this.eventSwitchModes(e);
        }
    }

    eventSwitchModes(e) {
        e.preventDefault();
        var changeModeTo = 1 - this.state.selectedMode;
        this.setState((state, props) => ({
            selectedMode: changeModeTo,
        }));
        this.state.documentViewModeChangeDelegate([constants.DVM_CHRONO, constants.DVM_HIERARCH][changeModeTo]);
    }

    render() {
        var classNameForWrapper = ["left-picked", "right-picked"][this.state.selectedMode];
        var eventSwitchModes = this.eventSwitchModes.bind(this);
        return (
            <div id="settings-documentViewMode" className={classNameForWrapper} onClick={eventSwitchModes}>
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
        this.documentViewModeSelector = (
            <ModeSelector documentViewModeChangeDelegate={this.documentViewModeChangeHandler.bind(this)}></ModeSelector>
        );

        this.state = {
            overviewIsLoaded: false,
            response: null,
            error: null,
            documentViewMode: constants.DVM_CHRONO,
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
        // console.log("app render");
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
                {/* <a id="home-button" className="floating-corner left top" href="index.html">
                    Logbook
                </a> */}

                {/* <div id="sheet-settings" className="floating-corner right top dark">
                    <div>Share</div>

                    <div>Sync</div>
                </div> */}

                {/* <div id="active-task-details" className="floating-corner left bottom light">
                    History for active task
                    <div className="task">PAM for SSH</div>
                    <div className="task">ACL - Redis</div>
                    <div className="task">TOTP for SSH</div>
                </div> */}

                {/* <div id="date-anchors" className="floating-corner right bottom light">
                    <a href="#august-10-2021">10th August</a>
                    <a href="#august-12-2021">12th August</a>
                    <a href="#august-13-2021">13th August</a>
                    <a href="#august-14-2021">Active Tasks</a>
                    <a href="#august-14-2021">To-do Drawer</a>
                </div> */}

                <div id="debug" className="floating-corner right bottom light">
                    Welcome back.
                </div>

                <div id="settings" className="floating-corner left bottom">
                    {this.documentViewModeSelector}
                </div>

                {content}
            </div>
        );
    }
}

export default App;
