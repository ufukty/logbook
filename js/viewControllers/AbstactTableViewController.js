import { adoption, domElementReuseCollector, createElement, toggleAnimationWithClass, pSymbol } from "../utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { TableViewStructuredDataMedium } from "../dataSource.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";

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

/**
 * @param {Set} leftSet
 * @param {Set} rightSet
 * Returns a set represents intersection of two input sets. */
function setIntersect(leftSet, rightSet) {
    const intersection = new Set();
    for (const k of leftSet) {
        if (rightSet.has(k)) intersection.add(k);
    }
    return intersection;
}

export class AbstractTableViewCellContainerViewController extends AbstractViewController {
    constructor() {
        super();
        /** @type {HTMLElement} */
        this.container = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type { AbstractTableCellViewController }
         */
        this.cell = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type { Symbol }
         */
        this.reuseIdentifier = undefined;

        /** Filled by CellScrollerViewController. Don't modify that.
         * @type {Symbol}
         */
        this.objectSymbol = undefined;
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionX(newPosition, withAnimation) {
        console.error("abstract function is called directly");
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation) {
        console.error("abstract function is called directly");
    }
}

export class BasicTableCellContainerViewController extends AbstractTableViewCellContainerViewController {
    constructor() {
        super();
        this.container = createElement("div", ["abstract-cell-scroller-view-cell-positioner"]);
        this.cell = undefined; // should be assigned by callee
    }

    prepareForFree() {
        this.container.style.visibility = "hidden";
        if (this.animation) {
            this.animation.cancel();
            delete this.animation;
        }
        this.cell.prepareForFree();
        this.container.style.top = `0px`;
    }

    prepareForUse() {
        this.container.style.visibility = "visible";
        this.cell.prepareForUse();
    }

