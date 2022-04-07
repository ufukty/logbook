import "./css/app.css";
import "./css/document-view-mode-selector.css";
import "./css/infinite-sheet.css";

import React from "react";
import * as constants from "./constants";
import TaskPositioner from "./ui-components/task-group/task-list/task/Task";
import { toHaveAccessibleDescription } from "@testing-library/jest-dom/dist/matchers";
// import { autoFocusManager } from "./AutoFocusManager";

var endpoint_address = "http://192.168.1.44:8080";
// var endpoint_document_overview_hierarchical = "/document/overview/hierarchical";
var endpoint_document_overview_chronological = "/document/overview/chronological/";

function findFirstGreaterOrClosestItem(orderedList, searchItem) {
    var lastItemIndex = orderedList.length - 1;
    if (searchItem <= orderedList[0]) {
        // if searchItem is smaller than the smallest item on orderedList
        return 0;
    } else if (orderedList[lastItemIndex] <= searchItem) {
        // if searchItem is bigger than the biggest item on orderedList
        return lastItemIndex;
    } else {
        // if searchItem is in between first and last item of
        // orderedList, perform below instructions based on
        // binary search.
        var lo = 0,
            hi = lastItemIndex,
            mid = undefined;
        while (hi - lo > 1) {
            mid = Math.floor((lo + hi) / 2);
            if (orderedList[mid] <= searchItem) {
                lo = mid;
            } else {
                hi = mid;
            }
        }
        return lo;
    }
}

function Task(props) {
    // Properties that are fetched from server
    this.taskId = props.task_id; // example: "8557d156-3d00-4836-8323-a9bdd586719a"
    this.documentId = props.document_id; // example: "61bbc44a-c61c-4d49-8804-486181081fa7"
    this.parentId = props.parent_id; // example: "999c060e-d853-4271-b529-42c2655a4aae"
    this.content = props.content; // example: "Update redis/tf file according to prod.tfvars file"
    this.degree = props.degree; // example: 1
    this.depth = props.depth; // example: 4
    this.createdAt = props.created_at; // example: "2022-01-27T01:39:51.320386Z"
    this.completedAt = props.completed_at; // example: "2022-02-17"
    this.readyToPickUp = props.ready_to_pick_up; // example: true

    // Properties that are created and used by frontend
    this.effectiveDepth = 0;
}

