import { iota, symbolizer } from "../utilities.js";
import { AbstractLayoutPipe } from "./AbstractLayoutPipe.js";
import { AbstractLayoutCalculator, AbstractLayoutDecorator, AbstractLayoutMutator } from "./AbstractLayoutPipe.js";
import { UpdateScheduler } from "../UpdateScheduler.js";
import { Size } from "./Coordinates.js";
import { Align } from "./Mutators/Align.js";
import { DelegateRegistry } from "../DelegateRegistry.js";
import { itemMeasurer } from "../ItemMeasurer.js";

export const TRIGGER_NEW_CONTAINER_SIZE = symbolizer.symbolize("TRIGGER_NEW_CONTAINER_SIZE");
export const TRIGGER_PIPE_NEW_PLACEMENT = symbolizer.symbolize("TRIGGER_PIPE_NEW_PLACEMENT");
export const TRIGGER_PIPE_CONFIG_CHANGE = symbolizer.symbolize("TRIGGER_PIPE_CONFIG_CHANGE");
export const TRIGGER_ITEM_MEASURER_CALLBACK = symbolizer.symbolize("TRIGGER_ITEM_MEASURER_CALLBACK");

const NEW_LAYOUT_CALCULATED = symbolizer.symbolize("NEW_LAYOUT_CALCULATED");

class LayoutPipeRecalculationNeedChecker {
    constructor() {
        this._isFirstChangePassed = false;
    }

    /** @param {AbstractLayoutPipe} pipe */
    doesPipeNeedRecalculation(pipe) {
        if (!this._isFirstChangePassed) {
            if (pipe.controlledByEnvironment.pipeNeedsRefresh) {
                this._isFirstChangePassed = true;
            } else {
                return false;
            }
        }
        if (pipe.controlledByEnvironment.pipeNeedsRefresh) {
            pipe.controlledByEnvironment.pipeNeedsRefresh = false;
        }

        return true;
    }
}

export class Layout {
    constructor() {
        /** @private */
        this.private = {
            /** @type {Array.<AbstractLayoutPipe>} */
            pipeline: [],
            /** @type {Array.<Function>} */
            subscriber: [],
            /**  */
            delegates: new DelegateRegistry(),
            /** @type {UpdateScheduler} */
            updateScheduler: new UpdateScheduler(this._recalculate.bind(this), 60),
            processAtNextUpdate: {
                containerSize: undefined,
            },
        };

        this.passedThroughPipeline = {
            current: this._getTemplateForDataPassedThroughPipeline(),
            next: this._getTemplateForDataPassedThroughPipeline(),
        };

        this.environmentSymbol = symbolizer.symbolize(`environment#${iota()}`);

        itemMeasurer.subscribe(this.environmentSymbol, () => {
            this.scheduleRecalculation(TRIGGER_ITEM_MEASURER_CALLBACK);
        });
    }

    /** @private */
    _getTemplateForDataPassedThroughPipeline() {
        return {
            layout: {
                /** @type {Map.<Symbol, Area>} */
                positions: new Map(),
                /** @type {Map.<Symbol, number>} */
                scaling: new Map(),
            },
            /**  @type {Size} */
            contentBoundingBoxSize: new Size(),
            /**  @type {Size} */
            containerSize: this.private.processAtNextUpdate.containerSize,
        };
    }

    /** @private */
    _recalculate(trigger) {
        // const checker = new LayoutPipeRecalculationNeedChecker();
        this.passedThroughPipeline.next = this._getTemplateForDataPassedThroughPipeline();

        if (this.private.processAtNextUpdate.containerSize !== undefined) {
            this.passedThroughPipeline.next.containerSize = this.private.processAtNextUpdate.containerSize;
            // this.private.processAtNextUpdate.containerSize = undefined;
        }

        for (const pipe of this.private.pipeline) {
            // if (checker.doesPipeNeedRecalculation(pipe)) {
            pipe.passedThroughPipeline = this.passedThroughPipeline.next;
            pipe.perform();
            this.passedThroughPipeline.next = pipe.passedThroughPipeline;
            // }
        }
        this.private.updateScheduler.finished();
        this.private.delegates.notify(NEW_LAYOUT_CALCULATED);
    }

    scheduleRecalculation(trigger) {
        this.private.updateScheduler.schedule(trigger);
    }

    /**
     * @private
     * @param {AbstractLayoutPipe} pipe
     */
    _connectPipeToPipeline(pipe) {
        this.private.pipeline.push(pipe);
        pipe.controlledByEnvironment.environmentSymbol = this.environmentSymbol;
        pipe.controlledByEnvironment.environmentRef = this;
    }

    /** @param {AbstractLayoutCalculator} pipe */
    connectCalculator(pipe) {
        if (!pipe instanceof AbstractLayoutCalculator)
            console.error("Given pipe is not an instance of AbstractLayoutCalculator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /** @param {AbstractLayoutMutator} pipe */
    connectMutator(pipe) {
        if (!pipe instanceof AbstractLayoutMutator)
            console.error("Given pipe is not an instance of AbstractLayoutMutator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /** @param {AbstractLayoutDecorator} pipe */
    connectDecorator(pipe) {
        if (!pipe instanceof AbstractLayoutDecorator)
            console.error("Given pipe is not an instance of AbstractLayoutDecorator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /** @param {Size} newSize */
    setContainerSize(newSize) {
        this.private.processAtNextUpdate.containerSize = newSize;
    }

    /** @param {function} callback - This will be called after each layout update. */
    subscribe(callback) {
        this.private.delegates.add(NEW_LAYOUT_CALCULATED, callback);
    }
}
