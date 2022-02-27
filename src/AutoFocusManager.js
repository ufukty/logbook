import * as constants from "./constants";

function leastOfGreaterBinarySearch(list, item) {
    list.unshift(-1); // Add dummy -1 at the beginning of array
    var lo = 0,
        hi = list.length - 1,
        mid = -1;
    while (hi - lo > 1) {
        mid = Math.floor((lo + hi) / 2);
        if (list[mid] <= item) {
            lo = mid;
        } else {
            hi = mid;
        }
    }
    return lo - 1;
}

export function calculateShiftForItem(adjustmentDepth, itemDepth) {
    console.log("adjustmentDepth assumed: ", adjustmentDepth);
    return constants.AUTO_FOCUS_SHIFT_IN_PIXELS * (itemDepth - adjustmentDepth);
}

export class AutoFocusManager {
    constructor() {
        // super();
        this.lastFocusedElem = undefined;
        window.addEventListener("scroll", this.scrollEventListener.bind(this));
        this.callbackForFocusChange = undefined;
        this.callbackForAdjustmentDepthChange = undefined;
        this.effectiveDVM = constants.DVM_CHRONO;
    }

    detectDOMElements() {
        this.infiniteSheet = document.getElementsByClassName("auto-focus-shift-area")[0];
    }

    positionOfFieldOfFocus() {
        var upperEdgeOfScreen = window.scrollY,
            heightOfScreen = window.innerHeight;
        return upperEdgeOfScreen + 0.3 * heightOfScreen;
    }

    debug() {
        var debugActivated = true;
        if (debugActivated) {
            this.debugElem = document.querySelector("#debug");
            if (this.lastFocusedElem !== undefined) {
                this.debugElem.innerHTML =
                    "Foc.Dep.: " +
                    this.focusDepth +
                    "<br>" +
                    "Adj.Dep.: " +
                    this.adjustmentForDepth +
                    "<br>" +
                    "PosFieldOfFocus: " +
                    this.positionOfFieldOfFocus();
            }
        }
    }

    // focusElements should be a NodeList
    getBoundariesOfFocusAreas(focusElements) {
        return [...focusElements].map((task) => task.offsetTop);
    }

    getFocusableElements() {
        return document.querySelectorAll(".task");
    }

    scrollEventListener(e) {
        e.preventDefault();
        this.detectDOMElements();
        var elements = this.getFocusableElements();
        var boundaries = this.getBoundariesOfFocusAreas(elements);
        var focusIndex = leastOfGreaterBinarySearch(boundaries, this.positionOfFieldOfFocus());
        if (focusIndex !== -1) {
            this.changeTaskInFocus(elements[focusIndex]);
        }
        this.debug();
    }

    changeTaskInFocus(domObject) {
        if (this.lastFocusedElem === domObject) {
            // same object in focus area
            return;
        }
        if (this.lastFocusedElem !== undefined) {
            this.lastFocusedElem.classList.remove("focused-task");
        }
        domObject.classList.add("focused-task");
        this.lastFocusedElem = domObject;
        this.focusDepth = this.lastFocusedElem.getAttribute("task_depth");
        this.infiniteSheet.style.setProperty("--focus-depth", this.focusDepth - this.adjustmentForDepth);
        this.debug();
    }

    // getTheItemIndexShouldBeInFocus(boundaries) {
    //     var scrollPos = this.positionOfFieldOfFocus();
    //     for (const [itemIndex, boundary] of boundaries.entries()) {
    //         if (boundary.top <= scrollPos && scrollPos <= boundary.bottom) {
    //             return itemIndex;
    //         }
    //     }
    //     return -1;
    // }

    registerNextDVM(nextDVM) {
        if (this.effectiveDVM === nextDVM) {
            return;
        }
        this.effectiveDVM = nextDVM;
        if (nextDVM === constants.DVM_CHRONO) {
            this.beforeTurningAutoFocusOff();
        } else if (nextDVM === constants.DVM_HIERARCH) {
            this.beforeTurningAutoFocusOn();
        }
        this.debug();
    }

    beforeTurningAutoFocusOff() {
        this.adjustmentForDepth = this.focusDepth;
        console.log("adjustmentDepth recorded: ", this.adjustmentForDepth);
        this.notifyDelegatesAtAdjustmentDepthChange();
    }

    beforeTurningAutoFocusOn() {
        this.adjustmentForDepth = this.focusDepth;
        console.log("adjustmentDepth recorded: ", this.adjustmentForDepth);
        this.notifyDelegatesAtAdjustmentDepthChange();
    }

    registerDelegateForAdjustmentDepthChange(callback) {
        this.callbackForAdjustmentDepthChange = callback;
    }

    notifyDelegatesAtAdjustmentDepthChange() {
        // console.log("fired 2");
        this.callbackForAdjustmentDepthChange(this.adjustmentForDepth);
    }
}

export var autoFocusManager;
document.addEventListener("DOMContentLoaded", function () {
    autoFocusManager = new AutoFocusManager();
});
