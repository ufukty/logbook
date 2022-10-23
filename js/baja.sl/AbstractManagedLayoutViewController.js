import { Layout } from "./Layout/Layout.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { createElement, mergeMapKeys, setIntersect, setDifference } from "./utilities.js";
import { Area, Position } from "./Layout/Coordinates.js";
import { itemCellPairing } from "./ItemCellPairing.js";
import {
    AbstractManagedLayoutCellViewController,
    PRESENTATION_STATE_PLACEHOLDER,
} from "./AbstractManagedLayoutCellViewController.js";

export const TRIGGER_CONTENT_CHANGE = "TRIGGER_CONTENT_CHANGE";
export const TRIGGER_REPLACEMENT = "TRIGGER_REPLACEMENT";
export const TRIGGER_SCROLL_LISTENER = "TRIGGER_SCROLL_LISTENER";
export const TRIGGER_RESIZE_OBSERVER = "TRIGGER_RESIZE_OBSERVER";
export const TRIGGER_SCHEDULED_POSITION_TRANSLATE = "TRIGGER_SCHEDULED_POSITION_TRANSLATE";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} EnvironmentSymbol
 */

export class AbstractManagedLayoutViewController extends AbstractViewController {
    constructor() {
        super();

        this.dom = {
            ...this.dom,
            container: createElement("div", ["baja-sl-managed-layout-view-controller"]),
        };

        this.config = {
            /**
             * @type {HTMLElement}
             * That will be used to decide which potion of container is visible
             *   to user and cell reuse will be handled automatically at each
             *   scroll.
             */
            scrollElement: undefined,
            /**
             * @type {Layout}
             * Implementing class will use this for accessing up-to-date
             *   positions.
             */
            layout: undefined,
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
                preload: 1.1,
                parking: 1.2,
            },
            /**
             * focusedPosition = focusLevel * window.clientHeight
             */
            focusLevel: 0.5,
            updateMaxFrequency: 60,
        };

        /** @private */
        this.computedValues = {
            /** @type { Map.<Symbol, AbstractViewController> } */
            cellPositioners: new Map(),
            /** @type {} */
            timeoutToResizeCallback: undefined,
            /**
             * Scheduling unassignment gives cells to perform transition to
             *   their new place before unassign. When an item decided as
             *   "doesLeaveParking", it should start to transition and put
             *   into <transitioning>, a callback runs after transition also
             *   should put itemSymbol in <readyToUnassign>. Next update of VC,
             *   item will be unassigned from its cell and cell will be
             *   available for another item to move in.
             */
            unassignmentScheduling: {
                /** @type {Set.<Symbol>} */
                transitioning: new Set(),
                /** @type {Set.<Symbol>} */
                readyToUnassign: new Set(),
            },
            updateScheduling: {
                lastUpdateTime: undefined,
                ongoingUpdate: false,
                waitingForScheduledUpdate: undefined,
            },
            current: this._getTemplateForComputedValues(),
            next: this._getTemplateForComputedValues(),
        };

        /** @private */
        this.processAtNextUpdate = {
            /** @type {Area} */
            viewport: undefined,
        };

