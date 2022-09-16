import { adoption, domCollector, createElement, symbolizer, mergeMapKeys, checkCollision } from "./utilities.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";
import { BasicTableCellPositioner } from "./BasicTableCellPositioner.js";

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

        Object.assign(this.config, {
            debug: true,
            margins: {
                pageContent: {
                    before: 10,
                    after: 10,
                },
            },
            /**
             * Number values for zones represents how many times window height
             *   will be extended from up and down to find top and bottom edges of
             *   each zone. Suggestion: preload < parking
             *
             * Items that came inside of preload area will be assigned to a cell.
             * Items that went outside of parking area will be unassigned from
             *   a cell.
             *
             * Smaller numbers means less objects created at DOM and less memory
             *   usage but also the user will notice absense of items when scrolling
             *   fast.
             */
            zoneOffsets: {
                preload: 0.4,
                parking: 0.5,
            },
            /**
             * The ordering of sections and rows in them.
             * Each `Symbol` represents an `objectSymbol`
             * (either a `sectionID` or `rowID`).
             */
            placement: {
                /**
                 * Incomplete list of placement data.
                 * @type {Array.<Symbol>}
                 */
                objectIds: [],
                /**
                 * States what is the actual index of items[0]
                 * @type {number}
                 */
                offset: undefined,
                /**
                 * Total number of items in the document. That value is used
                 *   for estimation of full height of cell scroller for both
                 *   chronological and hierarchical view.
                 * @type {number}
                 */
                totalNumberOfItems: undefined,
            },
            /**
             * @type { Map.<Symbol, Symbol> }
             * Maps `objectIdSymbol` to correct reuse identifiers.
             *   Information will be used for requesting and
             *   sending cells to `domElementReuseCollector`.
             * Note that: Related constructors for each id
             *   given as key to this map, should've
             *   registered to `domElementReuseCollector` already.
             */
            objectReuseIdentifiers: new Map(),
        });

        this.computedValues = {
            /** @type { Map.<Symbol, AbstractTableCellPositioner> } */
            objectToCellContainers: new Map(),
            /** @type {Map.<Symbol, number>} */
            lastRecordedObjectHeight: new Map(),
            current: this._getTemplateForComputedValues(),
            next: this._getTemplateForComputedValues(),
        };

        document.addEventListener("scroll", this.updateView.bind(this));

        this.resizeObserver = new ResizeObserver((entries) => {
            var nothingIsChanged = true;
            entries.forEach((entry) => {
                const height = Math.ceil(entry.contentRect.height);
                const cellContainer_container = entry.target;
                const objectId = cellContainer_container.dataset["objectId"];
                const objectSymbol = symbolizer.symbolize(objectId);
                if (
                    !(
                        this.computedValues.lastRecordedObjectHeight.has(objectSymbol) &&
                        this.computedValues.lastRecordedObjectHeight.get(objectSymbol) === height
                    )
                ) {
                    nothingIsChanged = false;
                    this.computedValues.lastRecordedObjectHeight.set(objectSymbol, height);
                    this._debug("height is changed:", objectId, height);
                } else {
                    this._debug("height is NOT changed:", objectId, height);
                }
            });
            if (!nothingIsChanged) {
                this.updateView();
            }
        });
    }

    /**
     * Embeds the user-supplied cell constructor with a function
     * that creates a custom positioner view controller and wraps
     * the cell returned by user-supplied cell constructor with it.
     */
    registerCellIdentifier(cellIdentifier, cellConstructor) {
        domCollector.registerItemIdentifier(cellIdentifier, () => {
            const userProvidedCell = cellConstructor();
            const cellContainer = new BasicTableCellPositioner();
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
     * @returns {AbstractTableCellPositioner}
     */
    requestReusableCellContainer(cellIdentifier) {
        return domCollector.get(cellIdentifier);
    }

    /**
     * User should implement this method.
     * Request an empty cell from .getFreeCell()
     * with previously registered cellIdentifier
     * Then populate content accordingly to
     * specified objectSymbol.
     * @returns {AbstractTableCellPositioner}
     */
    getCellForObject(objectSymbol) {
        this._error("abstract function is called directly");
    }

    /** @param {Symbol} objectSymbol */
    getCellKindForObject(objectSymbol) {
        this._error("Abstract method has called directly.");
    }

    /**
     * Implementation of this method should check if content of cell
     * needs to get updated, then update it.
     * @param { Symbol } objectSymbol
     * @param {AbstractTableCellPositioner} cellContainer
     */
    updateCellIfNecessary(objectSymbol, cellContainer) {
        this._error("abstract function is called directly");
    }

    /** @returns {number} */
    getAverageHeightForAnObject() {
        return 20;
    }

    /**
     * Default height is important to estimate overall height of
     * the page and make the scrollbar much more useful.
     * @param {number} objectSymbol
     * @returns {number}
     */
    getDefaultHeightOfObject(objectSymbol) {
        this._error("abstract function is called directly");
    }

    _updateZoneBoundaries() {
        const preloadZoneOffset = Math.floor(this.config.zoneOffsets.preload * window.innerHeight);
        const parkingZoneOffset = Math.floor(this.config.zoneOffsets.parking * window.innerHeight);

        this.computedValues.next.boundaries = {
            viewport: {
                starts: window.scrollY,
                ends: window.scrollY + window.innerHeight,
            },
            /** an area which its height is 3 times of viewport (1 above, 1 below) */
            preload: {
                starts: window.scrollY - preloadZoneOffset,
                ends: window.scrollY + window.innerHeight + preloadZoneOffset,
            },
            /** an area which its height is 5 times of viewport (2 above, 2 below) */
            parking: {
                starts: window.scrollY - parkingZoneOffset,
                ends: window.scrollY + window.innerHeight + parkingZoneOffset,
            },
        };
    }

    // TODO: don't mind tasks that their parents are folded.
    _calculateComponentPositions() {
        let lastPosition = 0;
        let lastCellKind = undefined;
        lastPosition += this.config.margins.pageContent.before;

        let objectIds = this.config.placement.objectIds;
        let lastObjectIndex = objectIds.length - 1;
        for (const [objectIndex, objectId] of objectIds.entries()) {
            // apply "before/between/after" margins to the lastPosition

            const currentCellKind = this.getCellKindForObject(objectId);
            const marginsToApply = {
                beforePageContent: objectIndex === 0,
                afterPageContent: objectIndex === lastObjectIndex,
                betweenSameKind: lastCellKind && currentCellKind === lastCellKind,
                afterMarginForPreviousKind: lastCellKind && currentCellKind !== lastCellKind,
                beforeMarginForCurrentKind: lastCellKind && currentCellKind !== lastCellKind,
            };
            if (marginsToApply.beforePageContent) {
                const margin = this.config.margins.pageContent.before;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.beforeMarginForCurrentKind) {
                const margin = this.config.margins[currentCellKind].before;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.afterMarginForPreviousKind) {
                const margin = this.config.margins[lastCellKind].after;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.betweenSameKind) {
                const margin = this.config.margins[currentCellKind].between;
                lastPosition += margin ? margin : 0;
            }

            const cellHeight = this.computedValues.lastRecordedObjectHeight.has(objectId)
                ? this.computedValues.lastRecordedObjectHeight.get(objectId)
                : this.getDefaultHeightOfObject(objectId);

            // save object positions
            this.computedValues.next.positions.set(objectId, {
                starts: lastPosition,
                ends: lastPosition + cellHeight,
                height: cellHeight,
            });

            lastPosition += cellHeight;

            if (marginsToApply.afterPageContent) {
                const margin = this.config.margins.pageContent.after;
                lastPosition += margin ? margin : 0;
            }
            lastCellKind = currentCellKind;
        }

        this.computedValues.next.pageHeight = lastPosition;
    }

    /** @access private */
    _classifyObjectsByCollidedZones() {
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

    _mergeObjectSymbolsWithPreviousIteration() {
        this.computedValues.next.mergedObjectSymbols = mergeMapKeys(
            this.computedValues.current.positions,
            this.computedValues.next.positions
        );
    }

    _calculateFocusShift() {
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

    _classifyComponentsByUpdateTypes() {
        for (const objectSymbol of this.computedValues.next.mergedObjectSymbols) {
            // "to first-construct"
            if (
                !this.computedValues.objectToCellContainers.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inPreload.has(objectSymbol) &&
                !this.computedValues.current.positions.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toFirstConstruct.add(objectSymbol);
            }

            // "to repeating-construct"
            if (
                !this.computedValues.objectToCellContainers.has(objectSymbol) &&
                this.computedValues.next.zoneCollusions.inPreload.has(objectSymbol) &&
                this.computedValues.current.positions.has(objectSymbol)
            ) {
                this.computedValues.next.classifiedObjects.toRepeatingConstruct.add(objectSymbol);
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

    _updateComponents() {
        // "to destruct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDestruct) {
            // console.log(objectSymbol);
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            const reuseIdentifier = cellContainer.reuseIdentifier;
            domCollector.free(reuseIdentifier, cellContainer);
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

        // "to first-construct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toFirstConstruct) {
            const cellContainer = this.getCellForObject(objectSymbol);
            this.computedValues.objectToCellContainers.set(objectSymbol, cellContainer);
            cellContainer.objectSymbol = objectSymbol;
            cellContainer.container.dataset["objectId"] = symbolizer.desymbolize(objectSymbol);

            // this case is when item emerges from non-existance
            // TODO: animate emergence?
            // cellContainer.setPositionY(this.computedValues.next.positions.get(objectSymbol).starts + 10, false);
            cellContainer.setPositionY(this.computedValues.next.positions.get(objectSymbol).starts, false);

            this.computedValues.lastRecordedObjectHeight.set(objectSymbol, cellContainer.container.clientHeight);
        }

        // "to repeating-construct"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toRepeatingConstruct) {
            const cellContainer = this.getCellForObject(objectSymbol);
            this.computedValues.objectToCellContainers.set(objectSymbol, cellContainer);
            cellContainer.objectSymbol = objectSymbol;
            cellContainer.container.dataset["objectId"] = symbolizer.desymbolize(objectSymbol);

            // this case is when item emerges from non-existance
            cellContainer.setPositionY(this.computedValues.current.positions.get(objectSymbol).starts, false);
            cellContainer.setPositionY(this.computedValues.next.positions.get(objectSymbol).starts, true);

            this.computedValues.lastRecordedObjectHeight.set(objectSymbol, cellContainer.container.clientHeight);
        }

        // "to update position Y"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toUpdatePositionY) {
            const cellContainer = this.computedValues.objectToCellContainers.get(objectSymbol);
            const newPosition = this.computedValues.next.positions.get(objectSymbol).starts;
            cellContainer.setPositionY(newPosition, true);
        }

        // "to update position X"
    }

    _updateContainer() {
        this.container.style.height = `${this.computedValues.next.pageHeight}px`;
    }

    _getTemplateForComputedValues() {
        return {
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
                toFirstConstruct: new Set(),
                toRepeatingConstruct: new Set(),
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

    _debugUpdatedComponents() {
        if (!this.config.debug) return;

        console.log("AbstractTableViewController._debugUpdatedComponents");
        const classes = [
            "toFirstConstruct",
            "toRepeatingConstruct",
            "toAppear",
            "toUpdatePositionY",
            "toUpdatePositionX",
            "toUpdateFolding",
            "toUpdateExistance",
            "toDisappear",
            "toDestruct",
        ];
        classes.forEach((cls) => {
            if (this.computedValues.next.classifiedObjects[cls].size > 0) {
                console.log(cls, this.computedValues.next.classifiedObjects[cls]);
            }
        });
    }

    updateView() {
        this.computedValues.next = this._getTemplateForComputedValues();

        this._updateZoneBoundaries();
        this._calculateComponentPositions();
        this._classifyObjectsByCollidedZones();

        this._mergeObjectSymbolsWithPreviousIteration();
        this._calculateFocusShift();
        this._updateContainer();
        this._classifyComponentsByUpdateTypes();
        this._updateComponents();

        this._debugUpdatedComponents();

        delete this.computedValues.current;
        this.computedValues.current = this.computedValues.next;
    }

    /**
     * This function will be called for each cell that enters into the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellAppears(objectSymbol, cellPositioner) {
        this._error("abstract function is called directly");
    }

    /**
     * This function will be called for each cell that exits from the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} objectSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellDisappears(objectSymbol, cellPositioner) {
        this._error("abstract function is called directly");
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
