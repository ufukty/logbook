export function executeWhenDocumentIsReady(func) {
    if (
        document.readyState === "complete" ||
        document.readyState === "loaded" ||
        document.readyState === "interactive"
    ) {
        func();
    } else {
        document.addEventListener("DOMContentLoaded", func);
    }
}

class DOMElementReuseCollector {
    constructor() {
        /**  @type {AbstractViewController[]} */
        this._freeItems = {};
        this._constructors = {};
    }

    /** @param {String} itemTypeIdentifier */
    registerItemIdentifier(itemTypeIdentifier, constructor) {
        if (this._freeItems.hasOwnProperty(itemTypeIdentifier)) {
            return;
        }
        this._freeItems[itemTypeIdentifier] = [];
        this._constructors[itemTypeIdentifier] = constructor;
    }

    _createItemController(itemTypeIdentifier) {
        const _constructor = this._constructors[itemTypeIdentifier];
        const item = _constructor();
        return item;
    }

    /**
     * @param {String} itemTypeIdentifier
     * @returns {AbstractViewController}
     */
    get(itemTypeIdentifier) {
        let itemController;
        if (this._freeItems[itemTypeIdentifier].length === 0) {
            itemController = this._createItemController(itemTypeIdentifier);
        } else {
            itemController = this._freeItems[itemTypeIdentifier].pop();
        }
        if (typeof itemController.prepareForUse !== "undefined" && typeof itemController.prepareForUse === "function")
            itemController.prepareForUse();
        return itemController;
    }

    /**
     * @param {String} itemTypeIdentifier
     * @param {} itemController
     */
    free(itemTypeIdentifier, itemController) {
        if (typeof itemController.prepareForFree !== "undefined" && typeof itemController.prepareForFree === "function")
            itemController.prepareForFree();
        this._freeItems[itemTypeIdentifier].push(itemController);
    }
}

export const domCollector = new DOMElementReuseCollector();

/**
 * @param {String} tagName
 * @param {String[]} classList
 * @param {HTMLElement[]} childrenList
 */
export function createElement(tagName, classList, childrenList) {
    const element = document.createElement(tagName);
    if (classList !== undefined && classList.length > 0)
        classList.forEach((classStr) => {
            element.classList.add(classStr);
        });
    if (childrenList !== undefined && childrenList.length > 0)
        childrenList.forEach((childNode) => element.appendChild(childNode));
    return element;
}

/** @returns {HTMLElement} */
export function createElementInSVGNamespace(tag, classList, childrenList) {
    const element = document.createElementNS("http://www.w3.org/2000/svg", tag);
    if (classList !== undefined && classList.length > 0)
        classList.forEach((classStr) => {
            element.classList.add(classStr);
        });
    if (childrenList !== undefined && childrenList.length > 0)
        childrenList.forEach((childNode) => element.appendChild(childNode));
    return element;
}

/**
 * Adds classes and appends children nodes to <element>
 * @param {HTMLElement} element
 * @param {String[]} classList
 * @param {HTMLElement[]} childrenList
 * @returns {HTMLElement}
 */
export function configureElement(element, classList, childrenList) {
    if (classList !== undefined && classList.length > 0)
        classList.forEach((classStr) => {
            element.classList.add(classStr);
        });
    if (childrenList !== undefined && childrenList.length > 0)
        childrenList.forEach((childNode) => element.appendChild(childNode));
    return element;
}

export function configureElementWithProps(element, classList, style, nodeProperties, childrenList) {}

/**
 * Adds classes and appends children nodes to <element>
 * @param {HTMLElement} element
 * @param {HTMLElement[]} childrenList
 */
export function adoption(element, childrenList) {
    if (childrenList !== undefined) {
        if (Array.isArray(childrenList)) {
            if (childrenList.length > 0) childrenList.forEach((childNode) => element.appendChild(childNode));
        } else {
            element.appendChild(childrenList);
        }
    }
    return element;
}

/**
 * Adds classes and appends children nodes to <element>
 * @param {HTMLElement} element
 * @param {{}} props
 */
export function setStyleProperties(element, props) {
    for (const key in props) {
        if (Object.hasOwnProperty.call(props, key)) {
            const value = props[key];
            element.style.setProperty(key, value);
        }
    }
    return element;
}
/**
 * Adds classes and appends children nodes to <element>
 * @param {HTMLElement} element
 * @param {{}} props
 */