    /**
     * @param {number} newPosition
     * @param {boolean} withAnimation
     */
    setPositionY(newPosition, withAnimation = false) {
        if (withAnimation) {
            const objectSymbolAtAnimationStart = this.objectSymbol;

            const currentPos = this.container.style.top.match(/\d+/);
            const deltaPos = newPosition - currentPos;

            // prettier-ignore
            const keyframes = [
                { transform: "translateY(0px)" }, 
                { transform: `translateY(${deltaPos}px)` }
            ]
            const config = {
                duration: 200,
                iterations: 1,
                fill: "none",
                // easing: "cubic-bezier(0.5, 1.2, 0.8, 1.0)",
            };
            this.animation = this.container.animate(keyframes, config);
            this.animation.finished
                .then(() => {
                    if (objectSymbolAtAnimationStart === this.objectSymbol) {
                        this.container.style.top = `${newPosition}px`;
                    }
                })
                .catch(() => {
                    // console.log("animation is aborted");
                });
        } else {
            this.container.style.top = `${newPosition}px`;
        }
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
        this.container = createElement("div", ["abstract-cell-scroller-view"]);
        this.contentArea = createElement("div", ["abstract-cell-scroller-view-content-area"]);
        this.anchorPosition = createElement("div", ["abstract-cell-scroller-view-anchor-position"]);
        // prettier-ignore
        adoption(this.container, 
            adoption(this.contentArea, 
                this.anchorPosition
        ));

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
            /** @type { Map.<Symbol, AbstractTableViewCellContainerViewController> } */
            objectToCellContainers: new Map(),
            /** @type {Map.<Symbol, number>} */
            lastRecordedObjectHeight: new Map(),
            current: {
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

        document.addEventListener("scroll", this.updateView.bind(this));

        this.resizeObserver = new ResizeObserver((entries) => {
            entries.forEach((entry) => {
                const height = Math.ceil(entry.contentRect.height);

                const cellContainer_container = entry.target;
                const objectId = cellContainer_container.dataset["objectId"];
                const objectSymbol = pSymbol.get(objectId);

                this.computedValues.lastRecordedObjectHeight.set(objectSymbol, height);
            });
            this.updateView();
        });
    }

    /**
     * Embeds the user-supplied cell constructor with a function
     * that creates a custom positioner view controller and wraps
     * the cell returned by user-supplied cell constructor with it.
     */
    registerCellIdentifier(cellIdentifier, cellConstructor) {
        domElementReuseCollector.registerItemIdentifier(cellIdentifier, () => {
            const userProvidedCell = cellConstructor();
            const cellContainer = new BasicTableCellContainerViewController();
            cellContainer.cell = userProvidedCell;
            cellContainer.reuseIdentifier = cellIdentifier;
            // prettier-ignore
            adoption(this.anchorPosition,
                adoption(cellContainer.container,
                    userProvidedCell.container
            ));
            this.resizeObserver.observe(cellContainer.container);
            return cellContainer;
        });
    }

    /**
     * When user request a cell to populate with data, this method
     * only sends the nested user-supplied custom cell, instead
     * the positioner cell that wraps it from the constructor
     * registered by .registerCellConstructor().
     * @returns {AbstractTableViewCellContainerViewController}
     */
    requestReusableCellContainer(cellIdentifier) {
        return domElementReuseCollector.get(cellIdentifier);
    }

    /**
     * User should implement this method.
     * Request an empty cell from .getFreeCell()
     * with previously registered cellIdentifier
     * Then populate content accordingly to
     * specified objectSymbol.
     * @returns {AbstractTableViewCellContainerViewController}
     */
    getCellForObject(objectSymbol) {
        console.error("abstract function is called directly");
    }

    /**
     * Implementation of this method should check if content of cell
     * needs to get updated, then update it.
     * @param { Symbol } objectSymbol
     * @param {AbstractTableViewCellContainerViewController} cellContainer
     */
    updateCellIfNecessary(objectSymbol, cellContainer) {
        console.error("abstract function is called directly");
    }

    /**
     * Default height is important to estimate overall height of
     * the page and make the scrollbar much more useful.
     * @param { number } objectSymbol
     * @returns { number }
     */
    getDefaultHeightOfObject(objectSymbol) {
        console.error("abstract function is called directly");
    }

    updateZoneBoundaries() {
        const preloadZoneOffset = Math.floor(0.5 * window.innerHeight);
        const parkingZoneOffset = Math.floor(0.75 * window.innerHeight);

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

            const headerHeight = this.computedValues.lastRecordedObjectHeight.has(sectionID)
                ? this.computedValues.lastRecordedObjectHeight.get(sectionID)
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

                const itemHeight = this.computedValues.lastRecordedObjectHeight.has(rowID)
                    ? this.computedValues.lastRecordedObjectHeight.get(rowID)
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
        this.computedValues.next.mergedObjectSymbols = mergeMapKeys(
            this.computedValues.current.positions,
            this.computedValues.next.positions
        );
    }

    calculateFocusShift() {
        let totalScrollShift = 0;
        for (const objectSymbol of this.computedValues.next.mergedObjectSymbols) {
            // calculate scroll shift; if there is any change in height of object
            if (
                this.computedValues.current.positions.has(objectSymbol) &&
                this.computedValues.next.positions.has(objectSymbol)
            ) {
                if (
                    this.computedValues.current.positions.get(objectSymbol).height !==
                    this.computedValues.next.positions.get(objectSymbol).height
                ) {
                    if (
                        this.computedValues.next.positions.get(objectSymbol).ends <
                        this.computedValues.next.boundaries.viewport.starts
                    ) {
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
        for (const objectSymbol of this.computedValues.next.mergedObjectSymbols) {
            // "to construct"
            if (
                !this.computedValues.objectToCellContainers.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inPreload.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toConstruct.add(objectSymbol);
            }

            // "to appear"
            if (
                !this.computedValues.current.zoneCollusions.inViewport.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inViewport.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toAppear.add(objectSymbol);
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
                this.computedValues.next.classifiedObjects.toDisappear.add(objectSymbol);
            }

            // to destruct
            if (
                this.computedValues.objectToCellContainers.has(objectSymbol) &&
                !this.computedValues.next.zoneCollusions.inParking.has(objectSymbol)
            ) {
                // console.log("to destruct");
                this.computedValues.next.classifiedObjects.toDestruct.add(objectSymbol);
            }

            // for objects already allocated a cell and put in the page
            if (
                this.computedValues.objectToCellContainers.has(objectSymbol) &&
                !this.computedValues.next.classifiedObjects.toDestruct.has(objectSymbol)
            ) {
                // position change
                if (
                    !this.computedValues.current.positions.has(objectSymbol) ||
                    this.computedValues.current.positions.get(objectSymbol).starts !==
                        this.computedValues.next.positions.get(objectSymbol).starts
                ) {
                    this.computedValues.next.classifiedObjects.toUpdatePositionY.add(objectSymbol);
                }
            }
        }
    }

    updateComponents() {
        // "to destruct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDestruct) {
            // console.log(objectSymbol);
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            const reuseIdentifier = cellContainer.reuseIdentifier;
            domElementReuseCollector.free(reuseIdentifier, cellContainer);
            this.computedValues.objectToCellContainers.delete(objectSymbol);
        }

        // "to disappear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDisappear) {
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            this.cellDisappears(objectSymbol, cellContainer);
        }

        // existance change
        // folding change

        // "to appear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toAppear) {
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            this.cellAppears(objectSymbol, cellContainer);
        }

        // "to construct"
        this.computedValues.next.needsRelayout = false;
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toConstruct) {
            const cellContainer = this.getCellForObject(objectSymbol);
            this.computedValues.objectToCellContainers.set(objectSymbol, cellContainer);
            cellContainer.objectSymbol = objectSymbol;
            cellContainer.container.dataset["objectId"] = pSymbol.reverse(objectSymbol);

            let objectInitializationPositionY, objectInitializationPositionX;
            if (this.computedValues.current.positions.has(objectSymbol)) {
                objectInitializationPositionY = this.computedValues.current.positions.get(objectSymbol).starts;
                objectInitializationPositionX = this.computedValues.current.positions.get(objectSymbol).x;
            } else if (this.computedValues.next.positions.has(objectSymbol)) {
                objectInitializationPositionY = this.computedValues.next.positions.get(objectSymbol).starts;
                objectInitializationPositionX = this.computedValues.next.positions.get(objectSymbol).x;
            }
            cellContainer.setPositionY(objectInitializationPositionY, false);
            cellContainer.setPositionX(objectInitializationPositionX, false);
            this.computedValues.lastRecordedObjectHeight.set(objectSymbol, cellContainer.container.clientHeight);
            this.computedValues.next.needsRelayout = true;
        }

        // "to update position Y"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toUpdatePositionY) {
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            const newPosition = this.computedValues.next.positions.get(objectSymbol).starts;
            cellContainer.setPositionY(newPosition, true);
        }
        
        // "to update position X"
    }

