import { addEventListenerForNonTouchScreen, adoption, createElement } from "../baja.sl/utilities.js";
import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";

export class InfiniteSheetTask extends AbstractTableCellViewController {
    constructor() {
        super();
        this.dom = {
            container: createElement("div", ["task-container"]),
            updateBadge: createElement("div", ["task-update-badge"]),
            borderedArea: createElement("div", ["task-bordered-area"]),
            taskBody: createElement("div", ["task-body", "preload"]),
            textScroller: createElement("div", ["task-text-scroller"]),
            textArea: createElement("div", ["task-text-area"]),
            detailsScroller: createElement("div", ["task-details-scroller"]),
            detailsContainer: createElement("div", ["task-details-container"]),
            deadlineContainer: createElement("div", ["task-detail-container", "task-deadline-container"]),
            collaboratorAdditiveListContainer: createElement("div", [
                "task-detail-container",
                "task-collaborator-additive-container",
            ]),
            collaboratorExludeListContainer: createElement("div", [
                "task-detail-container",
                "task-collaborator-exclude-container",
            ]),
        };

        adoption(this.dom.container, [
            adoption(this.dom.borderedArea, [
                adoption(this.dom.taskBody, [adoption(this.dom.textScroller, [adoption(this.dom.textArea)])]),
                adoption(this.dom.detailsScroller, [
                    adoption(this.dom.detailsContainer, [
                        adoption(this.dom.deadlineContainer),
                        adoption(this.dom.collaboratorAdditiveListContainer),
                        adoption(this.dom.collaboratorExludeListContainer),
                    ]),
                ]),
            ]),
            adoption(this.dom.updateBadge),
        ]);

        this.hideControls();

        this.dom.deadlineContainer.innerText = "NRA: 12.11.2023";
        this.dom.collaboratorAdditiveListContainer.innerText = "Add to collaboration: () () ()";
        this.dom.collaboratorExludeListContainer.innerText = "Exclude from collaboration: () () ()";

        this.dom.updateBadge.innerText = "0";

        // this.dom.editableArea.contentEditable = true;

        this.config = {
            ...this.config,
            translationForDepth: 20,
        };

        this.animations = {
            /** @type {Animation} */
            highlight: undefined,
        };

        this.state = this._returnDefaultState();

        addEventListenerForNonTouchScreen(this.dom.taskBody, "click", this._clickEventReceiver.bind(this));
        this.dom.taskBody.addEventListener("touchend", this._clickEventReceiver.bind(this));
    }

    _clickEventReceiver() {
        if (this.state.isControlsPresented) this.hideControls();
        else this.showControls();
    }

    showControls() {
        this.dom.container.dataset.presentControls = "true";
        this.state.isControlsPresented = true;
    }

    hideControls() {
        this.dom.container.dataset.presentControls = "false";
        this.dom.detailsScroller.scrollLeft = 0;
        this.state.isControlsPresented = false;
    }

    _returnDefaultState() {
        return {
            ...this.state,
            itemSymbol: undefined,
            updateCount: 0,
            updateBadgeIsVisible: false,
            isControlsPresented: false,
        };
    }

    prepareForFree() {
        this.dom.textArea.innerText = "";
        this.dom.updateBadge.innerText = "0";
        this.dom.updateBadge.style.visibility = "hidden";
        this.hideControls();
        if (this.animations.highlight) this.animations.highlight.cancel();
        this.state = this._returnDefaultState();
    }

    prepareForUse() {}

    setContent(newContent) {
        this.dom.textArea.scrollLeft = 0;
        this.dom.textArea.innerText = newContent;
    }

    setDegree(degree) {}

    setDepth(depth) {
        this.dom.container.style.transform = `translateX(${this.config.translationForDepth * depth}px)`;
    }

    enableEditMode() {
        this.dom.textArea.contentEditable = true;
        this.dom.textArea.focus();
    }

    /**
     * @param {number} nextUpdateCount
     * @param {boolean} withAnimation
     */
    setUpdateCount(nextUpdateCount, withAnimation = true) {
        this.state.updateCount = nextUpdateCount;

        if (nextUpdateCount === 0 && this.state.updateBadgeIsVisible) {
            this.dom.updateBadge.innerText = nextUpdateCount.toString();

            this.state.updateBadgeIsVisible = false;

            if (withAnimation) {
                this.dom.updateBadge
                    .animate(
                        [
                            { transformOrigin: "-50% 50%", transform: "scale(1)", opacity: "1" },
                            { transformOrigin: "-50% 50%", transform: "scale(0.5)", opacity: "0" },
                        ],
                        {
                            iterations: 1,
                            duration: 200,
                            fill: "forwards",
                            easing: "ease-out",
                        }
                    )
                    .finished.then(() => {
                        this.dom.updateBadge.style.visibility = "hidden";
                    });
            } else {
                this.dom.updateBadge.style.visibility = "hidden";
            }
        } else if (nextUpdateCount > 0 && !this.state.updateBadgeIsVisible) {
            this.dom.updateBadge.innerText = nextUpdateCount.toString();
            this.dom.updateBadge.style.visibility = "visible";
            this.state.updateBadgeIsVisible = true;

            const computedProps = getComputedStyle(this.dom.updateBadge);
            const computedWidth = parseFloat(computedProps.getPropertyValue("width"));
            const computedHeight = parseFloat(computedProps.getPropertyValue("height"));

            this.dom.updateBadge.style.minWidth = `${Math.max(computedWidth, computedHeight)}px`;
            this.dom.updateBadge.style.borderRadius = `${computedHeight / 2}px`;

            if (withAnimation) {
                this.dom.updateBadge.animate(
                    [
                        { transformOrigin: "-50% 50%", transform: "scale(0.5)", opacity: "0" },
                        { transformOrigin: "-50% 50%", transform: "scale(1)", opacity: "1" },
                    ],
                    {
                        iterations: 1,
                        duration: 200,
                        fill: "forwards",
                        easing: "cubic-bezier(0.4, 1.5, 0.8, 1.0)",
                    }
                );
            } else {
                this.dom.updateBadge.style.opacity = "1";
            }
        } else {
            this.dom.updateBadge.innerText = nextUpdateCount.toString();
            if (withAnimation) {
                this.dom.updateBadge.animate(
                    [
                        { transformOrigin: "50% 50%", transform: "scale(1.2)" },
                        { transformOrigin: "50% 50%", transform: "scale(1)" },
                    ],
                    {
                        duration: 300,
                        fill: "none",
                        easing: "ease-out",
                        iterations: 1,
                    }
                );
            }
        }
    }

    highlight() {
        this.animations.highlight = this.dom.textArea.animate([{ opacity: 0 }, { opacity: 1 }], {
            duration: 2000,
            fill: "none",
            iterations: 1,
            easing: "cubic-bezier(0.1, 0.6, 1.0, 1.0)",
        });
        this.animations.highlight.onfinish = () => {
            this.animations.highlight.cancel();
        };
    }
}
