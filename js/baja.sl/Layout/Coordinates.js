import { lerp } from "../utilities.js";

export class Size {
    /**
     * @param {number} width
     * @param {number} height
     */
    constructor(width, height) {
        this.width = width;
        this.height = height;
    }

    /**
     * @param {number} width
     * @param {number} height
     */
    isEqual(width, height) {
        return width === this.width && height === this.height;
    }

    /**
     * @param {Size} other
     */
    isEqualWith(other) {
        return other.width === this.width && other.height === this.height;
    }
}

export class Position {
    constructor(x, y) {
        this.set(x, y);
    }

    set(x, y) {
        this.x = x;
        this.y = y;
    }

    add(x, y) {
        this.x += x;
        this.y += y;
    }

    /** @param {Position} position  */
    addFrom(position) {
        this.x += position.x;
        this.y += position.y;
    }

    subtract(x, y) {
        this.x -= x;
        this.y -= y;
    }

    /** @param {Position} position  */
    subtractFrom(position) {
        this.x -= position.x;
        this.y -= position.y;
    }

    deltaComp(x, y) {
        return [this.x - x, this.y - y];
    }

    /**
     * @param {Position} position
     * @returns {Size}
     */
    deltaCompFrom(position) {
        return Size(this.x - position.x, this.y - position.y);
    }

    delta(x, y) {
        const [dx, dy] = deltaComp(x, y);
        return Math.sqrt(Math.pow(dx, 2) + Math.pow(dy, 2));
    }

    /**
     * @param {Position} position
     * @returns {number}
     */
    deltaFrom(position) {
        const dimensions = deltaCompFrom(position);
        return Math.sqrt(Math.pow(dx, 2) + Math.pow(dy, 2));
    }
}

export class Anchor {
    constructor(horizontal, vertical) {
        this.horizontal = horizontal;
        this.vertical = vertical;
    }

    /**
     * @param {Size} areaSize
     * @returns {Position}
     */
    interpolate(areaSize) {
        return new Position(lerp(0, areaSize.width, this.horizontal), lerp(0, areaSize.height, this.vertical));
    }
}

export class Distance {
    constructor(horizontal, vertical) {
        this.horizontal = horizontal;
        this.vertical = vertical;
    }
}

export class Spacing {
    constructor(before, between, after) {
        this.before = before;
        this.after = after;
        this.between = between;
    }
}

export class Area {
    /**
     * @param {number} x0
     * @param {number} y0
     * @param {number} x1
     * @param {number} y1
     */
    constructor(x0, y0, x1, y1) {
        this.x0 = x0;
        this.y0 = y0;
        this.x1 = x1;
        this.y1 = y1;
        // this.size = new Size(x1 - x0, y1 - y0);
    }

    /** @param {Area} area  */
    isCollidingWith(area) {
        return this.x0 <= area.x1 && area.x0 <= this.x1 && this.y0 <= area.y1 && area.y0 <= this.y1;
    }

    /**
     * @param {number} factor
     * @param {Anchor} transformOrigin
     * This method updates the start and end positions of represented area,
     *   keeping original transformOrigin at same position.
     */
    scale(factor, transformOrigin = new Anchor(0.5, 0.5)) {
        const origin = transformOrigin.interpolate(new Size(this.width(), this.height()));
        origin.add(this.x0, this.y0);
        this.x0 = factor * (this.x0 - origin.x) + origin.x;
        this.y0 = factor * (this.y0 - origin.y) + origin.y;
        this.x1 = factor * (this.x1 - origin.x) + origin.x;
        this.y1 = factor * (this.y1 - origin.y) + origin.y;
        // this.size.width = this.x1 - this.x0;
        // this.size.height = this.y1 - this.y0;
        return this;
    }

    translate(horizontal = 0, vertical = 0) {
        this.x0 += horizontal;
        this.x1 += horizontal;
        this.y0 += vertical;
        this.y1 += vertical;
        return this;
    }

    moveTo(x = undefined, y = undefined) {
        if (x) {
            const width = this.width();
            this.x0 = x;
            this.x1 = x + width;
        }
        if (y) {
            const height = this.height();
            this.y0 = y;
            this.y1 = y + height;
        }
        return this;
    }

    height() {
        return this.y1 - this.y0;
    }

    width() {
        return this.x1 - this.x0;
    }
}
