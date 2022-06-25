import { adoption, domElementReuseCollector, createElement, toggleAnimationWithClass } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";

function inBetween(a, b, c) {
    if (a <= b && c <= c) return true;
    else return false;
}

function checkCollusion(item_y1, item_y2, viewport_y1, viewport_y2) {
    /*
            * * * * * * *  (y1)                     * * * * * * *  (y1)                
            *           *                           *           *              
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y1)
        +   *           *                           *           *           +
        +   * * * * * * *  (y2)                     *           *           +
        +                                           *           *           +
        +                                           *           *           +       <=  viewport
        +                                           *           *           +
        +                 * * * * * * *  (y1)       *           *           +
        +                 *           *             *           *           +
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y2)
                          *           *             *           *                              
                          * * * * * * *  (y2)       * * * * * * *  (y2)                               
    */
    if (item_y2 < viewport_y1 || item_y1 > viewport_y2)
        // if item starts after viewport ends, or item ends before viewport starts,
        // then the item is not in viewport.
        return false;
    else return true;
}

/**
 * @returns {Set.<string>}
 */
function mergeMapKeys() {
    console.log(arguments);
    debugger;
    let set_ = new Set();
    for (let i = 0; i < arguments.length; i++) {
        for (const key of arguments[i].keys()) set_.add(key);
    }
    return set_;
}

export class AbstractTableViewController extends AbstractViewController {
    constructor() {
        super();
        this.container = document.getElementById("infinite-sheet");
        this.contentArea = createElement("div", ["content-area"]);
        this.anchorPosition = createElement("div", ["anchor-position"]);
        adoption(this.container, [adoption(this.contentArea, [this.anchorPosition])]);

        // this.allocatedSectionElements = {};
        // this.allocatedItemElements = {};
        // this.visibleHeaderElements = {};
        // this.visibleRowElements = {};
        // this.computedHeights = {};

        this.config = {
            /** AbstractDataSource */
            dataSourceRef: undefined,
            margins: {
                pageContent: {
                    before: 10,
                    after: 10,
                },
                section: {
                    before: 10,
                    between: 10,
                },
                row: {
                    before: 10,
                    between: 10,
                },
            },
            /** @type { Map.<Symbol, number> } */
            defaultHeightForReuseId: new Map(),
            /** The ordering of sections and rows in them.
             * Each `Symbol` represents an `objectID`
             * (either a `sectionID` or `rowID`). */
            placement: {
                /** @type { Array.<Symbol> } */
                sections: [],
                /** @type { Map.<Symbol, Array.<Symbol>> } */
                rows: new Map(),
            },
            /**
             * @type { Map.<Symbol, Symbol> }
             * Maps `objectIDs` to correct reuse identifiers.
             * Information will be used for requesting and
             * sending cells to `domElementReuseCollector`.
             * > **Note that**: Related constructors for each id
             * given as key to this map, should've
             * registered to `domElementReuseCollector` already.
             */
            objectReuseIdentifiers: new Map(),
        };

        this.computedValues = {
            /**
             * @type { Map.<Symbol, AbstractTableCellViewController> }
             * Bind objectID and allocated cell on domElementReuseCollector.get() call
             */
            allocatedCells: new Map(),
            pageHeight: undefined,
            /** set and use when nodes above viewport changes their sizings */
            scrollShift: undefined,
            /**  */
            positions: {
                /** @type {Map.<Symbol, { starts: number, ends: number, height: number }>} */
                current: new Map(),
                /** @type {Map.<Symbol, { starts: number, ends: number, height: number }>} */
                next: undefined, // initialized later
            },
        };

        // document.addEventListener("scroll", this.renderVisible.bind(this));

        // this.resizeObserver = new ResizeObserver((entries) => {
        //     entries.forEach((entry) => {
        //         const element = entry.target;
        //         const section = element.dataset.section;
        //         const row = element.dataset.row;
        //         const height = entry.contentRect.height;
        //         this.updateComputedHeightOfElement(this.getReferenceOfAllocatedRowElement(section, row), section, row);
        //     });
        //     this.calculateElementBounds();
        //     this.rePosition();
        //     // this.renderVisible();
        // });
    }

    /**
     * @param {{sections: [], rows: {}}} placements
     * @param {{}} heights
     * @param {{}}
     */
    calculateComponentPositions() {
        let lastPosition = 0;
        lastPosition += this.config.margins.pageContent.before;

        // for each section
        for (const [sectionIndex, sectionID] of this.config.placement.sections.entries()) {
            // spacing before & between sections
            if (sectionIndex === 0) lastPosition += this.config.margins.section.before;
            else lastPosition += this.config.margins.section.between;

            lastPosition += this.config.margins.section.before;

            const headerHeight = heights[sectionID];

            // save object positions
            this.computedValues.positions.next.set(sectionID, {
                starts: lastPosition,
                ends: lastPosition + headerHeight,
                height: headerHeight,
            });

            lastPosition += headerHeight;

            for (const [rowIndex, rowID] of this.config.placement.rows.get(sectionID).entries()) {
                // spacing before & between rows
                if (rowIndex === 0) lastPosition += this.config.margins.row.before;
                else lastPosition += this.config.margins.row.between;

                const itemHeight = heights[rowID];

                // save object positions
                this.computedValues.positions.next.set(rowID, {
                    starts: lastPosition,
                    ends: lastPosition + itemHeight,
                    height: itemHeight,
                });

                lastPosition += itemHeight;
            }
        }

        lastPosition += this.config.margins.pageContent.after;
        this.computedValues.pageHeight = lastPosition;
    }

