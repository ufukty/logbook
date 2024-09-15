import { createElement, toggleAnimationWithClass } from "../baja.sl/utilities.js";

import { AbstractViewController } from "../baja.sl/AbstractViewController.js";

export class ContextMenu extends AbstractViewController {
    constructor() {
        super();
        this.container = document.getElementById("context-menu");
        this.delegates = {
            delete: () => {}, // should be assigned by callee
        };
        const deleteButton = document.getElementById("context-menu-delete");
        deleteButton.addEventListener("click", (e) => {
            e.preventDefault();
            // alert("delete");
            this.delegates.delete(this.activeTaskId);
        });
        deleteButton.addEventListener("touchend", () => {
            alert("delete");
        });
    }

    setActiveTaskId(taskId) {
        this.activeTaskId = taskId;
    }

    setPosition(posX, posY) {
        this.container.style.top = `${posY}px`;
        this.container.style.left = `${posX}px`;
    }

    setTransformOrigin(transformOriginX, transformOriginY) {
        this.container.style.transformOrigin = `${transformOriginX}% ${transformOriginY}%`;
    }

    playAnimation() {
        const animatedElement = document.getElementById("context-menu");
        const animationName = "context-menu-appear";
        const triggerClass = "appearing";

        animatedElement.addEventListener("animationend", function eventHandler(e) {
            if (e.animationName === animationName) {
                animatedElement.classList.remove(triggerClass);
                animatedElement.removeEventListener("animationend", eventHandler);
            }
        });
        animatedElement.classList.add(triggerClass);
    }

    deleteButtonHandler() {}

    show() {
        this.container.style.visibility = "visible";
        toggleAnimationWithClass(this.container, "appearing", "context-menu-appear");
    }

    hide() {
        this.container.style.visibility = "hidden";
    }
}

export default ContextMenu;
