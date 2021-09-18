import "./App.css";

// import ChronologicalView from "./views/chronological-view/Chronological";
// import TreeView from "./views/tree-view/TreeView";

import Task from "./ui-components/task-group/task-list/task/Task";
import React from "react";

// var task_events = [
//     {
//         log_id: 1,
//         event_type: "new child",
//         event_time: "198821988219882",
//     },
// ];

var endpoint_address = "http://localhost:8080";
var endpoint_document_overview = "/document/overview/";

class App extends React.Component {
    constructor() {
        super();
        this.state = {
            overviewIsLoaded: false,
            tasks: null,
            error: null,
        };
    }

    componentDidMount() {
        var documentId = "24304d83-633a-4f6a-81a6-d32bdc604d9f";
        fetch(endpoint_address + endpoint_document_overview + documentId)
            .then((result) => result.json())
            .then(
                (result) => {
                    console.log(result.resource);
                    var tasks = result.resource.map((task) => (
                        <Task key={task.task_id} data={task} />
                    ));

                    // console.log(tasks);

                    this.setState({
                        overviewIsLoaded: true,
                        response: result,
                        tasks: tasks,
                    });
                },
                (error) => {
                    this.setState({
                        overviewIsLoaded: true,
                        error: error,
                    });
                }
            );
    }

    render() {
        var content;

        // if (this.state.mode === "chronological") {
        //     content = <ChronologicalView dataset={dataset} />;
        // } else if (this.state.mode === "tree") {
        //     content = <TreeView />;
        // }

        content = (function (state) {
            const { overviewIsLoaded, tasks, error } = state;
            if (error) {
                return <div>{error.message}</div>;
            } else if (!overviewIsLoaded) {
                return <div>Loading...</div>;
            } else {
                return <div>{tasks}</div>;
            }
        })(this.state);

        return (
            <div className="document-sheet">
                <a
                    id="home-button"
                    className="floating-corner left top"
                    href="index.html"
                >
                    Logbook
                </a>

                <div
                    id="sheet-settings"
                    className="floating-corner right top dark"
                >
                    <div>Share</div>

                    <div>Sync</div>
                </div>

                <div
                    id="active-task-details"
                    className="floating-corner left bottom"
                >
                    History for active task
                    <div className="task">PAM for SSH</div>
                    <div className="task">ACL - Redis</div>
                    <div className="task">TOTP for SSH</div>
                </div>

                <div id="date-anchors" className="floating-corner right bottom">
                    <a href="#august-10-2021">10th August</a>
                    <a href="#august-12-2021">12th August</a>
                    <a href="#august-13-2021">13th August</a>
                    <a href="#august-14-2021">Active Tasks</a>
                    <a href="#august-14-2021">To-do Drawer</a>
                </div>

                {content}
            </div>
        );
    }
}

export default App;
