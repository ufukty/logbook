import { adoption, createElement, domElementReuseCollector, executeWhenDocumentIsReady } from "./utilities.js";

import ModeSelector from "./viewControllers/ModeSelector.js";
import InfiniteSheet from "./viewControllers/InfiniteSheet.js";
import ContextMenu from "./viewControllers/ContextMenu.js";
import InfiniteSheetTask from "./viewControllers/InfiniteSheetTask.js";
import InfiniteSheetHeader from "./viewControllers/InfiniteSheetHeader.js";

class App {
    constructor() {
        this.state = {
            documentMode: undefined,
        };

        this.modeSelector = new ModeSelector(this.updateMode.bind(this));
        this.infiniteSheet = new InfiniteSheet();
        this.contextMenu = new ContextMenu();

        // prettier-ignore
        domElementReuseCollector.registerItemIdentifier("infiniteSheetRow", function () {
            const cell = new InfiniteSheetTask();
            adoption(this.infiniteSheet.anchorPosition, [cell.container]);            
            cell.container.addEventListener("contextmenu", this.contextmenuEventReceiver.bind(this));
            return cell;
        }.bind(this));
        // prettier-ignore
        domElementReuseCollector.registerItemIdentifier("infiniteSheetHeader", function () {
            const cell = new InfiniteSheetHeader();
            adoption(this.infiniteSheet.anchorPosition, [cell.container]);
            return cell;
        }.bind(this));

        document.addEventListener("click", this.clickEventReceiver.bind(this));
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

    /**
     * @param {HTMLElement} taskElement
     * @param {int} section
     * @param {int} row
     */
    openContextMenuOnTask(taskElement, section, row) {}

    contextmenuEventReceiver(e) {
        e.preventDefault();
        this.closeContextMenuOnce();

        const taskElement = e.currentTarget;
        const posY = e.pageY;
        const posX = e.pageX;
        const section = taskElement.dataset["section"];
        const row = taskElement.dataset["row"];

        this.contextMenu.setPosition(posX, posY);
        this.contextMenu.show();

        this.infiniteSheet.container.classList.add("context-menu-open");
        taskElement.classList.add("context-menu-focused-on");

        this.taskElementThatIsContextMenuFocusedOn = taskElement;
    }

    closeContextMenuOnce() {
        if (typeof this.taskElementThatIsContextMenuFocusedOn === "undefined") {
            return;
        }

        this.contextMenu.hide();

        this.infiniteSheet.container.classList.remove("context-menu-open");
        this.taskElementThatIsContextMenuFocusedOn.classList.remove("context-menu-focused-on");
    }

    clickEventReceiver(e) {
        e.preventDefault();
        this.closeContextMenuOnce();
    }
}

executeWhenDocumentIsReady(function () {
    const app = new App();
    app.build();
    // const body = document.getElementsByTagName("body")[0];
    // body.appendChild(app.container)
});
