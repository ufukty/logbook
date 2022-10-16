import { inverseLerp, lerp } from "../js/baja.sl/utilities.js";

/** @type {Array.<HTMLElement>} */
const routeContainers = [...document.getElementsByClassName("route-container")];

/** @param {HTMLElement} routeContainer*/
const initializes = (routeContainer) => {
    /** @type {HTMLElement} */
    const carLight = routeContainer.querySelector(".route-car-light");
    /** @type {HTMLElement} */
    const reflectiveLayers = routeContainer.querySelectorAll(".reflective-layers");

    const relativePosX = routeContainer.offsetWidth / 2;
    const relativePosY = routeContainer.offsetHeight;

    carLight.style.top = `${relativePosY}px`;
    carLight.style.left = `${relativePosX}px`;

    const maskPositionX = `${-routeContainer.offsetWidth / 2 + relativePosX}`;
    const maskPositionY = `${-routeContainer.offsetHeight / 2 + relativePosY}`;

    reflectiveLayers[0].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[0].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
};
/** @param {MouseEvent} e */
const mouseEnters = (e) => {
    /** @type {HTMLElement} */
    const routeContainer = e.currentTarget;
    /** @type {HTMLElement} */
    const carLight = routeContainer.querySelector(".route-car-light");
    /** @type {HTMLElement} */
    const reflectiveLayers = routeContainer.querySelectorAll(".reflective-layers");

    carLight.style.transition = "none";
    reflectiveLayers[0].style.transition = "none";
    reflectiveLayers[1].style.transition = "none";
};

/** @param {MouseEvent} e */
const mouseMoves = (e) => {
    /** @type {HTMLElement} */
    const routeContainer = e.currentTarget;
    const relativePosX = e.clientX - routeContainer.offsetLeft - window.scrollX;
    // const relativePosX = routeContainer.offsetWidth / 2;
    const relativePosY = e.clientY - routeContainer.offsetTop + window.scrollY;

    /** @type {HTMLElement} */
    const carLight = routeContainer.querySelector(".route-car-light");
    carLight.style.top = `${relativePosY}px`;
    carLight.style.left = `${relativePosX}px`;

    /** @type {HTMLElement} */
    const reflectiveLayers = routeContainer.querySelectorAll(".reflective-layers");

    const maskPositionX = `${-routeContainer.offsetWidth / 2 + relativePosX}`;
    const maskPositionY = `${-routeContainer.offsetHeight / 2 + relativePosY}`;

    reflectiveLayers[0].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[0].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
};

/** @param {MouseEvent} e */
const mouseLeaves = (e) => {
    /** @type {HTMLElement} */
    const routeContainer = e.currentTarget;
    /** @type {HTMLElement} */
    const carLight = routeContainer.querySelector(".route-car-light");
    /** @type {HTMLElement} */
    const reflectiveLayers = routeContainer.querySelectorAll(".reflective-layers");

    carLight.style.transition = "all 200ms ease-out";
    reflectiveLayers[0].style.transition = "all 200ms ease-out";
    reflectiveLayers[1].style.transition = "all 200ms ease-out";

    const relativePosX = routeContainer.offsetWidth / 2;
    const relativePosY = routeContainer.offsetHeight;

    carLight.style.top = `${relativePosY}px`;
    carLight.style.left = `${relativePosX}px`;

    const maskPositionX = `${-routeContainer.offsetWidth / 2 + relativePosX}`;
    const maskPositionY = `${-routeContainer.offsetHeight / 2 + relativePosY}`;

    reflectiveLayers[0].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[0].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
};
routeContainers.forEach((routeContainer) => {
    initializes(routeContainer);
    routeContainer.addEventListener("mouseenter", mouseEnters);
    routeContainer.addEventListener("mousemove", mouseMoves);
    // routeContainer.addEventListener("touchmove", mouseMoves);
    routeContainer.addEventListener("mouseleave", mouseLeaves);
});
