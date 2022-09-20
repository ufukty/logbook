import { adoption, domCollector, createElement, symbolizer, mergeMapKeys, checkCollision } from "./utilities.js";
import { AbstractTableCellPositioner } from "./AbstractTableCellPositioner.js";
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
             * Smaller numbers means less cells created at DOM and less memory
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
             * Each `Symbol` represents an `itemSymbol`
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
             * Maps `itemIdSymbol` to correct reuse identifiers.
             * Information will be used for requesting and
             *   sending cells to `domElementReuseCollector`.
             * Note that: Related constructors for each id
             *   given as key to this map, should've
             *   registered to `domElementReuseCollector` already.
             */
            itemReuseIdentifiers: new Map(),
            /**
             * @type {HTMLElement}
             * This element is the one that will be listened for scroll elements
             *   on, and will be used to device which items to show in view.
             */
            scrollElement: this.container,
        };

        this.computedValues = {
            /** @type { Map.<Symbol, AbstractTableCellPositioner> } */
            cellPositioners: new Map(),
            /** @type {Map.<Symbol, number>} */
            lastRecordedCellHeightOfItem: new Map(),
            /** @type {} */
            timeoutToResizeCallback: undefined,
            current: this._getTemplateForComputedValues(),
            next: this._getTemplateForComputedValues(),
        };

        this.container.addEventListener("scroll", (e) => {
            console.log(e);
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
            const itemId = cellPositioner.dataset["itemId"];
            const itemSymbol = symbolizer.symbolize(itemId);

            const isStillAllocated = this.computedValues.cellPositioners.has(itemSymbol);
            if (!isStillAllocated) return;

            const isSameHeight =
                this.computedValues.lastRecordedCellHeightOfItem.has(itemSymbol) &&
                this.computedValues.lastRecordedCellHeightOfItem.get(itemSymbol) === height;
            if (isSameHeight) return;

            this.computedValues.lastRecordedCellHeightOfItem.set(itemSymbol, height);
            // changedItems.push(itemId);
            ignoreChanges = false;

            // const isInSight = this.computedValues.current.zoneCollusions.inViewport.has(itemSymbol);
            // if (isInSight) {
            //     ignoreChanges = false;
            // }
        });
        // console.log("resizeObserver detected height change for:", changedItems);
        if (ignoreChanges) return;
        // to avoid infinite resize loops
        if (this.computedValues.timeoutToResizeCallback === undefined) {
            this.computedValues.timeoutToResizeCallback = setTimeout(() => {
                requestAnimationFrame(() => {
                    this.computedValues.timeoutToResizeCallback = undefined;
                    this.updateView(TRIGGER_RESIZE_OBSERVER);
                });
            }, 500);
        }
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
     * specified itemSymbol.
     * @returns {AbstractTableCellPositioner}
     */
    getCellForItem(itemSymbol) {
        this._error("abstract function is called directly");
    }

    /** @param {Symbol} itemSymbol */
    getCellKindForItem(itemSymbol) {
        this._error("Abstract method has called directly.");
    }

    /**
     * Implementation of this method should check if content of cell
     * needs to get updated, then update it.
     * @param { Symbol } itemSymbol
     * @param {AbstractTableCellPositioner} cellContainer
     */
    updateCellIfNecessary(itemSymbol, cellContainer) {
        this._error("abstract function is called directly");
    }

    /** @returns {number} */
    getAverageHeightForAnItem() {
        return 20;
    }

    /**
     * Default height is important to estimate overall height of
     * the page and make the scrollbar much more useful.
     * @param {number} itemSymbol
     * @returns {number}
     */
    getDefaultHeightOfItem(itemSymbol) {
        this._error("abstract function is called directly");
    }

    _updateZoneBoundaries() {
        const scrollAreaHeight = this.config.scrollElement.clientHeight;
        const scrollTop = this.config.scrollElement.scrollTop;

        const preloadZoneOffset = Math.floor(this.config.zoneOffsets.preload * scrollAreaHeight);
        const parkingZoneOffset = Math.floor(this.config.zoneOffsets.parking * scrollAreaHeight);

        this.computedValues.next.boundaries = {
            focus: {
                starts: scrollTop + this.config.focusPoint * scrollAreaHeight,
                ends: scrollTop + this.config.focusPoint * scrollAreaHeight,
            },
            viewport: {
                starts: scrollTop,
                ends: scrollTop + scrollAreaHeight,
            },
            /** an area which its height is 3 times of viewport (1 above, 1 below) */
            preload: {
                starts: scrollTop - preloadZoneOffset,
                ends: scrollTop + scrollAreaHeight + preloadZoneOffset,
            },
            /** an area which its height is 5 times of viewport (2 above, 2 below) */
            parking: {
                starts: scrollTop - parkingZoneOffset,
                ends: scrollTop + scrollAreaHeight + parkingZoneOffset,
            },
        };
    }

    _filterPlacement() {
        this.config.placement.symbols.forEach((symbol) => {
            if (!this.config.placement.ignore.has(symbol)) this.computedValues.next.filteredPlacement.push(symbol);
        });
    }

    _calculateItemPositions() {
        let lastPosition = this.config.margins.pageContent.before,
            lastCellKind = undefined,
            lastItemIndex = this.computedValues.next.filteredPlacement.length - 1;

        for (const [itemIndex, itemSymbol] of this.computedValues.next.filteredPlacement.entries()) {
            // apply "before/between/after" margins to the lastPosition

            const currentCellKind = this.getCellKindForItem(itemSymbol);
            const marginsToApply = {
                beforePageContent: itemIndex === 0,
                afterPageContent: itemIndex === lastItemIndex,
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

            const cellHeight = this.computedValues.lastRecordedCellHeightOfItem.has(itemSymbol)
                ? this.computedValues.lastRecordedCellHeightOfItem.get(itemSymbol)
                : this.getDefaultHeightOfItem(itemSymbol);

            // save item positions
            this.computedValues.next.positions.set(itemSymbol, {
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
    _classifyItemsByCollidedZones() {
        for (const [itemSymbol, itemPos] of this.computedValues.next.positions) {
            var inFocus = false,
                inViewport = false,
                inPreload = false,
                inParking = false;
            inParking = checkCollision(
                itemPos.starts,
                itemPos.ends,
                this.computedValues.next.boundaries.parking.starts,
                this.computedValues.next.boundaries.parking.ends
            );
            if (inParking)
                inPreload = checkCollision(
                    itemPos.starts,
                    itemPos.ends,
                    this.computedValues.next.boundaries.preload.starts,
                    this.computedValues.next.boundaries.preload.ends
                );
            if (inPreload)
                inViewport = checkCollision(
                    itemPos.starts,
                    itemPos.ends,
                    this.computedValues.next.boundaries.viewport.starts,
                    this.computedValues.next.boundaries.viewport.ends
                );
            if (inViewport)
                inFocus = checkCollision(
                    itemPos.starts,
                    itemPos.ends,
                    this.computedValues.next.boundaries.focus.starts,
                    this.computedValues.next.boundaries.focus.ends
                );
            if (inFocus) this.computedValues.next.zoneCollisions.inFocus.add(itemSymbol);
            if (inViewport) this.computedValues.next.zoneCollisions.inViewport.add(itemSymbol);
            if (inPreload) this.computedValues.next.zoneCollisions.inPreload.add(itemSymbol);
            if (inParking) this.computedValues.next.zoneCollisions.inParking.add(itemSymbol);
        }
    }

    _mergeItemSymbolsWithPreviousIteration() {
        this.computedValues.next.mergedItemSymbols = mergeMapKeys(
            this.computedValues.current.positions,
            this.computedValues.next.positions
        );
    }

    _calculateFocusShift() {
        let totalScrollShift = 0;
        for (const itemSymbol of this.computedValues.next.mergedItemSymbols) {
            // calculate scroll shift; if there is any change in height of item
            if (
                this.computedValues.current.positions.has(itemSymbol) &&
                this.computedValues.next.positions.has(itemSymbol)
            ) {
                if (
                    this.computedValues.current.positions.get(itemSymbol).height !==
                    this.computedValues.next.positions.get(itemSymbol).height
                ) {
                    if (
                        this.computedValues.next.positions.get(itemSymbol).ends <
                        this.computedValues.next.boundaries.viewport.starts
                    ) {
                        totalScrollShift +=
                            this.computedValues.next.positions.get(itemSymbol).height -
                            this.computedValues.current.positions.get(itemSymbol).height;
                    }
                }
            }
        }
        this.computedValues.scrollShift = totalScrollShift;
    }

    _counterScrollToJustifyHeightChange() {
        this.container;
    }

    _classifyItemsByUpdateTypes() {
        for (const itemSymbol of this.computedValues.next.mergedItemSymbols) {
            /**
             * Two reasons one item to be classified as "toAssign":
             *   1. When an item enters to placement array
             *   2. When an item persisting in placement enters to viewport
             *      from outside of viewport, e.g. after user scrolls
             */

            /**
             * 1. items persisting in placement but changes their positions
             *    1.1. in-view  -> in-view  : items in view translates some position also in view
             *    1.2. out-view -> in-view  : items enter view
             *    1.3. in-view  -> out-view : items leave view
             *    1.4. out-view -> out-view : items that enter view from one edge and leave from other edge
             * 2. items enter placement
             * 3. items leave placement
             */

            // to help decide assign/unassign "with or without" translation animation

            const doesEnterPlacement =
                !this.computedValues.current.positions.has(itemSymbol) &&
                this.computedValues.next.positions.has(itemSymbol);

            const doesLeavePlacement =
                this.computedValues.current.positions.has(itemSymbol) &&
                !this.computedValues.next.positions.has(itemSymbol);

            const isPersistingInPlacement =
                this.computedValues.current.positions.has(itemSymbol) &&
                this.computedValues.next.positions.has(itemSymbol);

            // to help decide if item will be assigned/unassigned from cell

            const isPersistingInPreload =
                this.computedValues.cellPositioners.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(itemSymbol);

            const doesEnterPreload =
                !this.computedValues.cellPositioners.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(itemSymbol);

            const doesLeaveParking =
                this.computedValues.cellPositioners.has(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inParking.has(itemSymbol);

            const didLeaveParking =
                this.computedValues.cellPositioners.has(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inParking.has(itemSymbol); // FIXME:

            // if translation animation should be considered for item

            const isPositionChanged =
                isPersistingInPlacement &&
                this.computedValues.current.positions.get(itemSymbol).starts !==
                    this.computedValues.next.positions.get(itemSymbol).starts;

            // if item collides viewport

            const isPersistingInViewport =
                this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            const doesEnterViewport =
                !this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            const doesLeaveViewport =
                this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            // specifying necessary updates for the item

            if (isPersistingInPlacement) {
                if (isPersistingInPreload) {
                    if (isPositionChanged) {
                        this.computedValues.next.classifiedItems.toPositionTranslate.add(itemSymbol);
                        if (doesEnterViewport) {
                            this.computedValues.next.classifiedItems.toAppear.add(itemSymbol);
                        } else if (doesLeaveViewport) {
                            this.computedValues.next.classifiedItems.toDisappear.add(itemSymbol);
                        }
                    }
                } else if (doesEnterPreload) {
                    this.computedValues.next.classifiedItems.toAssign.add(itemSymbol);
                    this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                } else if (doesLeaveParking) {
                    this.computedValues.next.classifiedItems.toPositionTranslate.add(itemSymbol); // TODO: schedule/mark for unassign
                } else if (didLeaveParking) {
                    // this.computedValues.next.classifiedItems.toUnassign.add(itemSymbol);
                }
            } else if (doesEnterPlacement) {
                if (doesEnterPreload) {
                    this.computedValues.next.classifiedItems.toAssign.add(itemSymbol);
                    this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                }
                if (doesEnterViewport) {
                    this.computedValues.next.classifiedItems.toAppear.add(itemSymbol);
                }
            } else if (doesLeavePlacement) {
                // this.computedValues.next.classifiedItems.toUnassign.add(itemSymbol);
                if (doesLeaveViewport) {
                    this.computedValues.next.classifiedItems.toDisappear.add(itemSymbol);
                }
            }
        }
    }

    /**
     * @param {string} trigger
     */
    _updateComponents(trigger) {
        /**
         * execution order of updates:
         *   - assign
         *   - position set
         *   - position transition
         *   - position translate and unassign
         *   - appear
         *   - unassign
         *   - disappear
         */

        const classes = this.computedValues.next.classifiedItems;

        // assign
        for (const itemSymbol of classes.toAssign) {
            const cellPositioner = this.getCellForItem(itemSymbol);
            this.computedValues.cellPositioners.set(itemSymbol, cellPositioner);
            cellPositioner.itemSymbol = itemSymbol;
            cellPositioner.container.dataset["itemId"] = symbolizer.desymbolize(itemSymbol);
        }

        // position set
        for (const itemSymbol of classes.toPositionSet) {
            const cellPositioner = this.computedValues.cellPositioners.get(itemSymbol);
            const newPosition = this.computedValues.next.positions.get(itemSymbol).starts;
            cellPositioner.setPosition(newPosition, false);
        }

        // position transition
        for (const itemSymbol of classes.toPositionTranslate) {
            const cellPositioner = this.computedValues.cellPositioners.get(itemSymbol);
            const newPosition = this.computedValues.next.positions.get(itemSymbol).starts;
            cellPositioner.setPosition(newPosition, true);
        }

        // appear
        for (const itemSymbol of classes.toAppear) {
            this.cellAppears(itemSymbol);
        }

        // unassign
        for (const itemSymbol of classes.toUnassign) {
            const cellPositioner = this.computedValues.cellPositioners.get(itemSymbol);
            const reuseIdentifier = cellPositioner.reuseIdentifier;
            domCollector.free(reuseIdentifier, cellPositioner);
            this.computedValues.cellPositioners.delete(itemSymbol);
        }

        // disappear
        for (const itemSymbol of classes.toDisappear) {
            this.cellDisappears(itemSymbol);
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
             * Holds the set of item symbols from current and next iterations of update.
             * Intended to be used by update functions.
             * @type {Set.<Symbol>}
             */
            mergedItemSymbols: undefined,
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

            classifiedItems: {
                // toAssign: new Set(),
                // toAppear: new Set(),
                // toUpdatePositionY: new Set(),
                // toUpdatePositionX: new Set(), // indentation
                // toUpdateFolding: new Set(),
                // toUpdateExistance: new Set(),
                // toDisappear: new Set(),
                // toUnassign: new Set(),
                /** @type {Set.<Symbol>} */
                toAppear: new Set(),
                /** @type {Set.<Symbol>} */
                toAssign: new Set(),
                /** @type {Set.<Symbol>} */
                toDisappear: new Set(),
                /** @type {Set.<Symbol>} */
                toPositionSet: new Set(),
                /** @type {Set.<Symbol>} */
                toPositionTranslate: new Set(),
                /** @type {Set.<Symbol>} */
                toUnassign: new Set(),
                /** @type {Set.<Symbol>} */
                toUnassign: new Set(),
            },
        };
    }

    _debugUpdateStart(trigger) {
        console.groupCollapsed(`AbstractTableViewController.updateView(${trigger})`);
    }

    _debugClassifiedComponents() {
        const classes = [
            "toAppear",
            "toAssign",
            "toDisappear",
            "toPositionSet",
            "toPositionTranslate",
            "toUnassign",
            "toUnassign",
        ];
        classes.forEach((cls) => {
            if (this.computedValues.next.classifiedItems[cls].size > 0) {
                console.log(cls, this.computedValues.next.classifiedItems[cls]);
            }
        });
    }

    _debugUpdateEnd(trigger) {
        console.groupEnd(`AbstractTableViewController.updateView(${trigger})`);
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
        this._debugUpdateStart(trigger);

        this._updateZoneBoundaries();
        this._filterPlacement();
        this._calculateItemPositions();
        this._classifyItemsByCollidedZones();

        this._mergeItemSymbolsWithPreviousIteration();
        this._calculateFocusShift();
        this._updateContainer();
        this._classifyItemsByUpdateTypes();
        this._debugClassifiedComponents();

        console.group("updateComponents");
        this._updateComponents(trigger);
        console.groupEnd("updateComponents");

        this._debugUpdateEnd(trigger);

        delete this.computedValues.current;
        this.computedValues.current = this.computedValues.next;
    }

    /**
     * This function will be called for each cell that enters into the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} itemSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellAppears(itemSymbol, cellPositioner) {
        this._error("abstract function is called directly");
    }

    /**
     * This function will be called for each cell that exits from the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} itemSymbol
     * @param {AbstractTableCellViewController} cellPositioner
     */
    cellDisappears(itemSymbol, cellPositioner) {
        this._error("abstract function is called directly");
    }

    /**
     * Calling this function will trigger getCellForItem() method implemented by
     * subclass if those items are in preload area
     * @param {Set.<Symbol>} itemSymbols - Symbols of items
     */
    requestContentUpdateForItemIfNecessary(itemSymbols) {
        const intersect = new Set();
        for (const itemSymbolAllocated of this.computedValues.cellPositioners.keys()) {
            if (itemSymbols.has(itemSymbolAllocated)) intersect.add(itemSymbolAllocated);
        }
        for (const itemSymbol of intersect) {
            this.updateCellIfNecessary(itemSymbol, this.computedValues.cellPositioners.get(itemSymbol));
        }
        if (intersect.size > 0) {
            this.updateView();
        }
    }
}
