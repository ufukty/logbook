import React from "react";

import * as constants from "../utility/constants";
import * as misc from "../utility/misc";

import TaskPositioner from "./Task";
import ContextMenu from "./ContextMenu";
import DayHeader from "./DayHeader";
import { toHaveAccessibleDescription } from "@testing-library/jest-dom/dist/matchers";

class InfiniteSheet extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
            chronologicalOrdering: props.chronologicalOrdering,
            paneShiftTotal: 0,
            effectiveDepthForDayHeaders: 0,
            contextMenuEnabled: false,
            contextMenuPosY: 0,
            contextMenuPosX: 0,
        };
        this.autoFocusSetup();
        this.contextMenuSetup();
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
        this._isMounted = true;
        this.currentlyFocusedTask_DOMObject = this.findTasksFromDOM[0];
        this.focusHandler();
    }

    componentWillUnmount() {
        this._isMounted = false;
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

        var initiallyFocusedDepth = this.state.paneShiftTotal * -1; // this.getCurrentlyFocusedDepth();
        this.setState((prevState) => {
            var updatedTasks = Object.assign({}, prevState.tasks);
            for (const taskId in updatedTasks) {
                if (Object.hasOwnProperty.call(updatedTasks, taskId)) {
                    var task = updatedTasks[taskId];
                    task.effectiveDepth = initiallyFocusedDepth;
                }
            }
            return {
                tasks: updatedTasks,
                effectiveDepthForDayHeaders: initiallyFocusedDepth,
            };
        });
        this.debug();
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
        this.debug();
    }

    positionOfFieldOfFocus() {
        var upperEdgeOfScreen = window.scrollY,
            heightOfScreen = window.innerHeight;
        return upperEdgeOfScreen + 0.5 * heightOfScreen;
    }

    debug() {
        this.debugElem = document.getElementById("debug");
        if (this.debugElem !== null) {
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
                "paneShiftTotal: " +
                this.state.paneShiftTotal;

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

    getAverageEffectiveFocusDepthOfFocusedArea(tasksFromDOM, nextFocusedTask_DOMObject) {
        // var taskIdOfFocusedTask = nextFocusedTask_DOMObject.getAttribute("task_id");
        var orderOfFocusedTask = -1;
        for (let index = 0; index < tasksFromDOM.length; index++) {
            const taskDOMObject = tasksFromDOM[index];
            if (taskDOMObject.getAttribute("task_id") === nextFocusedTask_DOMObject.getAttribute("task_id")) {
                orderOfFocusedTask = index;
                break;
            }
        }
        var effectiveDepths = [];
        for (let offset = -2; offset <= 2; offset++) {
            var neighbourTaskOrder = orderOfFocusedTask + offset;
            if (0 >= neighbourTaskOrder || neighbourTaskOrder >= tasksFromDOM.length) {
                continue;
            }
            var neighbourTaskId = tasksFromDOM[neighbourTaskOrder].getAttribute("task_id");
            var neighbourTaskEffectiveDepth = this.state.tasks[neighbourTaskId].effectiveDepth;
            effectiveDepths.push(neighbourTaskEffectiveDepth);
        }
        return misc.averageInt(effectiveDepths);
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
                nextFocusIndex = 1;
            }
        }

        var nextFocusedTask_DOMObject = tasksFromDOM[nextFocusIndex];

        // calculate paneShiftCurrent
        // focused depth will be the average effective depth of 5 tasks
        // around the centered one
        this.currentlyFocusedTask_DOMObject = nextFocusedTask_DOMObject;
        var effectiveDepthOfTaskInFocus = this.getAverageEffectiveFocusDepthOfFocusedArea(
            tasksFromDOM,
            nextFocusedTask_DOMObject
        );
        this.paneShiftCurrent = effectiveDepthOfTaskInFocus - this.focusDepthOnTransition;

        // applyFocusedDepth
        var paneShiftTotal;
        if (this.effectiveDVM === constants.DVM_HIERARCH) {
            paneShiftTotal = -1 * (this.paneShiftPrior + this.paneShiftCurrent);
        } else {
            paneShiftTotal = -1 * this.paneShiftPrior;
        }
        this.setState({ paneShiftTotal: paneShiftTotal });
        this.debug();
    }

    detectFocusedTaskIndex(elements) {
        var boundaries = this.getBoundariesOfFocusAreas(elements);
        var focusIndex = misc.findFirstGreaterOrClosestItem(boundaries, this.positionOfFieldOfFocus());
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
            return this.getAverageEffectiveFocusDepthOfFocusedArea(
                this.findTasksFromDOM(),
                this.currentlyFocusedTask_DOMObject
            );
        }
    }

    contextMenuSetup() {

        const contextMenuEventHandler = function (e) {
            console.log(e)
            if (e.type !== "contextmenu" || !e.target.classList.contains("task")) {
                return
            }
            e.preventDefault();
            this.setState({ contextMenuEnabled: true, contextMenuPosY: e.pageY, contextMenuPosX: e.pageX })
            // console.log(e.target.getAttribute("task_id"))
        }

        const clickEventListener = function (e) {
            if (e.type !== "click") {
                return
            }
            this.setState({ contextMenuEnabled: false })
        }

        document.addEventListener("contextmenu", contextMenuEventHandler.bind(this));
        document.addEventListener("click", clickEventListener.bind(this));
    }

    content() {
        var content = [];
        console.log(this.state.effectiveDepthForDayHeaders);
        const dateStampsOrdered = Object.keys(this.state.chronologicalOrdering).sort();
        for (const dateStamp of dateStampsOrdered) {
            if (this.state.documentViewMode === constants.DVM_CHRONO) {
                content.push(
                    <DayHeader
                        key={dateStamp}
                        dateStamp={dateStamp}
                        effectiveDepth={this.state.effectiveDepthForDayHeaders}
                    ></DayHeader>
                );
            }
            const taskIDsOfDay = this.state.chronologicalOrdering[dateStamp];
            for (const taskID of taskIDsOfDay) {
                content.push(<TaskPositioner key={taskID} task={this.state.tasks[taskID]} />);
            }
        }
        return content;
    }

    render() {
        var content = this.content();
        console.log(this.state.contextMenuEnabled)
        return (
            <div>
                <div
                    id="infinite-sheet"
                    style={{
                        transform:
                            "translateX(calc(" +
                            this.state.paneShiftTotal +
                            " * var(--infinite-sheet-pixels-for-each-shift) - 50%))",
                    }}
                >
                    {content}
                </div>

                <ContextMenu
                    enabled={this.state.contextMenuEnabled}
                    posX={this.state.contextMenuPosX}
                    posY={this.state.contextMenuPosY}
                ></ContextMenu>
            </div>
        );
    }
}


export default InfiniteSheet;