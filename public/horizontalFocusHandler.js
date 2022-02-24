function leastOfGreaterBinarySearch(list, item) {
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
    return lo;
}

class HorizontalFocusHandler {
    constructor() {
        this.lastFocusedElem = null;
        window.addEventListener("scroll", this.scrollEventListener.bind(this));
    }

    detectDOMElements() {
        this.infiniteSheet = document.getElementsByClassName("auto-focus-shift-area")[0];
    }

    positionOfFieldOfFocus() {
        var upperEdgeOfScreen = window.scrollY,
            heightOfScreen = window.innerHeight;
        return upperEdgeOfScreen + 0.3 * heightOfScreen;
    }

    debug(elements, boundaries, lastFocusedElem, position) {
        this.debugElem = document.querySelector("#debug");
        console.log({
            elements: elements,
            boundaries: boundaries,
            lastFocusedElem: lastFocusedElem,
            position: position,
        });
        // console.log("centerY", centerY, " focus area boundaries: ", this.getBoundariesOfFocusAreas());
        // this.debugElem.innerHTML = centerY;
        if (this.lastFocusedElem != undefined) {
            this.debugElem.innerHTML =
                "Focused: " +
                this.lastFocusedElem.innerHTML +
                "<br>Focus Depth: " +
                this.lastFocusedElem.getAttribute("task_depth") +
                "<br>Position: " +
                position;
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
        this.detectDOMElements();
        var elements = this.getFocusableElements();
        var boundaries = this.getBoundariesOfFocusAreas(elements);
        var focusIndex = leastOfGreaterBinarySearch(boundaries, this.positionOfFieldOfFocus());
        if (focusIndex != -1) {
            this.changeTaskInFocus(elements[focusIndex]);
        }

        this.debug(elements, boundaries, this.lastFocusedElem, this.positionOfFieldOfFocus());
    }

    changeTaskInFocus(domObject) {
        if (this.lastFocusedElem != undefined) {
            this.lastFocusedElem.classList.remove("focused-task");
        }
        domObject.classList.add("focused-task");
        this.lastFocusedElem = domObject;
        this.focusDepth = this.lastFocusedElem.getAttribute("task_depth");
        this.infiniteSheet.style.setProperty("--focus-depth", this.focusDepth);
    }

    getTheItemIndexShouldBeInFocus(boundaries) {
        var scrollPos = this.positionOfFieldOfFocus();
        for (const [itemIndex, boundary] of boundaries.entries()) {
            if (boundary.top <= scrollPos && scrollPos <= boundary.bottom) {
                return itemIndex;
            }
        }
        return -1;
    }
}

var horizontalFocusHandler;
document.addEventListener("DOMContentLoaded", function () {
    horizontalFocusHandler = new HorizontalFocusHandler();
});
