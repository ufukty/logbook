import {
    adoption,
    addEventListenerForNonTouchScreen,
    domElementReuseCollector,
    executeWhenDocumentIsReady,
} from "./utilities.js";

import ModeSelector from "./viewControllers/ModeSelector.js";
import InfiniteSheet from "./viewControllers/InfiniteSheet.js";
import ContextMenu from "./viewControllers/ContextMenu.js";
import InfiniteSheetTask from "./viewControllers/InfiniteSheetTask.js";
import InfiniteSheetHeader from "./viewControllers/InfiniteSheetHeader.js";

import { Task } from "./models/Task.js";
import { ChronologicalDocumentOverview } from "./models/DocumentOverviewChronological.js";

import { classifyTasksByDays } from "./dateTime.js";

/* Moving 2 pixels left is treated as same with 1 pixel up + 1 pixel left. */
function manhattanDistance(x1, y1, x2, y2) {
    return Math.abs(x1 - x2) + Math.abs(y1 - y2);
}

class UserInputResolver {
    constructor() {
        this.openContextMenuOnceCallback = () => {}; // should be assigned by user
        this.closeContextMenuOnceCallback = () => {}; // should be assigned by user

        this.touchMoveDistance = 0;
        this.isFingerMoved = false;
        this.touchMoveLastPoint = { x: 0, y: 0 };
    }

    /** @param {MouseEvent} e */
    clickEventReceiverNonTouchScreen(e) {
        this.closeContextMenuOnceCallback();
    }

    /** @param {MouseEvent} e */
    contextMenuEventReceiver(e) {
        e.preventDefault();
        // this.closeContextMenuOnceCallback();
        const targetElement = e.target;

        // console.log(e);
        if (targetElement.classList.contains("task")) {
            e.stopPropagation();

            const taskElement = e.target.parentNode;

            if (this.taskElementThatIsContextMenuFocusedOn === taskElement) {
                this.closeContextMenuOnceCallback();
            } else {
                const section = taskElement.dataset["section"];
                const row = taskElement.dataset["row"];
                const clickPosX = e.pageX;
                const clickPosY = e.pageY;

                this.closeContextMenuOnceCallback();
                this.openContextMenuOnceCallback(taskElement, section, row, clickPosX, clickPosY);
            }
        } else {
            this.closeContextMenuOnceCallback();
        }
    }

    /** @param {MouseEvent} e */
    clickEventReceiver(e) {
        // e.preventDefault();
        // this.closeContextMenuOnceCallback();
        const targetElement = e.target;

        // console.log(e);
        if (targetElement.classList.contains("task")) {
            e.stopPropagation();

            const taskElement = e.target.parentNode;

            if (this.taskElementThatIsContextMenuFocusedOn === taskElement) {
                this.closeContextMenuOnceCallback();
            } else {
                const section = taskElement.dataset["section"];
                const row = taskElement.dataset["row"];
                const clickPosX = e.pageX;
                const clickPosY = e.pageY;

                this.closeContextMenuOnceCallback();
                this.openContextMenuOnceCallback(taskElement, section, row, clickPosX, clickPosY);
            }
        } else {
            this.closeContextMenuOnceCallback();
        }
    }

    /** @param {TouchEvent} e */
    touchStartEventReceiver(e) {
        const taskElement = e.currentTarget;
        // const touchStartTime = Date.now();
        // console.log(e);
        this.touchMoveDistance = 0.0;
        this.isFingerMoved = false;
        this.touchMoveLastPoint.x = e.changedTouches[0].screenX;
        this.touchMoveLastPoint.y = e.changedTouches[0].screenY;
    }

    /** @param {TouchEvent} e */
    touchMoveEventReceiver(e) {
        // console.log(e);
        if (this.isFingerMoved) return;
        const touch = e.changedTouches[0];
        const last = this.touchMoveLastPoint;
        const lastMovementDistance = manhattanDistance(last.x, last.y, touch.screenX, touch.screenY);
        this.touchMoveDistance += lastMovementDistance;
        // console.log(this.touchMoveDistance);
        if (this.touchMoveDistance > 10) this.isFingerMoved = true;
        this.touchMoveLastPoint.x = touch.screenX;
        this.touchMoveLastPoint.y = touch.screenY;
    }

    /** @param {TouchEvent} e */
    touchEndEventReceiver(e) {
        const element = e.target;

        if (element.classList.contains("task")) {
            const taskElement = e.target.parentNode;

            if (this.isFingerMoved) return;

            if (this.taskElementThatIsContextMenuFocusedOn === taskElement) {
                this.closeContextMenuOnceCallback();
            } else {
                const section = taskElement.dataset["section"];
                const row = taskElement.dataset["row"];
                const clickPosX = e.changedTouches[0].pageX;
                const clickPosY = e.changedTouches[0].pageY;
                this.closeContextMenuOnceCallback();
                this.openContextMenuOnceCallback(taskElement, section, row, clickPosX, clickPosY);
            }
        } else {
            this.closeContextMenuOnceCallback();
        }
    }
}

