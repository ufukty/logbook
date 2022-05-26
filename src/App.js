import "./css/colors.css";
import "./css/app.css";
import "./css/document-view-mode-selector.css";
import "./css/infinite-sheet.css";
import "./css/tasks.css";

import React from "react";

import * as constants from "./utility/constants";
import { classifyTasksByDays } from "./utility/dateTime";

import InfiniteSheet from "./viewControllers/InfiniteSheet";
import ModeSelector from "./viewControllers/ModeSelector";
import ContextMenu from "./viewControllers/ContextMenu";

import { TaskModel } from "./models/Task";

var endpoint_address = "http://192.168.1.66:8080";
// var endpoint_document_overview_hierarchical = "/document/overview/hierarchical";
var endpoint_document_overview_chronological = "/document/overview/chronological/";

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
        this.fetchDocumentFromServer(documentId);
    }

    organizeTasksByCreationDay(tasks) {
        tasks.forEach((task) => {
            var s = task.createdAt;
        });
    }

    fetchDocumentFromServer(documentId) {
        fetch(endpoint_address + endpoint_document_overview_chronological + documentId + "?limit=1000&offset=0")
            .then((result) => result.json())
            .then(
                (result) => {
                    // Create < associative array || key-value list >

                    var tasks = {};
                    result.resource.forEach((resource) => {
                        tasks[resource.task_id] = new TaskModel(resource);
                    });

                    var chronologicalOrdering = classifyTasksByDays(tasks);

                    this.setState((state, props) => ({
                        overviewIsLoaded: true,
                        response: result,
                        tasks: tasks,
                        chronologicalOrdering: chronologicalOrdering,
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
            content = (
                <div
                    style={{
                        position: "absolute",
                        left: "50vw",
                        top: "50vh",
                        transform: "translate(-50%)",
                    }}
                >
                    {this.state.error.message}
                </div>
            );
        } else if (!this.state.overviewIsLoaded) {
            content = (
                <div
                    style={{
                        position: "absolute",
                        left: "50vw",
                        top: "50vh",
                        transform: "translate(-50%)",
                    }}
                >
                    Loading...
                </div>
            );
        } else {
            content = (
                <InfiniteSheet
                    // key="ae0bcf02-427f-5f9d-8cc9-2aa969c8e273"
                    tasks={this.state.tasks}
                    chronologicalOrdering={this.state.chronologicalOrdering}
                    documentViewMode={this.state.documentViewMode}
                ></InfiniteSheet>
            );
        }

        return (
            <div className="document-sheet">
                <a id="home-button" className="floating-corner left top" href="index.html">
                    Logbook
                </a>

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

                {/* <div id="debug" className="floating-corner right  bottom light">
                    Welcome back.
                </div> */}

                <div id="settings" className="floating-corner left bottom">
                    {this.documentViewModeSelector}
                </div>

                {content}

                <ContextMenu></ContextMenu>
            </div>
        );
    }
}

export default App;
