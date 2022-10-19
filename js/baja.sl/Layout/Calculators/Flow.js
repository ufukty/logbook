import { iota } from "../../utilities.js";
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
export const HORIZONTAL = iota();
/** Meant to be used for `Flow.config.direction` */
export const VERTICAL = iota();

export class Flow extends AbstractLayoutCalculator {
    constructor() {
        super();
        this.config = {
            ...this.config,
            /** @type {Map.<CellTypeSymbol, Spacing>} */
            spacing: new Map(),
            /** Either `VERTICAL` or `HORIZONTAL` */
            direction: VERTICAL,
        };
    }

    perform() {
        console.log("flow-perform", this.config.placement);
        var lastPosition = 0;
        var lastCellKind = undefined;
        // var crossAxisMaxPosition = 0;

        if (this.config.direction === VERTICAL) {
            const averageHeight = this.config.averageSizeForUnplacedItem.height;
            const beforePlacementHeight = averageHeight * this.config.offset;
            lastPosition += beforePlacementHeight;
        } else if (this.config.direction === HORIZONTAL) {
            const averageWidth = this.config.averageSizeForUnplacedItem.width;
            const beforePlacementWidth = averageWidth * this.config.offset;
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

            // save item position
            if (this.config.direction === VERTICAL)
                area = new Area(0, lastPosition, itemSize.width, lastPosition + itemSize.height);
            else if (this.config.direction === HORIZONTAL)
                area = new Area(lastPosition, 0, lastPosition + itemSize + width, itemSize.height);

            this.passedThroughPipeline.layout.positions.set(itemSymbol, area);

            if (this.config.direction === VERTICAL) {
                // // crossAxisMinPosition = Math.min(crossAxisMinPosition, itemSize.width);
                // // crossAxisMaxPosition = Math.max(crossAxisMaxPosition, itemSize.width);
                lastPosition += itemSize.height;
            } else if (this.config.direction === HORIZONTAL) {
                // // crossAxisMinPosition = Math.min(crossAxisMinPosition, itemSize.height);
                // // crossAxisMaxPosition = Math.max(crossAxisMaxPosition, itemSize.height);
                lastPosition += itemSize.width;
            }

            // if (marginsToApply.afterPageContent) {
            //     const margin = this.config.spacing.container.after;
            //     lastPosition += margin ? margin : 0;
            // }
            lastCellKind = currentCellKind;
        }

        const remainingUnplacedItems =
            this.config.totalNumberOfItems - this.config.offset - this.config.placement.length;

        if (this.config.direction === VERTICAL) {
            const averageHeight = this.config.averageSizeForUnplacedItem.height;
            const afterPlacementHeight = averageHeight * remainingUnplacedItems;
            lastPosition += afterPlacementHeight;
        } else if (this.config.direction === HORIZONTAL) {
            const averageWidth = this.config.averageSizeForUnplacedItem.width;
            const afterPlacementWidth = averageWidth * remainingUnplacedItems;
            lastPosition += afterPlacementWidth;
        }

        // crossAxisMinPosition = Infinity;
        // crossAxisMaxPosition = -Infinity;
        this.passedThroughPipeline.layout.pageHeight = lastPosition;
    }
}
