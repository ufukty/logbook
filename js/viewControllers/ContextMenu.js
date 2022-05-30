import { createElement, toggleAnimationWithClass } from "../utilities.js";

import AbstractViewController from "./AbstractViewController.js";

class ContextMenu extends AbstractViewController {
    constructor() {
        super();
        this.container = document.getElementById("context-menu");
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

    show() {
        this.container.style.visibility = "visible";
        toggleAnimationWithClass(this.container, "appearing", "context-menu-appear");
    }

    hide() {
        this.container.style.visibility = "hidden";
    }
}

export default ContextMenu;