    /**
     * Compares previous and next states for positioning, existance, folding etc.
     * Then classifies tasks by operations to perform on them.
     * @param {Map.<string, {x: number, y: number, height: number}>} positions_current
     * @param {Map.<string, {x: number, y: number, height: number}>} positions_next
     * @param {Set} foldObject_current
     * @param {Set} foldObjects_next
     * @param {Set} objectsInViewbox_current
     * @param {Set} objectsInViewbox_next
     * @param {{viewport: {starts: number, ends: number}, preloadArea: {starts: number, ends: number}, releaseArea: {starts: number, ends: number}, focusPoint: number}} viewportPositions
     */
    updateComponents(
        foldObject_current,
        foldObjects_next,
        constructedObjects_current,
        constructedObjects_next,
        objectsInDestructionArea_next,
        objectsInViewbox_current,
        objectsInViewbox_next,
        viewportPositions
    ) {
        let mergedObjectIDs = mergeMapKeys(this.computedValues.positions.current, this.computedValues.positions.next);
        let totalScrollShift = 0;
        // debugger;

        for (const objectID of mergedObjectIDs) {
            // calculate scroll shift; if there is any change in height of object
            if (positions_current.has(objectID) && positions_next.has(objectID)) {
                if (positions_current.get(objectID).height !== positions_next.get(objectID).height) {
                    if (positions_next.get(objectID).y < viewportPositions.focusPoint) {
                        totalScrollShift +=
                            positions_next.get(objectID).height - positions_current.get(objectID).height;
                    }
                }
            }

            let waitForTransitionEnd = false;
            /** @type {AbstractTableCellViewController} object */
            let object;

            // "to construct"
            if (!constructedObjects_current.has(objectID) && constructedObjects_next.has(objectID)) {
                object = domElementReuseCollector.get(this.objectReuseIdentifiers[objectID]);

                let objectInitializationPositionY, objectInitializationPositionX;
                if (positions_current.has(objectID)) {
                    objectInitializationPositionY = positions_current.get(objectID).y;
                    objectInitializationPositionX = positions_current.get(objectID).x;
                } else if (positions_next.has(objectID)) {
                    objectInitializationPositionY = positions_next.get(objectID).y;
                    objectInitializationPositionX = positions_next.get(objectID).x;
                }
                object.setPositionY(objectInitializationPositionY, false);
                object.setPositionX(objectInitializationPositionX, false);
                object.setContent(this.config.dataSourceRef.getTextContent(objectID));
            }

            // "to appear"
            if (!objectsInViewbox_current.has(objectID) && objectsInViewbox_next.has(objectID)) {
            }

            // existance change
            if (positions_current.has(objectID) && !positions_next.has(objectID)) {
                // "to delete" content
                // TODO:
                waitForTransitionEnd = true;
            } else if (positions_current.has(objectID) && !positions_next.has(objectID)) {
                // "to create" content
                // TODO:
                waitForTransitionEnd = true;
            }

            // position change
            if (positions_current.get(objectID).y !== positions_next.get(objectID).y) {
                // TODO:
                waitForTransitionEnd = true;
                object.setPosition(positions_next.get(objectID).y, true);
            }

            // folding change
            if (foldObject_current.has(objectID) && !foldObjects_next.has(objectID)) {
                // unfold
                // TODO:
                object.unfold();
            } else if (!foldObject_current.has(objectID) && foldObjects_next.has(objectID)) {
                // fold
                // TODO:
                object.fold();
            }

            // to disappear
            if (objectsInViewbox_current.has(objectID) && !objectsInViewbox_next.has(objectID)) {
            }

            // "to destruct"
            if (objectsInDestructionArea_next.has(objectID)) {
                var f = () => {
                    domElementReuseCollector.free(this.objectReuseIdentifiers[objectID], object);
                };
                if (waitForTransitionEnd) object.container.addEventListener("transitionend", f, { once: true });
                else f();
            }
        }
    }

    updateContainer() {
        this.container.style.height = `${this.computedValues.pageHeight}px`;
    }

    updateViewFromData() {
        // debugger;
        this.computedValues.positions.next = new Map();

        this.calculateComponentPositions();
        this.updateContainer();
        this.updateComponents();

        delete this.computedValues.positions.current; // forget positions computed on previous call
        this.computedValues.positions.current = this.computedValues.positions.next;
    }

    build() {
        this.updateViewFromData();
    }
}
