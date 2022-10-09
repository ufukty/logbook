import { AbstractLayoutCalculator, AbstractLayoutMutator } from "../AbstractLayoutCalculator.js";

/** Meant to be used for `Flow.config.direction` */
export const HORIZONTAL = iota();
/** Meant to be used for `Flow.config.direction` */
export const VERTICAL = iota();

export class Flow extends AbstractLayoutCalculator {
    constructor() {
        super();
        this.config = {
            ...this.config,
            spacing: {
                container: new Spacing(10, 0, 10),
            },
            direction: VERTICAL,
        };
    }

    _performForDirectionVertical() {
        var lastPosition = this.config.spacing.container.before;
        var lastCellKind = undefined;
        var lastItemIndex = this.computedValues.layout.positions.length - 1;

        const averageHeight = this.getAverageHeightForAnItem();
        const beforePlacementHeight = averageHeight * this.config.placement.offset;
        lastPosition += beforePlacementHeight;

        for (const [itemIndex, itemSymbol] of this.computedValues.layout.positions.entries()) {
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
                const margin = this.config.spacing.container.before;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.beforeMarginForCurrentKind) {
                const margin = this.config.spacing[currentCellKind].before;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.afterMarginForPreviousKind) {
                const margin = this.config.spacing[lastCellKind].after;
                lastPosition += margin ? margin : 0;
            }
            if (marginsToApply.betweenSameKind) {
                const margin = this.config.spacing[currentCellKind].between;
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
                const margin = this.config.spacing.container.after;
                lastPosition += margin ? margin : 0;
            }
            lastCellKind = currentCellKind;
        }

        this.computedValues.next.pageHeight = lastPosition;
    }

    perform() {
        if (this.config.direction === VERTICAL) {
            this._performForDirectionVertical();
        } else if (this.config.direction === HORIZONTAL) {
            this._performForDirectionHorizontal();
        }
    }
}
