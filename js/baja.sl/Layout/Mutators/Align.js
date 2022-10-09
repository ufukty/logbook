import { AbstractLayoutMutator } from "../AbstractLayoutCalculator.js";
import { itemAccountant } from "../ItemAccountant.js";

/** Meant to be used for `Align.alignOn`. To align left edge of items to the left edge of container. */
export const HORIZONTAL_LEFT = iota();
/** Meant to be used for `Align.alignOn`. To align horizontal center of items to the horizontal center of container. */
export const HORIZONTAL_CENTER = iota();
/** Meant to be used for `Align.alignOn`. To align right edge of items to the right edge of container. */
export const HORIZONTAL_RIGHT = iota();

/** Meant to be used for `Align.alignOn`. To align top edge of items to the top edge of container. */
export const VERTICAL_TOP = iota();
/** Meant to be used for `Align.alignOn`. To align vertical center of items to the vertical center of container. */
export const VERTICAL_CENTER = iota();
/** Meant to be used for `Align.alignOn`. To align left edge of items to the left edge of container. */
export const VERTICAL_BOTTOM = iota();

export class Align extends AbstractLayoutMutator {
    constructor() {
        super();
        this.config = {
            ...this.config,
            /**
             * This describes which point of items and container will be aligned.
             * Look for exported constants by file.
             */
            alignOn: VERTICAL_CENTER,
        };
    }

    _justifyOnHorizontalLeft() {
        const containerWidth = this.config.containerSize.width;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            position.x = 0;
        }
    }

    _justifyOnHorizontalCenter() {
        const containerWidth = this.config.containerSize.width;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            const itemWidth = itemAccountant.getSize(itemSymbol, this.config.environmentSymbol).width;
            position.x = (containerWidth - itemWidth) / 2;
        }
    }

    _justifyOnHorizontalRight() {
        const containerWidth = this.config.containerSize.width;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            const itemWidth = itemAccountant.getSize(itemSymbol, this.config.environmentSymbol).width;
            position.x = containerWidth - itemWidth;
        }
    }

    _justifyOnVerticalTop() {
        const containerHeight = this.config.containerSize.height;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            position.y = 0;
        }
    }

    _justifyOnVerticalCenter() {
        const containerHeight = this.config.containerSize.height;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            const itemHeight = itemAccountant.getSize(itemSymbol, this.config.environmentSymbol).height;
            position.y = (containerHeight - itemHeight) / 2;
        }
    }

    _justifyOnVerticalBottom() {
        const containerHeight = this.config.containerSize.height;
        for (const [itemSymbol, position] of this.config.layout.positions.entries()) {
            const itemHeight = itemAccountant.getSize(itemSymbol, this.config.environmentSymbol).height;
            position.y = containerHeight - itemHeight;
        }
    }

    perform() {
        if (this.config.alignOn === HORIZONTAL_LEFT) {
            this._justifyOnHorizontalLeft();
        } else if (this.config.alignOn === HORIZONTAL_CENTER) {
            this._justifyOnHorizontalCenter();
        } else if (this.config.alignOn === HORIZONTAL_RIGHT) {
            this._justifyOnHorizontalRight();
        } else if (this.config.alignOn === VERTICAL_TOP) {
            this._justifyOnVerticalTop();
        } else if (this.config.alignOn === VERTICAL_BOTTOM) {
            this._justifyOnVerticalCenter();
        } else if (this.config.alignOn === VERTICAL_CENTER) {
            this._justifyOnVerticalBottom();
        }
    }
}
