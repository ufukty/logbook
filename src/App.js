import "./App.css";

// import ChronologicalView from "./views/chronological-view/Chronological";
// import TreeView from "./views/tree-view/TreeView";

// import Task from "./ui-components/task-group/task-list/task/Task";

import VCDoneTasks from "./DoneTasks";

import React from "react";

var endpoint_address = "http://192.168.1.44:8080";
// var endpoint_document_overview_hierarchical = "/document/overview/hierarchical";
var endpoint_document_overview_chronological = "/document/overview/chronological/";

// class DataStore {
//     constructor() {}
// }

class InfiniteSheet extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }

    static getDerivedStateFromProps(props, state) {
        return {
            tasksRawData: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }

    

    prepareChildren() {
        var preparedChildrenDivs = [];

        var doneTasks = this.state.tasksRawData.filter((item) => item.completed_at != null);
        if (doneTasks.length > 0) {
            preparedChildrenDivs.push(
                <VCDoneTasks
                    key="done-tasks"
                    tasks={doneTasks}
                    documentViewMode={this.state.documentViewMode}
                ></VCDoneTasks>
            );
        }

        // TODO: other types of tasks such as; ready-to-pick-up, to-do, active, paused etc..

        return preparedChildrenDivs;
    }

    render() {
        var childrenDivs = this.prepareChildren();
        return (
            <div
                id="infinite-sheet"
                style={{
                    "--focus-depth": "1",
                    background: "url('img/dot-background.png')",
                }}
            >
                {childrenDivs}
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