        this.dom.container.style.position = "relative";
    }

    /** @private */
    _getTemplateForComputedValues() {
        return {
            /** @type {string} */
            updateTrigger: undefined,
            /** @type {number} */
            pageHeight: undefined,
            /** set and use when nodes above viewport changes their sizings */
            scrollShift: undefined,
            /** @type {Map.<Symbol, Area>} */
            positions: new Map(),
            /**
             * Holds the set of item symbols from current and next iterations of update.
             * Intended to be used by update functions.
             * @type {Set.<Symbol>}
             */
            mergedItemSymbols: undefined,
            /** @type { { inFocus: Set.<Symbol>, inViewport: Set.<Symbol>, inPreload: Set.<Symbol>, inParking: Set.<Symbol> } } */
            zoneCollisions: {
                inFocus: new Set(),
                inViewport: new Set(),
                inPreload: new Set(),
                inParking: new Set(),
            },
            boundaries: {
                /** @type {Area} */
                focus: undefined,
                /** @type {Area} */
                view: undefined,
                /**
                 * @type {Area}
                 * an area which its height is 3 times of
                 * viewport (1 above, 1 below) */
                preload: undefined,
                /**
                 * @type {Area}
                 * an area which its height is 5 times of
                 * viewport (2 above, 2 below) */
                parking: undefined,
            },
            classifiedItems: {
                /** @type {Set.<Symbol>} */
                toAppear: new Set(),
                /** @type {Set.<Symbol>}
                 * When an item is in the placement but not in viewport */
                toPlaced: new Set(),
                /** @type {Set.<Symbol>} */
                toAssign: new Set(),
                /** @type {Set.<Symbol>} */
                toDisappear: new Set(),
                /** @type {Set.<Symbol>} */
                toPositionSet: new Set(),
                /** @type {Set.<Symbol>} */
                toPositionTranslate: new Set(),
                /** @type {Set.<Symbol>}
                 * When an item persisting in placement and changes its position
                 * with entering preload, but current position is not set yet. */
                toPositionTranslateFromCurrentPosition: new Set(),
                /** @type {Set.<Symbol>}
                 * When transition to another position required before unassigning
                 * item from its cell.
                 */
                toScheduledPositionTranslate: new Set(),
                /** @type {Set.<Symbol>} */
                toUnassign: new Set(),
                /** @type {Set.<Symbol>} */
                toCancelUnassign: new Set(),
            },
        };
    }

    /** @private */
    _updateZoneBoundaries() {
        if (this.processAtNextUpdate.viewport) {
            const viewport = this.processAtNextUpdate.viewport;
            this.processAtNextUpdate.viewport = undefined;

            this.computedValues.next.boundaries = {
                focus: new Area(viewport.x0, viewport.y0, viewport.x0, viewport.y0),
                view: viewport,
                preload: Object.assign(viewport, {}).scale(this.config.zoneOffsets.preload),
                parking: Object.assign(viewport, {}).scale(this.config.zoneOffsets.parking),
            };
        } else if (this.computedValues.current.boundaries.view) {
            this.computedValues.next.boundaries = this.computedValues.current.boundaries;
        } else {
            console.error(
                "AbstractManagedLayoutViewController is requested to get updated without viewport specified."
            );
        }
    }

    /** @private */
    _copyLayout() {
        this.computedValues.next.positions = this.config.layout.passedThroughPipeline.layout.positions;
    }

    /** @private */
    _classifyItemsByCollidedZones() {
        for (const [itemSymbol, item] of this.computedValues.next.positions) {
            var inFocus = false;
            var inViewport = false;
            var inPreload = false;
            var inParking = false;

            inParking = item.isCollidingWith(this.computedValues.next.boundaries.view);
            if (inParking) {
                inPreload = item.isCollidingWith(this.computedValues.next.boundaries.preload);
                if (inPreload) {
                    inViewport = item.isCollidingWith(this.computedValues.next.boundaries.view);
                    if (inViewport) {
                        inFocus = item.isCollidingWith(this.computedValues.next.boundaries.focus);
                    }
                }
            }

            if (inFocus) this.computedValues.next.zoneCollisions.inFocus.add(itemSymbol);
            if (inViewport) this.computedValues.next.zoneCollisions.inViewport.add(itemSymbol);
            if (inPreload) this.computedValues.next.zoneCollisions.inPreload.add(itemSymbol);
            if (inParking) this.computedValues.next.zoneCollisions.inParking.add(itemSymbol);
        }
    }

    /** @private */
    _mergeItemSymbolsWithPreviousIteration() {
        this.computedValues.next.mergedItemSymbols = mergeMapKeys(
            this.computedValues.current.positions,
            this.computedValues.next.positions
        );
    }

    /** @private */
    _updateContainerToTheContentBoundingBoxSize() {
        this.dom.container.style.width = `${this.config.layout.passedThroughPipeline.contentBoundingBoxSize.width}px`;
        this.dom.container.style.height = `${this.config.layout.passedThroughPipeline.contentBoundingBoxSize.height}px`;
    }

    /** @private */
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
                isPersistingInPlacement &&
                itemCellPairing.isItemAssignedToACell(itemSymbol) &&
                this.computedValues.current.zoneCollisions.inParking.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inParking.has(itemSymbol);

            const doesEnterPreload =
                !itemCellPairing.isItemAssignedToACell(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(itemSymbol);

            const doesLeavePreload =
                !itemCellPairing.isItemAssignedToACell(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inPreload.has(itemSymbol);

            const doesLeaveParking =
                itemCellPairing.isItemAssignedToACell(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inParking.has(itemSymbol);

            const isLeavingParking =
                !this.computedValues.unassignmentScheduling.transitioning.has(itemSymbol) &&
                itemCellPairing.isItemAssignedToACell(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inParking.has(itemSymbol);

            const isReadyToUnassign =
                this.computedValues.unassignmentScheduling.transitioning.has(itemSymbol) &&
                this.computedValues.unassignmentScheduling.readyToUnassign.has(itemSymbol);

            const isReturningToParking =
                this.computedValues.unassignmentScheduling.transitioning.has(itemSymbol) &&
                !this.computedValues.current.zoneCollisions.inParking.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inParking.has(itemSymbol);

            // if translation animation should be considered for item

            const isPositionChanged =
                isPersistingInPlacement &&
                this.computedValues.current.positions.get(itemSymbol).starts !==
                    this.computedValues.next.positions.get(itemSymbol).starts;

            // if item collides viewport

            const isPersistingInViewport =
                isPersistingInPlacement &&
                this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            const doesEnterViewport =
                !this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            const doesLeaveViewport =
                this.computedValues.current.zoneCollisions.inViewport.has(itemSymbol) &&
                !this.computedValues.next.zoneCollisions.inViewport.has(itemSymbol);

            // specifying necessary updates for the item

            /**
             * reason of "leave"/"enter" zones might be one of:
             *   - user-scroll
             *   - item's movement
             */

            // console.groupCollapsed(symbolizer.desymbolize(itemSymbol));
            // console.log("doesEnterPlacement:", doesEnterPlacement);
            // console.log("doesLeavePlacement:", doesLeavePlacement);
            // console.log("isPersistingInPlacement:", isPersistingInPlacement);
            // console.log("isPersistingInPreload:", isPersistingInPreload);
            // console.log("doesEnterPreload:", doesEnterPreload);
            // console.log("doesLeavePreload:", doesLeavePreload);
            // console.log("doesLeaveParking:", doesLeaveParking);
            // console.log("isLeavingParking:", isLeavingParking);
            // console.log("isReadyToUnassign:", isReadyToUnassign);
            // console.log("isReturningToParking:", isReturningToParking);
            // console.log("isPositionChanged:", isPositionChanged);
            // console.log("isPersistingInViewport:", isPersistingInViewport);
            // console.log("doesEnterViewport:", doesEnterViewport);
            // console.log("doesLeaveViewport:", doesLeaveViewport);
            // console.groupEnd(symbolizer.desymbolize(itemSymbol));

            if (isPersistingInPlacement) {
                if (isPersistingInPreload) {
                    if (isPositionChanged) {
                        this.computedValues.next.classifiedItems.toPositionTranslate.add(itemSymbol);
                    } else {
                        this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                    }
                    if (doesEnterViewport) {
                        this.computedValues.next.classifiedItems.toAppear.add(itemSymbol);
                    } else if (doesLeaveViewport) {
                        this.computedValues.next.classifiedItems.toDisappear.add(itemSymbol);
                    }
                } else if (doesEnterPreload) {
                    this.computedValues.next.classifiedItems.toAssign.add(itemSymbol);
                    if (isPositionChanged) {
                        this.computedValues.next.classifiedItems.toPositionTranslateFromCurrentPosition.add(itemSymbol);
                    } else {
                        this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                    }
                } else if (isLeavingParking) {
                    if (isPositionChanged) {
                        this.computedValues.next.classifiedItems.toScheduledPositionTranslate.add(itemSymbol);
                    } else {
                        this.computedValues.next.classifiedItems.toUnassign.add(itemSymbol);
                    }
                } else if (isReturningToParking) {
                    this.computedValues.next.classifiedItems.toCancelUnassign.add(itemSymbol);
                    if (isPositionChanged) {
                        this.computedValues.next.classifiedItems.toPositionTranslate.add(itemSymbol);
                    } else {
                        this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                    }
                } else if (isReadyToUnassign) {
                    this.computedValues.next.classifiedItems.toUnassign.add(itemSymbol);
                }
            } else if (doesEnterPlacement) {
                if (doesEnterPreload) {
                    this.computedValues.next.classifiedItems.toAssign.add(itemSymbol);
                    this.computedValues.next.classifiedItems.toPositionSet.add(itemSymbol);
                }
                if (doesEnterViewport) {
                    this.computedValues.next.classifiedItems.toAppear.add(itemSymbol);
                } else {
                    this.computedValues.next.classifiedItems.toPlaced.add(itemSymbol);
                }
            } else if (doesLeavePlacement) {
                if (doesLeaveParking) {
                    this.computedValues.next.classifiedItems.toUnassign.add(itemSymbol);
                    if (doesLeaveViewport) {
                        this.computedValues.next.classifiedItems.toDisappear.add(itemSymbol);
                    }
                }
            }
        }
    }

    /** @private */
    _debugPrintClassifiedItems() {
        Object.keys(this.computedValues.next.classifiedItems).forEach((cls) => {
            if (this.computedValues.next.classifiedItems[cls].size > 0) {
                console.log(cls, this.computedValues.next.classifiedItems[cls]);
            }
        });
    }

    /**
     * @private
     * @param {string} trigger
     */
    _updateCells() {
        /**
         * execution order of updates:
         *   - assign
         *   - position set
         *   - position transition
         *   - position translate and unassign
         *   - notify placement without appear
         *   - appear
         *   - unassign
         *   - disappear
         */

        const classes = this.computedValues.next.classifiedItems;

        // assign
        for (const itemSymbol of classes.toAssign) {
            const envSymbol = this.config.layout.environmentSymbol;
            const managedLayoutCellViewController = itemCellPairing.assign(itemSymbol, envSymbol);
            managedLayoutCellViewController.leveledPresentation(PRESENTATION_STATE_PLACEHOLDER);
            // this.populateCellForItem(managedLayoutCellViewController, itemSymbol);

            // const computedStyle = getComputedStyle(cellPositioner.cell.dom.container);
            // const computedHeight = parseFloat(computedStyle.getPropertyValue("height"));
            // this.computedValues.lastRecordedCellHeightOfItem.set(itemSymbol, computedHeight);
        }

        // position set
        for (const itemSymbol of classes.toPositionSet) {
            const managedLayoutCellViewController = itemCellPairing.getAssignedCellForItem(itemSymbol);
            const newPosition = this.computedValues.next.positions.get(itemSymbol);
            managedLayoutCellViewController.setPosition(new Position(newPosition.x0, newPosition.y0), false, false);
        }

        // position transition
        for (const itemSymbol of classes.toPositionTranslate) {
            const cellPositioner = itemCellPairing.getAssignedCellForItem(itemSymbol);
            var newPosition = this.computedValues.next.positions.get(itemSymbol).starts;
            cellPositioner.setPosition(newPosition, true, false);
        }

        // position transition from current position
        for (const itemSymbol of classes.toPositionTranslateFromCurrentPosition) {
            const cellPositioner = itemCellPairing.getAssignedCellForItem(itemSymbol);
            const currentPosition = this.computedValues.current.positions.get(itemSymbol).starts;
            const newPosition = this.computedValues.next.positions.get(itemSymbol).starts;
            cellPositioner.setPosition(currentPosition, false, true);
            cellPositioner.setPosition(newPosition, true, false);
        }

        // position transition with scheduling
        for (const itemSymbol of classes.toScheduledPositionTranslate) {
            const cellPositioner = itemCellPairing.getAssignedCellForItem(itemSymbol);
            const newPosition = this.computedValues.next.positions.get(itemSymbol).starts;
            this.computedValues.unassignmentScheduling.transitioning.add(itemSymbol);
            cellPositioner.setPosition(newPosition, true, false, () => {
                this.computedValues.unassignmentScheduling.readyToUnassign.add(itemSymbol);
                this.updateView(TRIGGER_SCHEDULED_POSITION_TRANSLATE);
            });
        }

        // placement without appear
        for (const itemSymbol of classes.toPlaced) {
            this.cellPlacedWithoutAppear(itemSymbol);
        }

        // appear
        for (const itemSymbol of classes.toAppear) {
            this.cellAppears(itemSymbol);
        }

        // cancel unassign
        for (const itemSymbol of classes.toCancelUnassign) {
            this.computedValues.unassignmentScheduling.transitioning.delete(itemSymbol);
            this.computedValues.unassignmentScheduling.readyToUnassign.delete(itemSymbol);
        }

        // unassign
        for (const itemSymbol of classes.toUnassign) {
            itemCellPairing.unassign(itemSymbol);
            this.computedValues.cellPositioners.delete(itemSymbol);
            this.computedValues.unassignmentScheduling.transitioning.delete(itemSymbol);
            this.computedValues.unassignmentScheduling.readyToUnassign.delete(itemSymbol);
        }

        // disappear
        for (const itemSymbol of classes.toDisappear) {
            this.cellDisappears(itemSymbol);
        }
    }

    /** @private
     * Note: This function re-schedules an update to future if update
     * condition is not satisfied currently. */
    _isUpdateNeeded(trigger) {
        const now = Date.now();

        // const doesParentElementHasHeight =
        //     (this.computedValues.current.boundaries.view &&
        //         this.computedValues.current.boundaries.view.height() !== 0) ||
        //     (this.processAtNextUpdate.viewport && this.processAtNextUpdate.viewport.height() !== 0);
        // if (!doesParentElementHasHeight) {
        //     // console.log("re-scheduling updateView due to scrollElement height is not being set");
        //     requestAnimationFrame(() => {
        //         this.updateView(trigger);
        //     });
        //     return false;
        // }

        if (this.computedValues.updateScheduling.lastUpdateTime) {
            const timePassedSinceLastUpdate = now - this.computedValues.updateScheduling.lastUpdateTime;
            const periodRequiredMS = 1000 / this.config.updateMaxFrequency;
            const remainingToNextUpdate = periodRequiredMS - timePassedSinceLastUpdate;

            if (this.computedValues.updateScheduling.ongoingUpdate || timePassedSinceLastUpdate <= periodRequiredMS) {
                if (!this.computedValues.updateScheduling.waitingForScheduledUpdate) {
                    this.computedValues.updateScheduling.waitingForScheduledUpdate = true;
                    setTimeout(() => {
                        this.updateView(trigger);
                    }, remainingToNextUpdate + 1);
                }
                // console.log("update is rejected");
                return false;
            }

            if (this.computedValues.updateScheduling.waitingForScheduledUpdate) {
                this.computedValues.updateScheduling.waitingForScheduledUpdate = undefined;
            }
        }
        // console.log("updating");

        this.computedValues.updateScheduling.lastUpdateTime = now;
        return true;
    }

    updateView(trigger) {
        if (!this._isUpdateNeeded(trigger)) return;
        this.computedValues.updateScheduling.ongoingUpdate = true;

        this.computedValues.next = this._getTemplateForComputedValues();
        this.computedValues.next.updateTrigger = trigger;

        console.log(trigger);

        this._updateZoneBoundaries();
        this._copyLayout();
        this._classifyItemsByCollidedZones();
        this._mergeItemSymbolsWithPreviousIteration();
        this._updateContainerToTheContentBoundingBoxSize();
        this._classifyItemsByUpdateTypes();
        // this._debugPrintClassifiedItems();

        requestAnimationFrame(() => {
            this._updateCells();
            delete this.computedValues.current;
            this.computedValues.current = this.computedValues.next;
            this.computedValues.updateScheduling.ongoingUpdate = undefined;
        });
    }

    /** @param {Area} newArea */
    setViewport(newArea) {
        this.processAtNextUpdate.viewport = newArea;
    }

    /**
     * @param {CellTypeSymbol} cellTypeSymbol
     * @param {EnvironmentSymbol} environmentSymbol
     * @param {function():AbstractManagedLayoutCellViewController} constructorFunction
     */
    registerCellViewControllerConstructor(cellTypeSymbol, constructorFunction) {
        const environmentSymbol = this.config.layout.environmentSymbol;
        itemCellPairing.registerCellViewControllerConstructor(cellTypeSymbol, environmentSymbol, () => {
            const cell = constructorFunction();
            this.dom.container.appendChild(cell.dom.managedLayoutPositioner);
            return cell;
        });
    }

    /**
     * @abstract
     * This function will be called for each cell that placed in placement but
     *   could not appear in viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} itemSymbol
     */
    cellPlacedWithoutAppear(itemSymbol) {
        // console.error("abstract function is called directly");
    }

    /**
     * @abstract
     * This function will be called for each cell that enters into the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} itemSymbol
     */
    cellAppears(itemSymbol) {
        // console.error("abstract function is called directly");
    }

    /**
     * @abstract
     * This function will be called for each cell that exits from the viewport.
     * Implementer can use this method to perform UI updates on rest of the cell.
     * @param {Symbol} itemSymbol
     */
    cellDisappears(itemSymbol) {
        // console.error("abstract function is called directly");
    }

    /**
     * Calling this function will trigger getCellForItem() method implemented by
     * subclass if those items are in preload area
     * @param {Set.<Symbol>} updatedItems - Symbols of items
     */
    requestContentUpdateForItemsIfNecessary(updatedItems) {
        const assignedItems = new Set(this.computedValues.cellPositioners.keys());
        const updatedAssignedItems = setIntersect(assignedItems, updatedItems);
        const updatedUnassignedItems = setDifference(updatedItems, updatedAssignedItems);

        for (const itemSymbol of updatedAssignedItems) {
            this.updateCellForUpdatedItem(itemSymbol, this.computedValues.cellPositioners.get(itemSymbol));
        }
        if (updatedAssignedItems.size > 0) {
            this.updateView(TRIGGER_CONTENT_CHANGE);
        }
        for (const itemSymbol of updatedUnassignedItems) {
            this.cellUpdateIsSkippedForUpdatedItem(itemSymbol);
        }
    }

    /**
     * @abstract
     * @param {managedLayoutCellViewController} managedLayoutCellViewController
     * @param {ItemSymbol} itemSymbol
     * Populate content of managedLayoutCellViewController according to itemSymbol
     */
    populateCellForItem(managedLayoutCellViewController, itemSymbol) {
        console.error("abstract method is called directly.");
    }
}
