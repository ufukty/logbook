import { adoption, domElementReuseCollector, createElement, toggleAnimationWithClass } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { TableViewStructuredDataMedium } from "../dataSource.js";

function inBetween(a, b, c) {
    if (a <= b && c <= c) return true;
    else return false;
}

function checkCollision(item_y1, item_y2, viewport_y1, viewport_y2) {
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
    let set_ = new Set();
    for (let i = 0; i < arguments.length; i++) {
        if (arguments[i]) for (const key of arguments[i].keys()) set_.add(key);
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
            current: {
                /** @type { Map.<Symbol, > } */
                computedHeights: new Map(),
                /**
                 * @type { Map.<Symbol, AbstractTableCellViewController> }
                 * Bind objectSymbol and allocated cell on domElementReuseCollector.get() call
                 */
                allocatedCells: new Map(),
                pageHeight: undefined,
                /** set and use when nodes above viewport changes their sizings */
                scrollShift: undefined,
                /** @type {Map.<Symbol, { starts: number, ends: number, height: number }>} */
                positions: new Map(),
                /**
                 * Holds the set of object symbols from current and next iterations of update.
                 * Intended to be used by update functions.
                 * @type {Set.<Symbol>}
                 */
                mergedObjectSymbols: undefined,
                /** @type { { inViewport: Set.<Symbol>, inPreload: Set.<Symbol>, inParking: Set.<Symbol> } } */
                zoneCollusions: {
                    inViewport: new Set(),
                    inPreload: new Set(),
                    inParking: new Set(),
                },
                /** @type { { viewport: { starts: number, ends: number }, preload: { starts: number, ends: number }, parking: { starts: number, ends: number } } } */
                boundaries: {
                    viewport: {},
                    /** an area which its height is 3 times of
                     * viewport (1 above, 1 below) */
                    preload: {},
                    /** an area which its height is 5 times of
                     * viewport (2 above, 2 below) */
                    parking: {},
                },
                classifiedObjects: {
                    toConstruct: new Set(),
                    toAppear: new Set(),
                    toUpdatePositionY: new Set(),
                    toUpdatePositionX: new Set(), // indentation
                    toUpdateFolding: new Set(),
                    toUpdateExistance: new Set(),
                    toDisappear: new Set(),
                    toDestruct: new Set(),
                },
            },
            next: {
                /** @type { Map.<Symbol, > } */
                computedHeights: new Map(),
                /**
                 * @type { Map.<Symbol, AbstractTableCellViewController> }
                 * Bind objectSymbol and allocated cell on domElementReuseCollector.get() call
                 */
                allocatedCells: new Map(),
                pageHeight: undefined,
                /** set and use when nodes above viewport changes their sizings */
                scrollShift: undefined,
                /** @type {Map.<Symbol, { starts: number, ends: number, height: number }>} */
                positions: new Map(),
                /**
                 * Holds the set of object symbols from current and next iterations of update.
                 * Intended to be used by update functions.
                 * @type {Set.<Symbol>}
                 */
                mergedObjectSymbols: undefined,
                /** @type { { inViewport: Set.<Symbol>, inPreload: Set.<Symbol>, inParking: Set.<Symbol> } } */
                zoneCollusions: {
                    inViewport: new Set(),
                    inPreload: new Set(),
                    inParking: new Set(),
                },
                /** @type { { viewport: { starts: number, ends: number }, preload: { starts: number, ends: number }, parking: { starts: number, ends: number } } } */
                boundaries: {
                    viewport: {},
                    /** an area which its height is 3 times of
                     * viewport (1 above, 1 below) */
                    preload: {},
                    /** an area which its height is 5 times of
                     * viewport (2 above, 2 below) */
                    parking: {},
                },
                classifiedObjects: {
                    toConstruct: new Set(),
                    toAppear: new Set(),
                    toUpdatePositionY: new Set(),
                    toUpdatePositionX: new Set(), // indentation
                    toUpdateFolding: new Set(),
                    toUpdateExistance: new Set(),
                    toDisappear: new Set(),
                    toDestruct: new Set(),
                },
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
     * Embeds the user-supplied cell constructor with a function
     * that creates a custom positioner view controller and wraps
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

    updateZoneBoundaries() {
        const preloadZoneOffset = 1 * window.innerHeight;
        const parkingZoneOffset = 2 * window.innerHeight;

        this.computedValues.next.boundaries = {
            viewport: {
                starts: window.scrollY,
                ends: window.scrollY + window.innerHeight,
            },
            /** an area which its height is 3 times of
             * viewport (1 above, 1 below) */
            preload: {
                starts: window.scrollY - preloadZoneOffset,
                ends: window.scrollY + window.innerHeight + preloadZoneOffset,
            },
            /** an area which its height is 5 times of
             * viewport (2 above, 2 below) */
            parking: {
                starts: window.scrollY - parkingZoneOffset,
                ends: window.scrollY + window.innerHeight + parkingZoneOffset,
            },
        };
    }

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

            const headerHeight = this.computedValues.next.computedHeights.has(sectionID)
                ? this.computedValues.next.computedHeights.get(sectionID)
                : this.getDefaultHeightOfObject(sectionID);

            // save object positions
            this.computedValues.next.positions.set(sectionID, {
                starts: lastPosition,
                ends: lastPosition + headerHeight,
                height: headerHeight,
            });

            lastPosition += headerHeight;

            for (const [rowIndex, rowID] of this.config.placement.rows.get(sectionID).entries()) {
                // spacing before & between rows
                if (rowIndex === 0) lastPosition += this.config.margins.row.before;
                else lastPosition += this.config.margins.row.between;

                const itemHeight = this.computedValues.next.computedHeights.has(rowID)
                    ? this.computedValues.next.computedHeights.get(rowID)
                    : this.getDefaultHeightOfObject(rowID);

                // save object positions
                this.computedValues.next.positions.set(rowID, {
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

    /** @access private */
    classifyObjectsByCollidedZones() {
        this.computedValues.next.zoneCollusions = {
            inViewport: new Set(),
            inPreload: new Set(),
            inParking: new Set(),
        };

        for (const [objectSymbol, objectPos] of this.computedValues.next.positions) {
            const inViewport = checkCollision(
                objectPos.starts,
                objectPos.ends,
                this.computedValues.next.boundaries.viewport.starts,
                this.computedValues.next.boundaries.viewport.ends
            );
            const inPreload = checkCollision(
                objectPos.starts,
                objectPos.ends,
                this.computedValues.next.boundaries.preload.starts,
                this.computedValues.next.boundaries.preload.ends
            );
            const inParking = checkCollision(
                objectPos.starts,
                objectPos.ends,
                this.computedValues.next.boundaries.parking.starts,
                this.computedValues.next.boundaries.parking.ends
            );
            if (inViewport) this.computedValues.next.zoneCollusions.inViewport.add(objectSymbol);
            if (inPreload) this.computedValues.next.zoneCollusions.inPreload.add(objectSymbol);
            if (inParking) this.computedValues.next.zoneCollusions.inParking.add(objectSymbol);
        }
    }

    mergeObjectSymbolsWithPreviousIteration() {
        this.computedValues.mergedObjectSymbols = mergeMapKeys(
            this.computedValues.current.positions,
            this.computedValues.next.positions
        );
    }

    calculateFocusShift() {
        let totalScrollShift = 0;
        for (const objectSymbol of this.computedValues.mergedObjectSymbols) {
            // calculate scroll shift; if there is any change in height of object
            if (
                this.computedValues.current.positions.has(objectSymbol) &&
                this.computedValues.next.positions.has(objectSymbol)
            ) {
                if (
                    this.computedValues.current.positions.get(objectSymbol).height !==
                    this.computedValues.next.positions.get(objectSymbol).height
                ) {
                    if (this.computedValues.next.positions.get(objectSymbol).y < viewportPositions.focusPoint) {
                        totalScrollShift +=
                            this.computedValues.next.positions.get(objectSymbol).height -
                            this.computedValues.current.positions.get(objectSymbol).height;
                    }
                }
            }
        }
        this.computedValues.scrollShift = totalScrollShift;
    }

    classifyComponentsByUpdateTypes() {
        for (const objectSymbol of this.computedValues.mergedObjectSymbols) {
            // "to construct"
            if (
                !this.computedValues.current.zoneCollusions.inPreload.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inPreload.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toConstruct.push(objectSymbol);
            }

            // "to appear"
            if (
                !this.computedValues.current.zoneCollusions.inViewport.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inViewport.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toAppear.push(objectSymbol);
            }

            // // existance change
            // if (
            //     this.computedValues.current.positions.has(objectSymbol) &&
            //     !this.computedValues.next.positions.has(objectSymbol)
            // ) {
            // } else if (
            //     this.computedValues.current.positions.get(objectSymbol) &&
            //     !this.computedValues.next.positions.has(objectSymbol)
            // ) {
            //     // "to create" content
            //     // TODO:
            //     waitForTransitionEnd = true;
            // }

            // // position change
            // if (
            //     this.computedValues.current.positions.get(objectSymbol).y !==
            //     this.computedValues.next.positions.get(objectSymbol).y
            // ) {
            //     // TODO:
            //     waitForTransitionEnd = true;
            //     cell.setPosition(this.computedValues.next.positions.get(objectSymbol).y, true);
            // }

            // // folding change
            // if (foldObject_current.has(objectSymbol) && !foldObjects_next.has(objectSymbol)) {
            //     // unfold
            //     // TODO:
            //     cell.unfold();
            // } else if (!foldObject_current.has(objectSymbol) && foldObjects_next.has(objectSymbol)) {
            //     // fold
            //     // TODO:
            //     cell.fold();
            // }

            // to disappear
            if (
                this.computedValues.current.zoneCollusions.inViewport.has(objectSymbol) &&
                !this.computedValues.next.zoneCollusions.inViewport.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toDisappear.push(objectSymbol);
            }

            // to destruct
            if (
                this.computedValues.current.zoneCollusions.inParking.has(objectSymbol) &&
                !this.computedValues.next.zoneCollusions.inParking.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toDestruct.push(objectSymbol);
            }
        }
    }

    updateComponents() {
        // "to construct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toConstruct) {
            cell = this.getCellForObject(objectSymbol);
            this.computedValues.allocatedCells.set(objectSymbol, cell);

            let objectInitializationPositionY, objectInitializationPositionX;
            if (this.computedValues.current.positions.has(objectSymbol)) {
                objectInitializationPositionY = this.computedValues.current.positions.get(objectSymbol).y;
                objectInitializationPositionX = this.computedValues.current.positions.get(objectSymbol).x;
            } else if (this.computedValues.next.positions.has(objectSymbol)) {
                objectInitializationPositionY = this.computedValues.next.positions.get(objectSymbol).y;
                objectInitializationPositionX = this.computedValues.next.positions.get(objectSymbol).x;
            }
            object.setPositionY(objectInitializationPositionY, false);
            object.setPositionX(objectInitializationPositionX, false);
            object.setContent(this.config.structuredDataMedium.getTextContent(objectSymbol));
        }

        // "to appear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toAppear) {
        }

        // existance change
        // position change
        // folding change

        // "to disappear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDisappear) {
        }

        // "to destruct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDestruct) {
        }
    }

    updateContainer() {
        this.container.style.height = `${this.computedValues.pageHeight}px`;
    }

    updateViewFromData() {
        this.computedValues.next = {
            /** @type { Map.<Symbol, number> } */
            computedHeights: new Map(),
            /**
             * @type { Map.<Symbol, AbstractTableCellViewController> }
             * Bind objectSymbol and allocated cell on domElementReuseCollector.get() call
             */
            allocatedCells: new Map(),
            pageHeight: undefined,
            /** set and use when nodes above viewport changes their sizings */
            scrollShift: undefined,
            /** @type {Map.<Symbol, { starts: number, ends: number, height: number }>} */
            positions: new Map(),
            /**
             * Holds the set of object symbols from current and next iterations of update.
             * Intended to be used by update functions.
             * @type {Set.<Symbol>}
             */
            mergedObjectSymbols: undefined,
            /** @type { { inViewport: Set.<Symbol>, inPreload: Set.<Symbol>, inParking: Set.<Symbol> } } */
            zoneCollusions: {
                inViewport: new Set(),
                inPreload: new Set(),
                inParking: new Set(),
            },
            /** @type { { viewport: { starts: number, ends: number }, preload: { starts: number, ends: number }, parking: { starts: number, ends: number } } } */
            boundaries: {
                viewport: {},
                /** an area which its height is 3 times of
                 * viewport (1 above, 1 below) */
                preload: {},
                /** an area which its height is 5 times of
                 * viewport (2 above, 2 below) */
                parking: {},
            },
            classifiedObjects: {
                toConstruct: new Set(),
                toAppear: new Set(),
                toUpdatePositionY: new Set(),
                toUpdatePositionX: new Set(), // indentation
                toUpdateFolding: new Set(),
                toUpdateExistance: new Set(),
                toDisappear: new Set(),
                toDestruct: new Set(),
            },
        };

        this.updateZoneBoundaries();
        this.calculateComponentPositions();
        this.classifyObjectsByCollidedZones();
        this.mergeObjectSymbolsWithPreviousIteration();
        this.calculateFocusShift();
        this.updateContainer();
        this.updateComponents();

        delete this.computedValues.current; // forget positions computed on previous call
        this.computedValues.current = this.computedValues.next;
    }

    build() {
        this.updateViewFromData();
    }
}
