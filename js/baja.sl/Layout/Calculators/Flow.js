import { iota, symbolizer } from "../../utilities.js";
import { AbstractLayoutCalculator, AbstractLayoutMutator } from "../AbstractLayoutPipe.js";
import { Area, Size, Spacing } from "../Coordinates.js";
import { itemMeasurer as itemMeasurer } from "../../ItemMeasurer.js";
import { itemCellPairing } from "../../ItemCellPairing.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} ViewControllerSymbol
 */

/** Meant to be used for `Flow.config.direction` */
export const HORIZONTAL = symbolizer.symbolize("HORIZONTAL");
/** Meant to be used for `Flow.config.direction` */
export const VERTICAL = symbolizer.symbolize("VERTICAL");

export class Flow extends AbstractLayoutCalculator {
    /** @param {direction} */
    constructor(direction = VERTICAL) {
        super();
        this.config = {
            ...this.config,
            /** @type {Map.<CellTypeSymbol, Spacing>} */
            spacing: new Map(),
            /** Either `VERTICAL` or `HORIZONTAL` */
            direction: direction,
        };
    }

    perform() {
        var lastPosition = 0;
        var lastCellKind = undefined;
        var crossAxisMaxPosition = 0;
        const direction = this.config.direction;
        const totalNumberOfItems = this.config.totalNumberOfItems ?? this.config.placement.length;
        const offsetOfFirstItem = this.config.offset ?? 0;
        const remainingUnplacedItems = totalNumberOfItems - offsetOfFirstItem - this.config.placement.length;

        if (direction === VERTICAL) {
            const averageHeight = itemMeasurer.getAverageSize(this.controlledByEnvironment.environmentSymbol).height;
            const beforePlacementHeight = averageHeight * offsetOfFirstItem;
            lastPosition += beforePlacementHeight;
        } else if (direction === HORIZONTAL) {
            const averageWidth = itemMeasurer.getAverageSize(this.controlledByEnvironment.environmentSymbol).width;
            const beforePlacementWidth = averageWidth * offsetOfFirstItem;
            lastPosition += beforePlacementWidth;
        }

        for (const [itemIndex, itemSymbol] of this.config.placement.entries()) {
            const currentCellKind = itemCellPairing.getCellTypeForItem(
                itemSymbol,
                this.controlledByEnvironment.environmentSymbol
            );
            const marginsToApply = {
                // beforePageContent: itemIndex === 0,
                // afterPageContent: itemIndex === lastItemIndex,
                betweenSameKind: lastCellKind && currentCellKind === lastCellKind,
                afterMarginForPreviousKind: lastCellKind && currentCellKind !== lastCellKind,
                beforeMarginForCurrentKind: lastCellKind && currentCellKind !== lastCellKind,
            };
            // if (marginsToApply.beforePageContent) {
            //     const margin = this.config.spacing.container.before;
            //     lastPosition += margin ?? 0;
            // }
            if (marginsToApply.beforeMarginForCurrentKind) {
                const margin = this.config.spacing.get(currentCellKind).before;
                lastPosition += margin ?? 0;
            }
            if (marginsToApply.afterMarginForPreviousKind) {
                const margin = this.config.spacing.get(lastCellKind).after;
                lastPosition += margin ?? 0;
            }
            if (marginsToApply.betweenSameKind) {
                const margin = this.config.spacing.get(currentCellKind).between;
                lastPosition += margin ?? 0;
            }

            const itemSize = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol);

            var area;
            if (direction === VERTICAL) {
                area = new Area(0, lastPosition, itemSize.width, lastPosition + itemSize.height);
            } else if (direction === HORIZONTAL) {
                area = new Area(lastPosition, 0, lastPosition + itemSize + width, itemSize.height);
            }

            this.passedThroughPipeline.layout.positions.set(itemSymbol, area);

            if (direction === VERTICAL) {
                // crossAxisMinPosition = Math.min(crossAxisMinPosition, itemSize.width);
                crossAxisMaxPosition = Math.max(crossAxisMaxPosition, itemSize.width);
                lastPosition += itemSize.height;
            } else if (direction === HORIZONTAL) {
                // crossAxisMinPosition = Math.min(crossAxisMinPosition, itemSize.height);
                crossAxisMaxPosition = Math.max(crossAxisMaxPosition, itemSize.height);
                lastPosition += itemSize.width;
            }

            // if (marginsToApply.afterPageContent) {
            //     const margin = this.config.spacing.container.after;
            //     lastPosition += margin ? margin : 0;
            // }
            lastCellKind = currentCellKind;
        }

        if (direction === VERTICAL) {
            const averageHeight = itemMeasurer.getAverageSize(this.controlledByEnvironment.environmentSymbol).height;
            const afterPlacementHeight = averageHeight * remainingUnplacedItems;
            lastPosition += afterPlacementHeight;
        } else if (direction === HORIZONTAL) {
            const averageWidth = itemMeasurer.getAverageSize(this.controlledByEnvironment.environmentSymbol).width;
            const afterPlacementWidth = averageWidth * remainingUnplacedItems;
            lastPosition += afterPlacementWidth;
        }

        if (direction === VERTICAL) {
            this.passedThroughPipeline.contentBoundingBoxSize.width = crossAxisMaxPosition;
            this.passedThroughPipeline.contentBoundingBoxSize.height = lastPosition;
        } else if (direction === HORIZONTAL) {
            this.passedThroughPipeline.contentBoundingBoxSize.width = lastPosition;
            this.passedThroughPipeline.contentBoundingBoxSize.height = crossAxisMaxPosition;
        }
    }

    setSpacing(CellTypeSymbol) {}
}
