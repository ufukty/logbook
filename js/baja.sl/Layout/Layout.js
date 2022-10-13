import { iota, symbolizer } from "../utilities.js";
import { AbstractLayoutPipe } from "./AbstractLayoutPipe.js";
import { AbstractLayoutCalculator, AbstractLayoutDecorator, AbstractLayoutMutator } from "./AbstractLayoutPipe.js";
import { UpdateScheduler } from "../UpdateScheduler.js";
import { Size } from "./Coordinates.js";

export class Layout {
    constructor() {
        /** @private */
        this.private = {
            /** @type {Array.<AbstractLayoutPipe>} */
            pipeline: [],
            /** @type {Array.<Function>} */
            subscriber: [],
            /** @type {UpdateScheduler} */
            updateScheduler: new UpdateScheduler(this._recalculate.bind(this), 60),
        };

        this.passedThroughPipeline = {
            /** @type {Layout} */
            layout: {
                /** @type {Map.<Symbol, Area>} */
                positions: new Map(),
                /** @type {Map.<Symbol, number>} */
                scaling: new Map(),
            },
            /**  @type {Size} */
            containerSize: new Size(),
        };

        this.environmentSymbol = symbolizer.symbolize(`environment#${iota()}`);
    }

    /** @private */
    _recalculate(trigger) {
        console.log("calculate");
        for (const pipe of this.private.pipeline) {
            pipe.passedThroughPipeline = this.passedThroughPipeline;
            pipe.perform();
            this.passedThroughPipeline = pipe.passedThroughPipeline;
        }
        this.private.updateScheduler.finished();
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

    /**
     *
     */
    start() {}

    /**
     * It is important to call this method when recalculating position calculation
     * for each resizeObserver notification is not needed anymore.
     */
    stop() {}
}
