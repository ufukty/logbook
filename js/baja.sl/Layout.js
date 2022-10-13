import { Position, Size } from "./Layout/Coordinates.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { AbstractLayoutDecorator, AbstractLayoutMutator } from "./Layout/AbstractLayoutCalculator.js";
import { itemAccountant } from "./Layout/ItemAccountant.js";

export class PlaceholderViewController extends AbstractViewController {
    constructor() {
        super();

        this.dom = {
            ...this.dom,
            container: createElement("div", ["placeholder"]),
        };

        this.config = {
            ...this.config,
            /** @type {Size} */
            size: undefined,
        };
    }

    updateView() {
        this.dom.container.style.width = this.config.size.width;
        this.dom.container.style.height = this.config.size.height;
    }
}

export class Indentation extends AbstractLayoutMutator {
    constructor() {
        super();
    }

    calculate() {}
}

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
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            if (!isFocusedItemPassed) {
                totalDistanceToFocusedElement += position;
            }

            if (itemSymbol === this.config.focusedItemSymbol) isFocusedItemPassed = true;
        }
    }
}

export class AvatarLayout extends AbstractLayoutDecorator {
    constructor() {
        super();
    }
}


export class Padding extends AbstractLayoutMutator {
    /**
     * @param {number} top
     * @param {number} left
     * @param {number} right
     * @param {number} bottom
     */
    constructor(top = undefined, left = undefined, right = undefined, bottom = undefined) {
        super();

        this.top = top ?? 0;
        this.left = left ?? 0;
        this.right = right ?? 0;
        this.bottom = bottom ?? 0;
    }

    perform() {
        for (const position of this.passedThroughPipeline.layout.positions.values()) {
            position.x += this.left;
            position.y += this.top;
        }
        this.passedThroughPipeline.containerSize.height += this.top + this.bottom;
        this.passedThroughPipeline.containerSize.width += this.left + this.right;
    }
}