export function setAttributes(element, props) {
    for (const key in props) {
        if (Object.hasOwnProperty.call(props, key)) {
            const value = props[key];
            element.setAttribute(key, value);
        }
    }
    return element;
}

/**
 * @param {HTMLElement} element
 * @param {String} triggerClass
 * @param {HTMLElement} animatedElement
 * @param {String} animationName
 * @param {Function|undefined} callback
 */
export function toggleAnimationWithParentClass(
    parentElement,
    triggerClass,
    animatedElement,
    animationName,
    callback = undefined
) {
    animatedElement.addEventListener("animationend", function eventHandler(e) {
        if (e.animationName === animationName) {
            parentElement.classList.remove(triggerClass);
            animatedElement.removeEventListener("animationend", eventHandler);
            if (callback !== undefined) callback();
        }
    });
    parentElement.classList.add(triggerClass);
}

/**
 * @param {HTMLElement} element
 * @param {String} triggerClass
 * @param {String} animationName
 * @param {Function|undefined} callback
 */
export function toggleAnimationWithClass(element, triggerClass, animationName, callback = undefined) {
    toggleAnimationWithParentClass(element, triggerClass, element, animationName, callback);
}

/**
 * Unlike toggleAnimationWith(Parent)Class, startAnimationWith(Parent)Class functions don't remove
 * the class after animation stop
 * @param {HTMLElement} element Provided class name will be added and removed to this element
 * @param {String} triggerClass
 * @param {HTMLElement} animatedElement "animationend" event listener will be added to this element
 * @param {String} animationName
 * @param {Function|undefined} callback
 */
export function startAnimationWithParentClass(
    parentElement,
    triggerClass,
    animatedElement,
    animationName,
    callback = undefined
) {
    animatedElement.addEventListener("animationend", function eventHandler(e) {
        if (e.animationName === animationName) {
            animatedElement.removeEventListener("animationend", eventHandler);
            if (callback !== undefined) callback();
        }
    });
    parentElement.classList.add(triggerClass);
}

/**
 * Unlike toggleAnimationWith(Parent)Class, startAnimationWith(Parent)Class functions don't remove
 * the class after animation stop
 * @param {HTMLElement} element Provided class name will be added and removed to this element
 * @param {String} triggerClass
 * @param {String} animationName "animationend" event listener will be added to this element
 * @param {Function|undefined} callback
 */
export function startAnimationWithClass(element, triggerClass, animationName, callback = undefined) {
    startAnimationWithParentClass(element, triggerClass, element, animationName, callback);
}

export class CSSAnimationTrigger {
    constructor(element, triggerClassName, animationName, animatedElement) {
        this.callOnEnd = [];
    }

    onStop(callback) {
        if (this.status !== "end") {
            this.callOnEnd.push(callback);
        } else {
            callback();
        }
        return this;
    }

    start() {
        return this;
    }

    removeTrigger() {
        return this;
    }
}

export function arrayFromRange(limit) {
    return [...Array(limit).keys()];
}

/**
 * @param {String[]} keyList
 */
export function createAnObjectOfLists(keyList) {
    let objectOfLists = {};
    keyList.forEach((key) => {
        objectOfLists[key] = [];
    });
    return objectOfLists;
}

export function assert(statement, errorMessage) {
    if (statement === false) {
        console.error(errorMessage);
    }
}

export function setRootCSSVariable(prop, value) {
    document.documentElement.style.setProperty(prop, value);
}

function fillCSSVariablesForViewportSize() {
    let vh = window.innerHeight * 0.01;
    let vw = window.innerWidth * 0.01;
    // unitless versions are for using in properties that accept unitless
    // values such as transform: translate(0.5)
    setRootCSSVariable("--vw", `${vw}`);
    setRootCSSVariable("--vh", `${vh}`);
    setRootCSSVariable("--vw-px", `${vw}px`);
    setRootCSSVariable("--vh-px", `${vh}px`);
}

let historicalViewportDimensions = {
    "p": {
        // portrait orientation
        "min-vh": undefined,
        "max-vh": undefined,
        "min-vw": undefined,
        "max-vw": undefined,
    },
    "l": {
        // landscape orientation
        "min-vh": undefined,
        "max-vh": undefined,
        "min-vw": undefined,
        "max-vw": undefined,
    },
};

