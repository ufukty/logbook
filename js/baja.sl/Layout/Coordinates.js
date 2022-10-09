import { lerp } from "../utilities.js";

export class Size {
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
        set(x, y);
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
     */
    interpolate(areaSize) {
        return Position(lerp(0, areaSize.width, this.horizontal), lerp(0, areaSize.height, this.vertical));
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
