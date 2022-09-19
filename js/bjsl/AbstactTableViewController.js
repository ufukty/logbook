import { adoption, domCollector, createElement, symbolizer, mergeMapKeys, checkCollision } from "./utilities.js";
import {
    AbstractTableCellPositioner,
    POSITION_ANIMATE,
    POSITION_INSTANT,
    POSITION_REDIRECT_IF_PLAYING,
} from "./AbstractTableCellPositioner.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";
import { BasicTableCellPositioner } from "./BasicTableCellPositioner.js";

export const TRIGGER_REPLACEMENT = "TRIGGER_REPLACEMENT";
export const TRIGGER_SCROLL_LISTENER = "TRIGGER_SCROLL_LISTENER";
export const TRIGGER_RESIZE_OBSERVER = "TRIGGER_RESIZE_OBSERVER";

export class AbstractTableViewController {
    constructor() {
        this.container = createElement("div", ["abstract-cell-scroller-view"]);
        this.contentArea = createElement("div", ["abstract-cell-scroller-view-content-area"]);
        this.anchorPosition = createElement("div", ["abstract-cell-scroller-view-anchor-position"]);
        // prettier-ignore
        adoption(this.container, 
            adoption(this.contentArea, 
                this.anchorPosition
        ));

        this.config = {
            ...this.config,
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
             * focusPoint * window.clientHeight
             */
            focusPoint: 0.5,
            /**
             * The ordering of sections and rows in them.
             * Each `Symbol` represents an `objectSymbol`
             * (either a `sectionID` or `rowID`).
             */
            placement: {
                /**
                 * @type {Array.<Symbol>}
                 * Complete list of placement data, height-ignored items should
                 *   also be in this array.
                 */
                symbols: [],
                /**
                 * @type {Set.<Symbol>}
                 * Those items which their symbols pushed to this array, shall
                 *   be excluded from the calculation of item positions.
                 * Example use case is folding/collapsing a portion of table
                 *   while performing some transition on items before excluding
                 *   from placement (in order to perform animation before
                 *   unassign those items from their cells. )
                 */
                ignore: new Set(),
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
             * Information will be used for requesting and
             *   sending cells to `domElementReuseCollector`.
             * Note that: Related constructors for each id
             *   given as key to this map, should've
             *   registered to `domElementReuseCollector` already.
             */
            objectReuseIdentifiers: new Map(),
        };

        this.computedValues = {
            /** @type { Map.<Symbol, AbstractTableCellPositioner> } */
            cellPositioners: new Map(),
            /** @type {Map.<Symbol, number>} */
            lastRecordedObjectHeight: new Map(),
            current: this._getTemplateForComputedValues(),
            next: this._getTemplateForComputedValues(),
        };

        document.addEventListener("scroll", (e) => {
            // console.log(e.target.scrollingElement.scrollTop);
            this.updateView(TRIGGER_SCROLL_LISTENER);
        });

        this.resizeObserver = new ResizeObserver(this._resizeObserverNotificationHandler.bind(this));
    }

    /**
     * @param {Array.<ResizeObserverEntry>} entries
     */
    _resizeObserverNotificationHandler(entries) {
        // console.log(entries);
        var ignoreChanges = true;
        var changedItems = [];
        entries.forEach((entry) => {
            const height = entry.contentRect.height;
            const cellPositioner = entry.target;
            const objectId = cellPositioner.dataset["objectId"];
            const objectSymbol = symbolizer.symbolize(objectId);

            const isStillAllocated = this.computedValues.cellPositioners.has(objectSymbol);
            if (!isStillAllocated) return;

            const isSameHeight =
                this.computedValues.lastRecordedObjectHeight.has(objectSymbol) &&
                this.computedValues.lastRecordedObjectHeight.get(objectSymbol) === height;
            if (isSameHeight) return;

            this.computedValues.lastRecordedObjectHeight.set(objectSymbol, height);
            // changedItems.push(objectId);
            ignoreChanges = false;

            // const isInSight = this.computedValues.current.zoneCollusions.inViewport.has(objectSymbol);
            // if (isInSight) {
            //     ignoreChanges = false;
            // }
        });
        // console.log("resizeObserver detected height change for:", changedItems);
        if (ignoreChanges) return;
        // to avoid infinite resize loops
        requestAnimationFrame(() => {
            this.updateView(TRIGGER_RESIZE_OBSERVER);
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
            focus: {
                starts: window.scrollY + this.config.focusPoint * window.innerHeight,
                ends: window.scrollY + this.config.focusPoint * window.innerHeight,
            },
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

    _filterPlacement() {
        this.config.placement.symbols.forEach((symbol) => {
            if (!this.config.placement.ignore.has(symbol)) this.computedValues.next.filteredPlacement.push(symbol);
        });
    }

    _calculateComponentPositions() {
        let lastPosition = this.config.margins.pageContent.before,
            lastCellKind = undefined,
            lastObjectIndex = this.computedValues.next.filteredPlacement.length - 1;

        for (const [objectIndex, objectSymbol] of this.computedValues.next.filteredPlacement.entries()) {
            // apply "before/between/after" margins to the lastPosition

            const currentCellKind = this.getCellKindForObject(objectSymbol);
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

            const cellHeight = this.computedValues.lastRecordedObjectHeight.has(objectSymbol)
                ? this.computedValues.lastRecordedObjectHeight.get(objectSymbol)
                : this.getDefaultHeightOfObject(objectSymbol);

            // save object positions
            this.computedValues.next.positions.set(objectSymbol, {
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
        for (const [objectSymbol, objectPos] of this.computedValues.next.positions) {
            var inFocus = false,
                inViewport = false,
                inPreload = false,
                inParking = false;
            inParking = checkCollision(
                objectPos.starts,
                objectPos.ends,
                this.computedValues.next.boundaries.parking.starts,
                this.computedValues.next.boundaries.parking.ends
            );
            if (inParking)
                inPreload = checkCollision(
                    objectPos.starts,
                    objectPos.ends,
                    this.computedValues.next.boundaries.preload.starts,
                    this.computedValues.next.boundaries.preload.ends
                );
            if (inPreload)
                inViewport = checkCollision(
                    objectPos.starts,
                    objectPos.ends,
                    this.computedValues.next.boundaries.viewport.starts,
                    this.computedValues.next.boundaries.viewport.ends
                );
            if (inViewport)
                inFocus = checkCollision(
                    objectPos.starts,
                    objectPos.ends,
                    this.computedValues.next.boundaries.focus.starts,
                    this.computedValues.next.boundaries.focus.ends
                );
            if (inFocus) this.computedValues.next.zoneCollisions.inFocus.add(objectSymbol);
            if (inViewport) this.computedValues.next.zoneCollisions.inViewport.add(objectSymbol);
            if (inPreload) this.computedValues.next.zoneCollisions.inPreload.add(objectSymbol);
            if (inParking) this.computedValues.next.zoneCollisions.inParking.add(objectSymbol);
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

    _counterScrollToJustifyHeightChange() {
        this.container;
    }

    _classifyComponentsByUpdateTypes() {
        for (const objectSymbol of this.computedValues.next.mergedObjectSymbols) {
            /**
             * Two reasons one item to be classified as "toAssign":
             *   1. When an item enters to placement array
             *   2. When an item persisting in placement enters to viewport
             *      from outside of viewport, e.g. after user scrolls
             */

            /**
             * 1. objects persisting in placement but changes their positions
             *    1.1. in-view  -> in-view  : objects in view translates some position also in view
             *    1.2. out-view -> in-view  : objects enter view
             *    1.3. in-view  -> out-view : objects leave view
             *    1.4. out-view -> out-view : objects that enter view from one edge and leave from other edge
             * 2. objects enter placement
             * 3. objects leave placement
             */

            // to help decide assign/unassign "with or without" translation animation

            const doesEnterPlacement =
                !this.computedValues.current.positions.has(objectSymbol) &&
                this.computedValues.next.positions.has(objectSymbol);

            const doesLeavePlacement =
                this.computedValues.current.positions.has(objectSymbol) &&
                !this.computedValues.next.positions.has(objectSymbol);

            const isPersistingInPlacement =
                this.computedValues.current.positions.has(objectSymbol) &&
                this.computedValues.next.positions.has(objectSymbol);

            // to help decide if item will be assigned/unassigned from cell

            const isPersistingInPreload =
                this.computedValues.cellPositioners.has(objectSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(objectSymbol);

            const doesEnterPreload =
                !this.computedValues.cellPositioners.has(objectSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(objectSymbol);

            const doesLeaveParking =
                this.computedValues.cellPositioners.has(objectSymbol) &&
                !this.computedValues.next.zoneCollisions.inParking.has(objectSymbol); // FIXME:

            const didLeaveParking = this.computedValues.current.classifiedObjects; // FIXME:

            // if translation animation should be considered for item

            const isPositionChanged =
                isPersistingInPlacement &&
                this.computedValues.current.positions.get(objectSymbol).starts !==
                    this.computedValues.next.positions.get(objectSymbol).starts;

            // if item collides viewport

            const isPersistingInViewport =
                this.computedValues.current.zoneCollisions.inViewport.has(objectSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(objectSymbol);

            const doesEnterViewport =
                !this.computedValues.current.zoneCollisions.inViewport.has(objectSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(objectSymbol);

            const doesLeaveViewport =
                this.computedValues.current.zoneCollisions.inViewport.has(objectSymbol) &&
                !this.computedValues.next.zoneCollisions.inViewport.has(objectSymbol);

            // specifying necessary updates for the item

            /**
             * execution order of updates:
             *   - assign
             *   - position set
             *   - position transition
             *   - appear
             *   - unassign
             *   - disappear
             */

            if (isPersistingInPlacement) {
                if (isPersistingInPreload) {
                    if (isPositionChanged) {
                        // position translate
                        if (doesEnterViewport) {
                            // appear(item)
                        } else if (doesLeaveViewport) {
                            // disappear(item)
                        }
                    }
                } else if (doesEnterPreload) {
                    // to assign
                } else if (doesLeaveParking) {
                    // FIXME:
                    // position translate
                    // schedule unassign
                } else if (didLeaveParking) {
                    // FIXME:
                    // to unassign
                }
            } else if (doesEnterPlacement) {
                if (doesEnterPreload) {
                    // assign
                    // set position
                }
                if (doesEnterViewport) {
                    // appear(item)
                }
            } else if (doesLeavePlacement) {
                // unassign
                if (doesLeaveViewport) {
                    // disappear(item)
                }
            }
        }
    }

    /**
     * @param {string} trigger
     */
    _updateComponents(trigger) {
        // "to unassign"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toUnassign) {
            const cellPositioner = this.computedValues.cellPositioners.get(objectSymbol);
            const reuseIdentifier = cellPositioner.reuseIdentifier;
            domCollector.free(reuseIdentifier, cellPositioner);
            this.computedValues.cellPositioners.delete(objectSymbol);
        }

        // "to disappear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toDisappear) {
            const cellPositioner = this.computedValues.cellPositioners.get(objectSymbol);
            this.cellDisappears(objectSymbol, cellPositioner);
        }

        // existance change
        // folding change

        // "to appear"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toAppear) {
            const cellPositioner = this.computedValues.cellPositioners.get(objectSymbol);
            this.cellAppears(objectSymbol, cellPositioner);
        }

        // "to assign"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toAssign) {
            const cellPositioner = this.getCellForObject(objectSymbol);
            this.computedValues.cellPositioners.set(objectSymbol, cellPositioner);
            cellPositioner.objectSymbol = objectSymbol;
            cellPositioner.container.dataset["objectId"] = symbolizer.desymbolize(objectSymbol);
            // cellPositioner.setPositionY(this.computedValues.next.positions.get(objectSymbol).starts, POSITION_INSTANT);
        }

        // "to update position Y"
        for (const objectSymbol of this.computedValues.next.classifiedObjects.toUpdatePositionY) {
            const cellPositioner = this.computedValues.cellPositioners.get(objectSymbol);
            const currentPosition = this.computedValues.current.positions.get(objectSymbol).starts;
            const nextPosition = this.computedValues.next.positions.get(objectSymbol).starts;

            if (currentPosition !== nextPosition) {
                var animation;
                switch (trigger) {
                    case TRIGGER_RESIZE_OBSERVER:
                        animation = POSITION_REDIRECT_IF_PLAYING;
                        break;
                    case TRIGGER_REPLACEMENT:
                        animation = POSITION_ANIMATE;
                        break;
                    case TRIGGER_SCROLL_LISTENER:
                        // console.error("why there is location change?");
                        animation = POSITION_INSTANT;
                        break;
                    default:
                        animation = POSITION_INSTANT;
                }
                cellPositioner.setPositionY(nextPosition, animation);
            }
        }
    }

    _updateContainer() {
        this.contentArea.style.height = `${this.computedValues.next.pageHeight}px`;
    }

    _getTemplateForComputedValues() {
        return {
            /** @type {Array.<Symbol>} */
            filteredPlacement: [],
            /** @type {string} */
            updateTrigger: undefined,
            /** @type {number} */
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
            zoneCollisions: {
                inFocus: new Set(),
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
                toAssign: new Set(),
                toAppear: new Set(),
                toUpdatePositionY: new Set(),
                toUpdatePositionX: new Set(), // indentation
                toUpdateFolding: new Set(),
                toUpdateExistance: new Set(),
                toDisappear: new Set(),
                toUnassign: new Set(),
            },
        };
    }

    _debugUpdateStart() {
        console.group(`AbstractTableViewController.updateView(${this.computedValues.next.updateTrigger})`);
    }

    _debugClassifiedComponents() {
        const classes = [
            "toAssign",
            "toAppear",
            "toUpdatePositionY",
            "toUpdatePositionX",
            "toUpdateFolding",
            "toUpdateExistance",
            "toDisappear",
            "toUnassign",
        ];
        classes.forEach((cls) => {
            if (this.computedValues.next.classifiedObjects[cls].size > 0) {
                console.log(cls, this.computedValues.next.classifiedObjects[cls]);
            }
        });
    }

    _debugUpdateEnd() {
        console.groupEnd(`AbstractTableViewController.updateView(${this.computedValues.next.updateTrigger})`);
    }

    _isUpdateNeededForScroll(nextTrigger) {
        // const last = this.computedValues.scrollUpdates.lastUpdatedScrollPosition
        if (nextTrigger !== TRIGGER_SCROLL_LISTENER) return true;
        const lastTrigger = this.computedValues.current.updateTrigger;
        if (lastTrigger !== TRIGGER_SCROLL_LISTENER) return true;

        // TODO:
        // if the scrolled distance is not greater than half of the distance
        // between "preload" and "view" zones ignore scroll (return false)

        return true;
    }

    updateView(trigger) {
        if (!this._isUpdateNeededForScroll(trigger)) return;

        this.computedValues.next = this._getTemplateForComputedValues();
        this.computedValues.next.updateTrigger = trigger;
        this._debugUpdateStart();

        this._updateZoneBoundaries();
        this._filterPlacement();
        this._calculateComponentPositions();
        this._classifyObjectsByCollidedZones();

        this._mergeObjectSymbolsWithPreviousIteration();
        this._calculateFocusShift();
        this._updateContainer();
        this._classifyComponentsByUpdateTypes();
        this._debugClassifiedComponents();

        console.group("updateComponents");
        this._updateComponents(trigger);
        console.groupEnd("updateComponents");

        this._debugUpdateEnd();

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
    requestContentUpdateForObjectsIfNecessary(objectSymbols) {
        const intersect = new Set();
        for (const objectSymbolAllocated of this.computedValues.cellPositioners.keys()) {
            if (objectSymbols.has(objectSymbolAllocated)) intersect.add(objectSymbolAllocated);
        }
        for (const objectSymbol of intersect) {
            this.updateCellIfNecessary(objectSymbol, this.computedValues.cellPositioners.get(objectSymbol));
        }
        if (intersect.size > 0) {
            this.updateView();
        }
    }
}