function fillCSSVariablesForDynamicViewportSize() {
    let vh = window.innerHeight * 0.01;
    let vw = window.innerWidth * 0.01;
    // unitless versions are for using in properties that accept unitless
    // values such as transform: translate(0.5)
    setRootCSSVariable("--dynamic-vw", `${vw}`);
    setRootCSSVariable("--dynamic-vh", `${vh}`);
    setRootCSSVariable("--dynamic-vw-px", `${vw}px`);
    setRootCSSVariable("--dynamic-vh-px", `${vh}px`);

    let orientation;
    if (vw > vh) {
        orientation = "l";
    } else {
        orientation = "p";
    }

    if (historicalViewportDimensions[orientation]["max-vw"] === undefined) {
        historicalViewportDimensions[orientation]["min-vh"] = vh;
        historicalViewportDimensions[orientation]["max-vh"] = vh;
        historicalViewportDimensions[orientation]["min-vw"] = vw;
        historicalViewportDimensions[orientation]["max-vw"] = vw;
    }

    if (vw > historicalViewportDimensions[orientation]["max-vw"]) {
        historicalViewportDimensions[orientation]["max-vw"] = vw;
    }
    if (historicalViewportDimensions[orientation]["min-vw"] > vw) {
        historicalViewportDimensions[orientation]["min-vw"] = vw;
    }
    if (vh > historicalViewportDimensions[orientation]["max-vh"]) {
        historicalViewportDimensions[orientation]["max-vh"] = vh;
    }
    if (historicalViewportDimensions[orientation]["min-vh"] > vh) {
        historicalViewportDimensions[orientation]["min-vh"] = vh;
    }

    ["min-vh", "max-vh", "min-vw", "max-vw"].forEach((key) => {
        setRootCSSVariable(`--${key}`, historicalViewportDimensions[orientation][key]);
    });
}

// TODO: wrap those with executeWhenDocumentIsReady
fillCSSVariablesForViewportSize();
fillCSSVariablesForDynamicViewportSize();

window.addEventListener("resize", function () {
    fillCSSVariablesForDynamicViewportSize();
});

export function addEventListenerForNonTouchScreen(targetElement, eventType, callback, options) {
    executeWhenDocumentIsReady(function () {
        targetElement.addEventListener(eventType, callback, options);
        targetElement.addEventListener("touchstart", function () {
            targetElement.removeEventListener(eventType, callback, options);
        });
    });
}

class PersistentSymbolizer {
    constructor() {
        // TODO: Implement LRU cache to reduce memory allocation
        this.cache = new Map();
        this.cacheReverse = new Map();
    }

    /**
     * @param {string} value
     * @returns {Symbol}
     */
    symbolize(value) {
        if (this.cache.has(value)) {
            return this.cache.get(value);
        } else {
            const symbol = Symbol(value);
            this.cache.set(value, symbol);
            this.cacheReverse.set(symbol, value);
            return symbol;
        }
    }

    /**
     * @param {Symbol} value
     * @returns {string}
     */
    desymbolize(symbol) {
        return this.cacheReverse.get(symbol);
    }
}
export const symbolizer = new PersistentSymbolizer();

export function isInBetween(a, b, c) {
    if (a <= b && c <= c) return true;
    else return false;
}

export function checkCollision(item_y1, item_y2, viewport_y1, viewport_y2) {
    /*
            * * * * * * *  (y1)                     * * * * * * *  (y1)                
            *           *                           *           *              
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y1)
        +   *           *                           *           *           +
        +   * * * * * * *  (y2)                     *           *           +
        +                                           *           *           +
        +                                           *           *           +       <=  viewport
        +                                           *           *           +
        +                 * * * * * * *  (y1)       *           *           +
        +                 *           *             *           *           +
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y2)
                          *           *             *           *                              
                          * * * * * * *  (y2)       * * * * * * *  (y2)                               
    */
    if (item_y2 < viewport_y1 || item_y1 > viewport_y2)
        // if item starts after viewport ends, or item ends before viewport starts,
        // then the item is not in viewport.
        return false;
    else return true;
}

/** @returns {Set.<string>} */
export function mergeMapKeys() {
    let set_ = new Set();
    for (let i = 0; i < arguments.length; i++) {
        if (arguments[i]) for (const key of arguments[i].keys()) set_.add(key);
    }
    return set_;
}

/**
 * @param {Set} leftSet
 * @param {Set} rightSet
 * Returns a set represents intersection of two input sets. */
export function setIntersect(leftSet, rightSet) {
    const intersection = new Set();
    for (const k of leftSet) {
        if (rightSet.has(k)) intersection.add(k);
    }
    return intersection;
}
