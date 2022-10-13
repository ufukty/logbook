import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";
import { Area } from "../Coordinates.js";

export class MeasureContainer extends AbstractLayoutMutator {
    perform() {
        var container = new Area(Infinity, Infinity, -Infinity, -Infinity);

        for (const [itemSymbol, item] of this.passedThroughPipeline.layout.positions) {
            if (item.x0 < container.x0) container.x0 = item.x0;
            if (item.y0 < container.y0) container.y0 = item.y0;
            if (container.x1 < item.x1) container.x1 = item.x1;
            if (container.y1 < item.y1) container.y1 = item.y1;
        }

        this.passedThroughPipeline.containerSize = container;
    }
}
