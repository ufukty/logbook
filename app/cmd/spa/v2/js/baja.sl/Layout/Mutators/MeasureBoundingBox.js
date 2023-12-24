import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";
import { Area } from "../Coordinates.js";

export class MeasureBoundingBox extends AbstractLayoutMutator {
    perform() {
        var box = new Area(Infinity, Infinity, -Infinity, -Infinity);

        for (const [itemSymbol, position] of this.passedThroughPipeline.layout.positions) {
            if (position.x0 < box.x0) box.x0 = position.x0;
            if (position.y0 < box.y0) box.y0 = position.y0;
            if (box.x1 < position.x1) box.x1 = position.x1;
            if (box.y1 < position.y1) box.y1 = position.y1;
        }

        this.passedThroughPipeline.contentBoundingBoxSize = box;
    }
}
