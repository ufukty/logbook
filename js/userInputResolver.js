/* Moving 2 pixels left is treated as same with 1 pixel up + 1 pixel left. */
function manhattanDistance(x1, y1, x2, y2) {
    return Math.abs(x1 - x2) + Math.abs(y1 - y2);
}

export class UserInputResolver {
    constructor() {
        this.delegates = {
            openContextMenu: () => {},
            closeContextMenu: () => {},
        };
        this.flags = {
            contextMenuModeIsActive: false,
        };
        this.touchTracking = {
            totalDistanceTraveled: 0,
            isFingerMoved: false,
            lastPoint: { x: 0, y: 0 },
        };
    }

    enableContextMenuOnce(taskElement, clickPosX, clickPosY) {
        if (this.flags.contextMenuModeIsActive) return;
        this.flags.contextMenuModeIsActive = true;

        const taskPositionerElement = taskElement.parentNode;
        const objectId = taskPositionerElement.dataset["objectId"];

        this.taskElementThatIsContextMenuFocusedOn = taskPositionerElement;
        this.delegates.openContextMenu(taskPositionerElement, taskElement, objectId, clickPosX, clickPosY);
    }

    disableContextMenuOnce() {
        if (!this.flags.contextMenuModeIsActive) return;
        this.flags.contextMenuModeIsActive = false;

        this.delegates.closeContextMenu();
        this.taskElementThatIsContextMenuFocusedOn = undefined;
    }

    /** @param {MouseEvent} e */
    clickEventReceiverNonTouchScreen(e) {
        this.disableContextMenuOnce();
    }

    /** @param {MouseEvent} e */
    contextMenuEventReceiver(e) {
        e.preventDefault();
        // this.disableContextMenuOnce();
        const targetElement = e.target;

        console.log(e);
        // console.log(e);
        if (targetElement.classList.contains("infinite-sheet-task-content-area")) {
            e.stopPropagation();

            const taskElement = targetElement;
            const taskPositionerElement = taskElement.parentNode;

            if (this.taskElementThatIsContextMenuFocusedOn === taskPositionerElement) {
                this.disableContextMenuOnce();
            } else {
                this.disableContextMenuOnce();
                this.enableContextMenuOnce(taskElement, e.pageX, e.pageY);
            }
        } else {
            this.disableContextMenuOnce();
        }
    }

    /** @param {TouchEvent} e */
    touchStartEventReceiver(e) {
        const taskElement = e.currentTarget;
        // const touchStartTime = Date.now();
        // console.log(e);
        this.touchTracking.totalDistanceTraveled = 0.0;
        this.touchTracking.isFingerMoved = false;
        this.touchTracking.lastPoint.x = e.changedTouches[0].screenX;
        this.touchTracking.lastPoint.y = e.changedTouches[0].screenY;
    }

    /** @param {TouchEvent} e */
    touchMoveEventReceiver(e) {
        // console.log(e);
        if (this.touchTracking.isFingerMoved) return;
        const touch = e.changedTouches[0];
        const last = this.touchTracking.lastPoint;
        const lastMovementDistance = manhattanDistance(last.x, last.y, touch.screenX, touch.screenY);
        this.touchTracking.totalDistanceTraveled += lastMovementDistance;
        // console.log(this.touchTracking.totalDistanceTraveled);
        if (this.touchTracking.totalDistanceTraveled > 10) this.touchTracking.isFingerMoved = true;
        this.touchTracking.lastPoint.x = touch.screenX;
        this.touchTracking.lastPoint.y = touch.screenY;
    }

    /** @param {TouchEvent} e */
    touchEndEventReceiver(e) {
        const targetElement = e.target;

        if (targetElement.classList.contains("infinite-sheet-task-content-area")) {
            const taskElement = targetElement;
            const taskPositionerElement = taskElement.parentNode;

            if (this.touchTracking.isFingerMoved) return;

            if (this.taskElementThatIsContextMenuFocusedOn === taskPositionerElement) {
                this.disableContextMenuOnce();
            } else {
                this.disableContextMenuOnce();
                this.enableContextMenuOnce(taskElement, e.changedTouches[0].pageX, e.changedTouches[0].pageY);
            }
        } else {
            this.disableContextMenuOnce();
        }
    }
}
