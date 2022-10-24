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

export class Align extends AbstractLayoutMutator {
    /** @param {Symbol} alignTo */
    constructor(alignTo = HORIZONTAL_LEFT) {
        super();
        this.config = {
            ...this.config,
            /**
             * This describes which point of items and container will be aligned.
             * Look for exported constants by file.
             */
            alignTo: alignTo,
        };
    }

    _alignToHorizontalLeft() {
        const contentBoundingBoxSize = this.passedThroughPipeline.contentBoundingBoxSize.width;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            position.moveTo(0, undefined);
        }
    }

    _alignToHorizontalCenter() {
        const contentBoundingBoxSize = this.passedThroughPipeline.contentBoundingBoxSize.width;
        if (contentBoundingBoxSize === undefined) {
            console.error("can not align without contentBoundingBoxSize set");
            return;
        }
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemWidth = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).width;
            position.translate((contentBoundingBoxSize - itemWidth) / 2, 0);
        }
    }

    _alignToHorizontalRight() {
        const contentBoundingBoxSize = this.passedThroughPipeline.contentBoundingBoxSize.width;
        if (contentBoundingBoxSize === undefined) {
            console.error("can not align without contentBoundingBoxSize set");
            return;
        }
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemWidth = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).width;
            position.translate(contentBoundingBoxSize - itemWidth, 0);
        }
    }

    _alignToVerticalTop() {
        const containerHeight = this.passedThroughPipeline.contentBoundingBoxSize.height;
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            position.moveTo(undefined, 0);
        }
    }

    _alignToVerticalCenter() {
        const containerHeight = this.passedThroughPipeline.contentBoundingBoxSize.height;
        if (containerHeight === undefined) {
            console.error("can not align without contentBoundingBoxSize set");
            return;
        }
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemHeight = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).height;
            position.translate(0, (containerHeight - itemHeight) / 2);
        }
    }

    _alignToVerticalBottom() {
        const containerHeight = this.passedThroughPipeline.contentBoundingBoxSize.height;
        if (containerHeight === undefined) {
            console.error("can not align without contentBoundingBoxSize set");
            return;
        }
        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions.entries()) {
            const itemHeight = itemMeasurer.getSize(itemSymbol, this.controlledByEnvironment.environmentSymbol).height;
            position.translate(0, containerHeight - itemHeight);
        }
    }

    perform() {
        switch (this.config.alignTo) {
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
        this.config.alignTo = newAlignTo;
        this.controlledByEnvironment.environmentRef.scheduleRecalculation(TRIGGER_PIPE_CONFIG_CHANGE);
    }
}