class InfiniteSheet extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
            chronologicalOrdering: props.chronologicalOrdering,
            paneStyleLeftInPixels: 0,
        };
        this.autoFocusSetup();
        this.debug();
    }

    static getDerivedStateFromProps(props, state) {
        return {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
            chronologicalOrdering: props.chronologicalOrdering,
        };
    }

    autoFocusSetup() {
        window.addEventListener("scroll", this.scrollEventHandler.bind(this));

        this.currentlyFocusedTask_DOMObject = undefined;
        this.effectiveDVM = constants.DVM_CHRONO;

        // this.hierarchicalModeIsEnabled = false;
        this.focusDepthOnTransition = 0;
        this.paneShiftPrior = 0;
        this.paneShiftCurrent = 0;
    }

    componentDidMount() {
        this.focusHandler();
    }

    componentDidUpdate() {
        this.inspectDVMChange();
    }

    inspectDVMChange() {
        var nextDVM = this.state.documentViewMode;
        if (this.effectiveDVM === nextDVM) {
            return;
        } else {
            this.effectiveDVM = nextDVM;
        }

        if (nextDVM === constants.DVM_CHRONO) {
            this.switchDVMToChronological();
        } else if (nextDVM === constants.DVM_HIERARCH) {
            this.switchDVMToHierarchical();
        }
        this.debug();
    }

    switchDVMToChronological() {
        this.paneShiftPrior += this.paneShiftCurrent;
        this.paneShiftCurrent = 0;

        var initiallyFocusedDepth = this.getCurrentlyFocusedDepth();
        // var initiallyFocusedDepth =
        var updatedTasks;
        this.setState((prevState) => {
            updatedTasks = Object.assign({}, prevState.tasks);
            for (const taskId in updatedTasks) {
                if (Object.hasOwnProperty.call(updatedTasks, taskId)) {
                    var task = updatedTasks[taskId];
                    task.effectiveDepth = initiallyFocusedDepth;
                }
            }
            return {
                tasks: updatedTasks,
            };
        });
    }

    // save the focus depth
    switchDVMToHierarchical() {
        if (this.currentlyFocusedTask_DOMObject === undefined) {
            return;
        }
        this.focusDepthOnTransition = this.getCurrentlyFocusedDepth();

        var initiallyFocusedDepth = this.getCurrentlyFocusedDepth();
        var absoluteDepthOfFocusedTask =
            this.state.tasks[this.currentlyFocusedTask_DOMObject.getAttribute("task_id")].depth;
        var updatedTasks;
        this.setState((prevState) => {
            updatedTasks = Object.assign({}, prevState.tasks);
            for (const taskId in updatedTasks) {
                if (Object.hasOwnProperty.call(updatedTasks, taskId)) {
                    var task = updatedTasks[taskId];
                    var deltaDepth = task.depth - absoluteDepthOfFocusedTask;
                    task.effectiveDepth = deltaDepth + initiallyFocusedDepth;
                }
            }
            return {
                tasks: updatedTasks,
            };
        });
    }

    positionOfFieldOfFocus() {
        var upperEdgeOfScreen = window.scrollY,
            heightOfScreen = window.innerHeight;
        return upperEdgeOfScreen + 0.3 * heightOfScreen;
    }

    debug() {
        var debugActivated = false;
        if (debugActivated) {
            this.debugElem = document.getElementById("debug");
            // if (this.currentlyFocusedTask_DOMObject !== undefined) {
            this.debugElem.innerHTML =
                "focusDepthOnTransition: " +
                this.focusDepthOnTransition +
                "<br>" +
                "positionOfFieldOfFocus: " +
                this.positionOfFieldOfFocus().toFixed(2) +
                "<br>---<br>" +
                "paneShiftPrior: " +
                this.paneShiftPrior +
                "<br>" +
                "paneShiftCurrent: " +
                this.paneShiftCurrent +
                "<br>" +
                "paneTranslationInPixels: " +
                this.state.paneTranslationInPixels;

            if (this.currentlyFocusedTask_DOMObject !== undefined) {
                this.debugElem.innerHTML +=
                    "<br>---<br>" +
                    "currentlyFocusedTask.depth: " +
                    this.state.tasks[this.currentlyFocusedTask_DOMObject.getAttribute("task_id")].depth +
                    "<br>" +
                    "currentlyFocusedTask.effectiveDepth: " +
                    this.state.tasks[this.currentlyFocusedTask_DOMObject.getAttribute("task_id")].effectiveDepth +
                    "<br>" +
                    "currentlyFocusedTask.innerText: " +
                    "<br>" +
                    this.currentlyFocusedTask_DOMObject.innerText;
            }
        }
    }

    // focusElements should be a NodeList
    getBoundariesOfFocusAreas(focusElements) {
        return [...focusElements].map((task) => task.offsetTop);
    }

    findTasksFromDOM() {
        return document.querySelectorAll(".task");
    }

    scrollEventHandler(e) {
        e.preventDefault();
        this.focusHandler();
        this.debug();
    }

    focusHandler() {
        var tasksFromDOM = this.findTasksFromDOM();

        // Detect if there is a task in focus
        var nextFocusIndex = this.detectFocusedTaskIndex(tasksFromDOM);
        if (nextFocusIndex === -1) {
            if (this.currentlyFocusedTask_DOMObject !== undefined) {
                // if there is already a task in focus, then this is not the
                // first call of this function. so, we will skip the rest of
                // process and keep the task in focus
                return;
            } else {
                // if there is no task in focus, then this should be the first
                // call of this function. we will continue with focusing the
                // task at top
                nextFocusIndex = 0;
            }
        }

        var nextFocusedTask_DOMObject = tasksFromDOM[nextFocusIndex];

        // If this is the first time we focus on a task (first call of this function)
        if (this.currentlyFocusedTask_DOMObject === undefined) {
            // Register currently focused task
            this.currentlyFocusedTask_DOMObject = nextFocusedTask_DOMObject;
            this.addHighlightToCurrentlyFocusedTask();
        } else {
            this.focusDepth = this.getCurrentlyFocusedDepth();

            // Detect if the task in focus is the same with previous call
            if (this.currentlyFocusedTask_DOMObject === nextFocusedTask_DOMObject) {
                return;
            }
            this.removeHighlightFromFromCurrentlyFocusedTask();

            // Register currently focused task
            this.currentlyFocusedTask_DOMObject = nextFocusedTask_DOMObject;
            this.addHighlightToCurrentlyFocusedTask();
        }

        // calculatePaneShiftCurrent
        this.currentlyFocusedTask_DOMObject = nextFocusedTask_DOMObject;
        var effectiveDepthOfTaskInFocus =
            this.state.tasks[nextFocusedTask_DOMObject.getAttribute("task_id")].effectiveDepth;
        this.paneShiftCurrent = effectiveDepthOfTaskInFocus - this.focusDepthOnTransition;

        // applyFocusedDepth
        var paneTranslationInPixels;
        if (this.effectiveDVM === constants.DVM_HIERARCH) {
            paneTranslationInPixels =
                -1 * (this.paneShiftPrior + this.paneShiftCurrent) * constants.AUTO_FOCUS_SHIFT_IN_PIXELS;
        } else {
            paneTranslationInPixels = -1 * this.paneShiftPrior * constants.AUTO_FOCUS_SHIFT_IN_PIXELS;
        }
        this.setState({ paneTranslationInPixels: paneTranslationInPixels });
        this.debug();
    }

    detectFocusedTaskIndex(elements) {
        var boundaries = this.getBoundariesOfFocusAreas(elements);
        var focusIndex = findFirstGreaterOrClosestItem(boundaries, this.positionOfFieldOfFocus());
        return focusIndex;
    }

    addHighlightToCurrentlyFocusedTask() {
        this.currentlyFocusedTask_DOMObject.classList.add("focused-task");
    }

    removeHighlightFromFromCurrentlyFocusedTask() {
        this.currentlyFocusedTask_DOMObject.classList.remove("focused-task");
    }

    getCurrentlyFocusedDepth() {
        if (this.currentlyFocusedTask_DOMObject === undefined) {
            return 0;
        } else {
            return this.state.tasks[this.currentlyFocusedTask_DOMObject.getAttribute("task_id")].effectiveDepth;
        }
    }

    render() {
        var content = this.state.chronologicalOrdering.map((taskId) => {
            return <TaskPositioner key={taskId} task={this.state.tasks[taskId]} />;
        });

        // var className = "auto-focus-shift-area auto-focus-enabled";
        // if (this.state.documentViewMode === constants.DVM_HIERARCH) {
        //     className += " auto-focus-enabled";
        // }

        return (
            <div
                id="infinite-sheet"
                // className={className}
                style={{
                    transform: "translateX(calc(" + this.state.paneTranslationInPixels + "px - 50%))",
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
        this.setState({
            selectedMode: changeModeTo,
        });
        var delegate = this.state.documentViewModeChangeDelegate;
        delegate([constants.DVM_CHRONO, constants.DVM_HIERARCH][changeModeTo]);
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
        this.fetchDocumentFromServer(documentId);
    }

    fetchDocumentFromServer(documentId) {
        fetch(endpoint_address + endpoint_document_overview_chronological + documentId + "?limit=1000&offset=0")
            .then((result) => result.json())
            .then(
                (result) => {
                    // Create < associative array || key-value list >

                    var tasks = new Object();
                    result.resource.forEach((resource) => {
                        tasks[resource.task_id] = new Task(resource);
                    });

                    var chronologicalOrdering = result.resource.map((resource) => resource.task_id);

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
            content = <div>{this.state.error.message}</div>;
        } else if (!this.state.overviewIsLoaded) {
            content = <div>Loading...</div>;
        } else {
            content = (
                <InfiniteSheet
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
            </div>
        );
    }
}

export default App;
