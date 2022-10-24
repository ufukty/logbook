import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";

export class FocusStabilizer extends AbstractLayoutMutator {
    constructor() {
        super();

        this.config = {
            ...this.config,
            /** @type {Symbol} */
            focusedItemSymbol: undefined,
        };

        this.currentShift = 0;
    }

    calculate() {
        var isFocusedItemPassed = false;
        var totalDistanceToFocusedElement = 0;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            if (!isFocusedItemPassed) {
                totalDistanceToFocusedElement += position;
            }

            if (itemSymbol === this.config.focusedItemSymbol) isFocusedItemPassed = true;
        }
    }
}
