class Navigation {
    constructor() {
        const pageElements = document.querySelectorAll("[data-nav-handle]");
        this.pageOrdering = [];
        this.pageHandleToNode = new Map();
        pageElements.forEach((elem) => {
            this.pageHandleToNode.set(elem.dataset.navHandle, elem);
            this.pageOrdering.push(elem.dataset.navHandle);
        });

        this.pageSwitchAnimationConfig = {
            duration: 100,
            iterations: 1,
            fill: "both",
        };

        this.pageSwitchAnimationKeyframes = {
            "show-reverse": [
                { opacity: 0, transform: "translateX(-200px)" },
                { opacity: 1, transform: "translateX(0px)" },
            ],
            "show": [
                { opacity: 0, transform: "translateX(200px)" },
                { opacity: 1, transform: "translateX(0px)" },
            ],
            "hide-reverse": [
                { opacity: 1, transform: "translateX(0px)" },
                { opacity: 0, transform: "translateX(200px)" },
            ],
            "hide": [
                { opacity: 1, transform: "translateX(0px)" },
                { opacity: 0, transform: "translateX(-200px)" },
            ],
        };
        this.activePageIndex = -1;
        this.navigationStack = [];

        this.addEventListenersToPageElements();
    }

    addEventListenersToPageElements() {
        const elems = document.querySelectorAll("[data-nav-target]");
        elems.forEach((elem, num) => {
            elem.addEventListener("click", this.switchPages.bind(this, elem.dataset.navTarget));
        });
    }

    findIndexOfPageHandle(handle) {
        for (const [index, value] of this.pageOrdering.entries()) {
            if (value == handle) {
                return index;
            }
        }
        return -1;
    }

    switchPages(nextPageHandle) {
        if (nextPageHandle === "back") {
            if (this.navigationStack.length >= 2) {
                this.navigationStack.pop();
                nextPageHandle = this.navigationStack.pop();
            } else {
                console.error("can't go back because navigation stack is shorter than 1 pages");
                return;
            }
        } else if (nextPageHandle === "first") {
            nextPageHandle = this.navigationStack[0];
        }


        var nextPageIndex = this.findIndexOfPageHandle(nextPageHandle);
        if (nextPageIndex === -1) {
            if (this.activePageIndex !== -1) {
                return this.switchPages("first");
            } else {
                console.error("can't find index of this page handler: '" + nextPageHandle + "'");
                return;
            }
        }

        if (this.activePageIndex == nextPageIndex) {
            return;
        }

        const showPageCaller = () => {
            this.navigationStack.push(nextPageHandle);
            if (this.activePageIndex < nextPageIndex) {
                this.showPage(nextPageHandle);
            } else {
                this.showPage(nextPageHandle, true);
            }
            this.activePageIndex = nextPageIndex;
        };
        if (this.activePageIndex != -1) {
            var activePageHandle = this.pageOrdering[this.activePageIndex];
            if (this.activePageIndex < nextPageIndex) {
                this.hidePage(activePageHandle);
            } else {
                this.hidePage(activePageHandle, true);
            }
            setTimeout(showPageCaller, 100);
        } else {
            showPageCaller();
        }
    }

    showPage(pageHandle, reverseDirection) {
        const page = this.pageHandleToNode.get(pageHandle);
        page.classList.add("active-story-page");
        var keyframes;
        if (reverseDirection) {
            keyframes = this.pageSwitchAnimationKeyframes["show-reverse"];
        } else {
            keyframes = this.pageSwitchAnimationKeyframes["show"];
        }
        page.animate(keyframes, this.pageSwitchAnimationConfig);
    }

    hidePage(pageHandle, reverseDirection) {
        const page = this.pageHandleToNode.get(pageHandle);

        var keyframes;
        if (reverseDirection) {
            keyframes = this.pageSwitchAnimationKeyframes["hide-reverse"];
        } else {
            keyframes = this.pageSwitchAnimationKeyframes["hide"];
        }
        const animation = page.animate(keyframes, this.pageSwitchAnimationConfig);
        animation.finished.then(() => {
            page.classList.remove("active-story-page");
        });
    }
}

const nav = new Navigation();
nav.switchPages("landing");