    updateContainer() {
        this.container.style.height = `${this.computedValues.pageHeight}px`;
    }

    prepareComputedValuesForTheUpdate() {
        this.computedValues.next = {
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
    }

    updateView(trigger) {
        // console.log("update [start]");
        this.prepareComputedValuesForTheUpdate();

        this.updateZoneBoundaries();
        this.calculateComponentPositions();
        this.classifyObjectsByCollidedZones();

        this.mergeObjectSymbolsWithPreviousIteration();
        this.calculateFocusShift();
        this.updateContainer();
        this.classifyComponentsByUpdateTypes();
        this.updateComponents();

        delete this.computedValues.current; // forget positions computed on previous call
        this.computedValues.current = this.computedValues.next;

        // console.log("update [end]");
        if (this.computedValues.next.needsRelayout) this.updateView();
    }

    /**
     * This function will be called for each cell that enters into the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellAppears(objectSymbol, cellPositioner) {
        console.error("abstract function is called directly");
    }

    /**
     * This function will be called for each cell that exits from the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellDisappears(objectSymbol, cellPositioner) {
        console.error("abstract function is called directly");
    }

    /**
     * Calling this function will trigger getCellForObject() method implemented by subclass if those objects are in preload area
     * @param {Set.<Symbol>} symbolsOfObjectsToUpdate - Symbols of objects
     */
    requestContentUpdateForObjectsIfNecessary(symbolsOfObjectsToUpdate) {
        const intersect = new Set();
        for (const objectSymbolAllocated of this.computedValues.objectToCellContainers.keys()) {
            if (symbolsOfObjectsToUpdate.has(objectSymbolAllocated)) intersect.add(objectSymbolAllocated);
        }
        for (const objectSymbol of intersect) {
            this.updateCellIfNecessary(objectSymbol, this.computedValues.objectToCellContainers.get(objectSymbol));
        }
        if (intersect.size > 0) {
            this.updateView();
        }
    }
}
