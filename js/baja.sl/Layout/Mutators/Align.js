import { itemMeasurer } from "../../ItemMeasurer.js";
import { iota, symbolizer } from "../../utilities.js";
import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";
import { TRIGGER_PIPE_CONFIG_CHANGE } from "../Layout.js";

/** Meant to be used for `Align.alignTo`. To align left edge of items to the left edge of container. */
export const HORIZONTAL_LEFT = symbolizer.symbolize("HORIZONTAL_LEFT");
/** Meant to be used for `Align.alignTo`. To align horizontal center of items to the horizontal center of container. */
export const HORIZONTAL_CENTER = symbolizer.symbolize("HORIZONTAL_CENTER");
/** Meant to be used for `Align.alignTo`. To align right edge of items to the right edge of container. */
export const HORIZONTAL_RIGHT = symbolizer.symbolize("HORIZONTAL_RIGHT");

/** Meant to be used for `Align.alignTo`. To align top edge of items to the top edge of container. */
export const VERTICAL_TOP = symbolizer.symbolize("VERTICAL_TOP");
/** Meant to be used for `Align.alignTo`. To align vertical center of items to the vertical center of container. */
export const VERTICAL_CENTER = symbolizer.symbolize("VERTICAL_CENTER");
/** Meant to be used for `Align.alignTo`. To align left edge of items to the left edge of container. */
export const VERTICAL_BOTTOM = symbolizer.symbolize("VERTICAL_BOTTOM");

/** Items will be aligned accordingly to the content bounding box that is previously calculated in the pipeline. */
export const ALIGN_INTO_CONTENT_BOUNDING_BOX = symbolizer.symbolize("BASED_ON_CONTENT_BOUNDING_BOX");
/** Items will be aligned accordingly to the container. Make sure updating container size before layout calculation starts. */
export const ALIGN_INTO_CONTAINER = symbolizer.symbolize("BASED_ON_CONTAINER");

export class Align extends AbstractLayoutMutator {
    /** @param {Symbol} alignOn */
    constructor(alignOn = HORIZONTAL_LEFT, alignInto = ALIGN_INTO_CONTENT_BOUNDING_BOX) {
        super();
        this.config = {
            ...this.config,
            /**
             * This describes which point of items and container will be aligned.
             * Look for exported constants by file.
             */
            alignOn: alignOn,
            /**
             * Valid values are: `BASED_ON_CONTENT_BOUNDING_BOX` and `BASED_ON_CONTAINER`
             * @
             */
            alignInto: alignInto,
        };
    }

    /** @private */
    _returnBoxSize() {
        if (this.config.alignInto === ALIGN_INTO_CONTAINER) {
            const boxSize = this.passedThroughPipeline.containerSize;
            if (boxSize.width === undefined || boxSize.height === undefined)
                console.error("Alignment will be completed without containerSize is known");
            return boxSize;
        } else if (this.config.alignInto === ALIGN_INTO_CONTENT_BOUNDING_BOX) {
            const boxSize = this.passedThroughPipeline.contentBoundingBoxSize;
            if (boxSize.width === undefined || boxSize.height === undefined)
                console.error("Alignment will be completed without contentBoundingBoxSize is known");
            return boxSize;
        } else console.error("Invalid value for Align._returnBoxSize");
    }

    /** @private */
    _alignToHorizontalLeft() {
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            position.moveTo(0, undefined);
        }
    }

    /** @private */
    _alignToHorizontalCenter() {
        const boxSize = this._returnBoxSize().width;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemWidth = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).width;
            position.translate((boxSize - itemWidth) / 2, 0);
        }
    }

    /** @private */
    _alignToHorizontalRight() {
        const boxSize = this._returnBoxSize().width;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemWidth = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).width;
            position.translate(boxSize - itemWidth, 0);
        }
    }

    /** @private */
    _alignToVerticalTop() {
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            position.moveTo(undefined, 0);
        }
    }

    /** @private */
    _alignToVerticalCenter() {
        const boxSize = this._returnBoxSize().height;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemHeight = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).height;
            position.translate(0, (boxSize - itemHeight) / 2);
        }
    }

    /** @private */
    _alignToVerticalBottom() {
        const boxSize = this._returnBoxSize().height;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemHeight = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).height;
            position.translate(0, boxSize - itemHeight);
        }
    }

    perform() {
        switch (this.config.alignOn) {
            case HORIZONTAL_LEFT:
                this._alignToHorizontalLeft();
                break;
            case HORIZONTAL_CENTER:
                this._alignToHorizontalCenter();
                break;
            case HORIZONTAL_RIGHT:
                this._alignToHorizontalRight();
                break;
            case VERTICAL_TOP:
                this._alignToVerticalTop();
                break;
            case VERTICAL_BOTTOM:
                this._alignToVerticalBottom();
                break;
            case VERTICAL_CENTER:
                this._alignToVerticalCenter();
                break;
        }
    }

    /** @param {Symbol} newAlignTo */
    updateWithNewAlignment(newAlignTo) {
        this.controlledByEnvironment.pipeNeedsRefresh = true;
        this.config.alignOn = newAlignTo;
        this.controlledByEnvironment.environmentRef.scheduleRecalculation(TRIGGER_PIPE_CONFIG_CHANGE);
    }
}
