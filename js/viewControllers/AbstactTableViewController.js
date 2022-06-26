import { adoption, domElementReuseCollector, createElement, toggleAnimationWithClass } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { TableViewStructuredDataMedium } from "../dataSource.js";

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

class BasicTableCellContainerViewController extends AbstractViewController {
    constructor() {
        super();
        this.container = createElement("div", ["table-view-cell-container"]);
        this.userProvidedCell = undefined; // should be assigned by callee
    }

    prepareForFree() {
        this.userProvidedCell.prepareForFree();
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation) {
        // TODO: withAnimation option
        this.container.style.top = `${newPosition}px`;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionX(newPosition, withAnimation) {
        // TODO: withAnimation option
        this.container.style.left = `${newPosition}px`;
    }
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
            /** @type { TableViewStructuredDataMedium } */
            structuredDataMedium: undefined,
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
             * Each `Symbol` represents an `objectSymbol`
             * (either a `sectionID` or `rowID`). */
            placement: {
                /** @type { Array.<Symbol> } */
                sections: [],
                /** @type { Map.<Symbol, Array.<Symbol>> } */
                rows: new Map(),
            },
            /**
             * @type { Map.<Symbol, Symbol> }
             * Maps `objectISymbol` to correct reuse identifiers.
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
             * Bind objectSymbol and allocated cell on domElementReuseCollector.get() call
             */
            computedHeights: new Map(),
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
     * Embeds the user supplied cell constructor with a function
     * that creates custom positioner view controller and wraps
     * the cell returned by user-supplied cell constructor with it.
     */
    registerCellConstructor(reuseIdentifier, cellConstructor) {
        domElementReuseCollector.registerItemIdentifier(reuseIdentifier, () => {
            const userProvidedCell = cellConstructor();
            const cellContainer = new BasicTableCellContainerViewController();
            cellContainer.userProvidedCell = userProvidedCell;
            // prettier-ignore
            adoption(this.anchorPosition, [
                adoption(cellContainer.container, [
                    userProvidedCell.container
                ])
            ])
            // this.resizeObserver.observe(userProvidedCell.container);
        });
    }

    /**
     * When user request a cell to populate with data, this method
     * only sends the nested user-supplied custom cell, instead
     * the positioner cell that wraps it from the constructor
     * registered by .registerCellConstructor().
     * @returns { AbstractTableCellViewController }
     * */
    getRecycledCell(reuseIdentifier) {
        const cellContainer = domElementReuseCollector.get(reuseIdentifier);
        return cellContainer.userProvidedCell;
    }

    /**
     * Default height is important to estimate overall height of
     * the page and make the scrollbar much more useful.
     * @param { number } objectSymbol
     * @returns { number }
     */
    getDefaultHeightOfObject(objectSymbol) {
        console.error("abstract function is called");
    }

    /**
     * User should implement this method.
     * Request an empty cell from .getFreeCell()
     * with previously registered reuseIdentifier
     * Then populate content accordingly to
     * specified objectSymbol.
     * @returns { AbstractTableCellViewController }
     */
    getCellForObject(objectSymbol) {
        console.error("abstract function is called");
    }

    /**
     * @param {{sections: [], rows: {}}} placements
     * @param {{}} heights
     * @param {{}}
     */
    calculateComponentPositions() {
        // TODO: don't mind tasks that their parents are folded.
        let lastPosition = 0;
        lastPosition += this.config.margins.pageContent.before;

        // for each section
        for (const [sectionIndex, sectionID] of this.config.placement.sections.entries()) {
            // spacing before & between sections
            if (sectionIndex === 0) lastPosition += this.config.margins.section.before;
            else lastPosition += this.config.margins.section.between;

            lastPosition += this.config.margins.section.before;

            const headerHeight = this.computedValues.computedHeights.has(sectionID)
                ? this.computedValues.computedHeights.get(sectionID)
                : this.getDefaultHeightOfObject(sectionID);

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

                const itemHeight = this.computedValues.computedHeights.has(rowID)
                    ? this.computedValues.computedHeights.get(rowID)
                    : this.getDefaultHeightOfObject(rowID);

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

        for (const objectSymbol of mergedObjectIDs) {
            // calculate scroll shift; if there is any change in height of object
            if (positions_current.has(objectSymbol) && positions_next.has(objectSymbol)) {
                if (positions_current.get(objectSymbol).height !== positions_next.get(objectSymbol).height) {
                    if (positions_next.get(objectSymbol).y < viewportPositions.focusPoint) {
                        totalScrollShift +=
                            positions_next.get(objectSymbol).height - positions_current.get(objectSymbol).height;
                    }
                }
            }

            let waitForTransitionEnd = false;
            /** @type {AbstractTableCellViewController} object */
            let cell;

            // "to construct"
            if (!constructedObjects_current.has(objectSymbol) && constructedObjects_next.has(objectSymbol)) {
                cell = this.getCellForObjectId(objectSymbol);

                // let objectInitializationPositionY, objectInitializationPositionX;
                // if (positions_current.has(objectSymbol)) {
                //     objectInitializationPositionY = positions_current.get(objectSymbol).y;
                //     objectInitializationPositionX = positions_current.get(objectSymbol).x;
                // } else if (positions_next.has(objectSymbol)) {
                //     objectInitializationPositionY = positions_next.get(objectSymbol).y;
                //     objectInitializationPositionX = positions_next.get(objectSymbol).x;
                // }
                // object.setPositionY(objectInitializationPositionY, false);
                // object.setPositionX(objectInitializationPositionX, false);
                // object.setContent(this.config.structuredDataMedium.getTextContent(objectSymbol));
            }

            // "to appear"
            if (!objectsInViewbox_current.has(objectSymbol) && objectsInViewbox_next.has(objectSymbol)) {
            }

            // existance change
            if (positions_current.has(objectSymbol) && !positions_next.has(objectSymbol)) {
                // "to delete" content
                // TODO:
                waitForTransitionEnd = true;
            } else if (positions_current.has(objectSymbol) && !positions_next.has(objectSymbol)) {
                // "to create" content
                // TODO:
                waitForTransitionEnd = true;
            }

            // position change
            if (positions_current.get(objectSymbol).y !== positions_next.get(objectSymbol).y) {
                // TODO:
                waitForTransitionEnd = true;
                cell.setPosition(positions_next.get(objectSymbol).y, true);
            }

            // folding change
            if (foldObject_current.has(objectSymbol) && !foldObjects_next.has(objectSymbol)) {
                // unfold
                // TODO:
                cell.unfold();
            } else if (!foldObject_current.has(objectSymbol) && foldObjects_next.has(objectSymbol)) {
                // fold
                // TODO:
                cell.fold();
            }

            // to disappear
            if (objectsInViewbox_current.has(objectSymbol) && !objectsInViewbox_next.has(objectSymbol)) {
            }

            // "to destruct"
            if (objectsInDestructionArea_next.has(objectSymbol)) {
                var f = () => {
                    domElementReuseCollector.free(this.objectReuseIdentifiers[objectSymbol], cell);
                };
                if (waitForTransitionEnd) cell.container.addEventListener("transitionend", f, { once: true });
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