class App {
    constructor() {
        this.endpoint = {
            address: "http://192.168.1.66:8080",
            uri: {
                documentOverviewHierarchical: "/document/overview/hierarchical/",
                documentOverviewChronological: "/document/overview/chronological/",
            },
        };

        this.state = {
            documentMode: undefined,
        };

        this.userInputManager = new UserInputResolver();
        this.userInputManager.openContextMenuOnceCallback = this.openContextMenuOnce.bind(this);
        this.userInputManager.closeContextMenuOnceCallback = this.closeContextMenuOnce.bind(this);

        this.modeSelector = new ModeSelector(this.updateMode.bind(this));
        this.infiniteSheet = new InfiniteSheet();
        this.contextMenu = new ContextMenu();

        // prettier-ignore
        domElementReuseCollector.registerItemIdentifier("infiniteSheetRow", function () {
            const cell = new InfiniteSheetTask();
            adoption(this.infiniteSheet.anchorPosition, [cell.container]);            
            return cell;
        }.bind(this));

        // prettier-ignore
        domElementReuseCollector.registerItemIdentifier("infiniteSheetHeader", function () {
            const cell = new InfiniteSheetHeader();
            adoption(this.infiniteSheet.anchorPosition, [cell.container]);
            return cell;
        }.bind(this));

        // prettier-ignore
        addEventListenerForNonTouchScreen(document, "contextmenu", this.userInputManager.contextMenuEventReceiver.bind(this.userInputManager));
        // prettier-ignore
        addEventListenerForNonTouchScreen(document, "click", this.userInputManager.clickEventReceiverNonTouchScreen.bind(this.userInputManager));
        // prettier-ignore
        executeWhenDocumentIsReady(
            function () {
                document.addEventListener("touchstart", this.userInputManager.touchStartEventReceiver.bind(this.userInputManager), false);
                document.addEventListener("touchmove", this.userInputManager.touchMoveEventReceiver.bind(this.userInputManager), false);
                document.addEventListener("touchend", this.userInputManager.touchEndEventReceiver.bind(this.userInputManager), false);
            }.bind(this)
        );

        document.addEventListener("scroll", this.scrollEventReceiver.bind(this));

        // const documentId = "61bbc44a-c61c-4d49-8804-486181081fa7";
        // this.fetchDocumentFromServer(documentId);
    }

    updateMode(newMode) {
        this.state.documentMode = newMode;
        this.updateView();
    }

    updateView() {
        this.infiniteSheet.build();
    }

    build() {
        this.infiniteSheet.build();
    }

    /** @param {MouseEvent} e */
    openContextMenuOnce(activatedTaskElement, taskSection, taskRow, clickPosX, clickPosY) {
        if (this.isContextMenuOpen) return;
        this.isContextMenuOpen = true;

        const contextMenuBounds = this.contextMenu.container.getBoundingClientRect();
        const contextMenuWidth = Math.floor(contextMenuBounds.width);
        const contextMenuHeight = Math.floor(contextMenuBounds.height);
        const padding = 20;

        const lastSafePosY = window.scrollY + window.innerHeight - (contextMenuHeight + padding);
        const lastSafePosX = window.innerWidth - (contextMenuWidth + padding);
        const posY = clickPosY < lastSafePosY ? clickPosY : lastSafePosY;
        const posX = clickPosX < lastSafePosX ? clickPosX : lastSafePosX;

        const transformOriginX = Math.floor(((clickPosX - posX) / contextMenuWidth) * 100);
        const transformOriginY = Math.floor(((clickPosY - posY) / contextMenuHeight) * 100);

        // const section = taskElement.dataset["section"];
        // const row = taskElement.dataset["row"];

        this.contextMenu.setPosition(posX, posY);
        this.contextMenu.setTransformOrigin(transformOriginX, transformOriginY);
        this.contextMenu.show();

        this.infiniteSheet.container.classList.add("context-menu-open");
        activatedTaskElement.classList.add("context-menu-focused-on");

        this.taskElementThatIsContextMenuFocusedOn = activatedTaskElement;
    }

    closeContextMenuOnce() {
        if (!this.isContextMenuOpen) return;

        this.contextMenu.hide();

        this.infiniteSheet.container.classList.remove("context-menu-open");
        this.taskElementThatIsContextMenuFocusedOn.classList.remove("context-menu-focused-on");

        this.isContextMenuOpen = false;
        this.taskElementThatIsContextMenuFocusedOn = undefined;
    }

    infiniteSheetContextMenuEventReceiver(e) {
        e.preventDefault();
        this.closeContextMenuOnce();
    }

    scrollEventReceiver(e) {
        this.closeContextMenuOnce();
    }

    fetchDocumentFromServer(documentId) {
        fetch(
            this.endpoint.address +
                this.endpoint.uri.documentOverviewChronological +
                documentId +
                "?limit=1000&offset=0"
        )
            .then((result) => result.json())
            .then(
                (result) => {
                    // Create < associative array || key-value list >

                    var dsd = new ChronologicalDocumentOverview(result);
                    console.log(dsd);
                    var tasks = {};
                    result.resource.forEach((resource) => {
                        tasks[resource.task_id] = new Task(resource);
                    });

                    var chronologicalOrdering = classifyTasksByDays(tasks);

                    // this.setState((state, props) => ({
                    //     overviewIsLoaded: true,
                    //     response: result,
                    //     tasks: tasks,
                    //     chronologicalOrdering: chronologicalOrdering,
                    // }));
                },
                (error) => {
                    // this.setState((state, props) => ({
                    //     overviewIsLoaded: true,
                    //     error: error,
                    // }));
                }
            );
    }
}

executeWhenDocumentIsReady(function () {
    const app = new App();
    app.build();
    // const body = document.getElementsByTagName("body")[0];
    // body.appendChild(app.container)
});
