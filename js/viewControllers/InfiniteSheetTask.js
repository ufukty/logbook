import { adoption, createElement } from "../baja.sl/utilities.js";
import { AbstractTableCellViewController } from "../baja.sl/AbstractTableCellViewController.js";
import { DataSource } from "../dataSource.js";

export class InfiniteSheetTask extends AbstractTableCellViewController {
    constructor() {
        super();
        this.dom = {
            container: createElement("div", ["infinite-sheet-task-container"]),
            taskBody: createElement("div", ["infinite-sheet-task-body", "preload"]),
            editableArea: createElement("div", ["infinite-sheet-task-editable-area"]),
            updateBadge: createElement("div", ["infinite-sheet-task-update-badge"]),
        };
        // prettier-ignore
        adoption(this.dom.container, [
            adoption(this.dom.taskBody, [
                adoption(this.dom.editableArea),
            ]),
            // adoption(this.dom.updateBadge)
        ])

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

        this._setDefaultState();
    }

    _setDefaultState() {
        this.state = {
            ...this.state,
            itemSymbol: undefined,
            updateCount: 0,
            updateBadgeIsVisible: false,
        };
    }

    prepareForFree() {
        this.dom.editableArea.innerText = "";
        this.dom.updateBadge.style.visibility = "hidden";
        if (this.animations.highlight) this.animations.highlight.cancel();
        this._setDefaultState();
    }

    prepareForUse() {}

    // /**
    //  * This method, restores details presented previously by this (or another
    //  * cell) to user again, without animation/transition. This method should
    //  * be called by AbstractTableViewController implementer.
    //  * @param {Symbol} itemSymbol
    //  * @param {DataSource} dataSource
    //  */
    // restoreItem(itemSymbol, dataSource) {
    //     const taskDetails = dataSource.cache.tasks.get(itemSymbol)
    //     taskDetails.

    //     this.state.itemSymbol = itemSymbol;
    //     this.state.updateCount = dataSource.cache.updateCounts.get(itemSymbol) ?? 0;
    //     if (dataSource.cache.updateCounts.get(itemSymbol) > 0) {
    //         this.dom.updateBadge.innerText = nextUpdateCount.toString();
    //     }
    // }

    setContent(newContent) {
        this.dom.editableArea.scrollLeft = 0;
        this.dom.editableArea.innerText = newContent;
    }

    setDegree(degree) {}

    setDepth(depth) {
        this.dom.container.style.transform = `translateX(${this.config.translationForDepth * depth}px)`;
    }

    enableEditMode() {
        this.dom.editableArea.contentEditable = true;
        this.dom.editableArea.focus();
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
        this.animations.highlight = this.dom.editableArea.animate([{ opacity: 0 }, { opacity: 1 }], {
            duration: 2000,
            fill: "none",
            iterations: 1,
            easing: "cubic-bezier(0.1, 0.6, 1.0, 1.0)",
        });
    }
}
