import { addEventListenerForNonTouchScreen, domCollector, executeWhenDocumentIsReady } from "./bjsl/utilities.js";

import { ModeSelector } from "./viewControllers/ModeSelector.js";
import { InfiniteSheet } from "./viewControllers/InfiniteSheet.js";
import { ContextMenu } from "./viewControllers/ContextMenu.js";
import { UserInputResolver } from "./userInputResolver.js";
import { DataSource } from "./dataSource.js";

import { Task } from "./models/Task.js";
import { ChronologicalDocumentOverview } from "./models/DocumentOverviewChronological.js";

import { classifyTasksByDays } from "./dateTime.js";

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
        this.userInputManager.delegates = {
            closeContextMenu: this.closeContextMenuOnce.bind(this),
            openContextMenu: this.openContextMenuOnce.bind(this),
        };

        this.modeSelector = new ModeSelector(this.updateMode.bind(this));
        this.infiniteSheet = new InfiniteSheet();
        this.infiniteSheetWrapper = document.getElementById("infinite-sheet");
        this.infiniteSheetWrapper.appendChild(this.infiniteSheet.container);

        // // prettier-ignore
        // domCollector.registerItemIdentifier("infiniteSheetRow", function () {
        //     const cell = new InfiniteSheetTask();
        //     adoption(this.infiniteSheet.anchorPosition, [cell.container]);
        //     return cell;
        // }.bind(this));
        // // prettier-ignore
        // domCollector.registerItemIdentifier("infiniteSheetHeader", function () {
        //     const cell = new InfiniteSheetHeader();
        //     adoption(this.infiniteSheet.anchorPosition, [cell.container]);
        //     return cell;
        // }.bind(this));

        this.contextMenu = new ContextMenu();
        this.contextMenu.delegates = {
            delete: this.deleteTask.bind(this),
        };

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

        this.dataSource = new DataSource();
        this.dataSource.delegates.placementUpdate.push(this.placementUpdateFromData.bind(this));
        this.dataSource.delegates.objectUpdate.push((a) => {
            this.infiniteSheet.requestContentUpdateForObjectsIfNecessary(a);
        });
        // TODO:
        // TODO: this.dataSource.delegates.notifyGUI = this.infiniteSheet.updateData // updateData should be a callable
        // TODO: this.infiniteSheet.delegates.notifyDataSource = this.dataSource.updateData // updateData should be a callable
        this.infiniteSheet.config.structuredDataMedium = this.dataSource.medium.data; // FIXME:
        this.infiniteSheet.dataSource = this.dataSource;
        // FIXME: connect to server and fetch document/placement details
        this.dataSource.loadTestDataset2();

        // localSourceOfTruth.delegates.linearizedHierarchicalOrdering = this.hierarchicalOrderingUpdate.bind(this);
    }

    updateMode(newMode) {
        this.state.documentMode = newMode;
        this.updateView();
    }

    updateView() {
        this.infiniteSheet.build();
    }

    build() {
        // this.infiniteSheet.build();
    }

    /** @param {MouseEvent} e */
    openContextMenuOnce(taskPositionerElement, taskElement, taskId, taskSection, taskRow, clickPosX, clickPosY) {
        if (this.isContextMenuOpen) return;
        this.isContextMenuOpen = true;

        console.log(taskPositionerElement);

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

        this.contextMenu.setPosition(posX, posY);
        this.contextMenu.setTransformOrigin(transformOriginX, transformOriginY);
        this.contextMenu.setActiveTaskId(taskPositionerElement.parentNode.dataset["object-id"]);
        this.contextMenu.show();

        this.infiniteSheetWrapper.classList.add("context-menu-open");
        taskPositionerElement.classList.add("context-menu-focused-on");
    }

    closeContextMenuOnce(taskPositionerElement) {
        if (!this.isContextMenuOpen) return;

        this.contextMenu.hide();

        this.infiniteSheetWrapper.classList.remove("context-menu-open");
        taskPositionerElement.classList.remove("context-menu-focused-on");

        this.isContextMenuOpen = false;
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

    hierarchicalOrderingUpdate() {}

    placementUpdateFromData() {
        this.infiniteSheet.config.placement = this.dataSource.medium.data;
        this.infiniteSheet.updateView();
    }

    /** @param {string} taskId */
    deleteTask(taskId) {
        this.dataSource.medium.deleteRow(taskId);
        this.infiniteSheet.updateView();
    }
}

executeWhenDocumentIsReady(function () {
    const app = new App();
    app.build();
    // const body = document.getElementsByTagName("body")[0];
    // body.appendChild(app.container)
});
