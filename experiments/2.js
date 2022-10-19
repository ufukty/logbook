import { clamp, executeWhenDocumentIsReady, inverseLerp, lerp } from "../js/baja.sl/utilities.js";

/** @type {Array.<HTMLElement>} */
const routeContainers = [...document.getElementsByClassName("route-container")];

const elementUpdater = (routeContainer, carLight, reflectiveLayers, clampedMidPointRouteContainer) => {
    const relativePosX = routeContainer.offsetWidth / 2;
    const relativePosY = lerp(0, routeContainer.offsetHeight, clampedMidPointRouteContainer);

    const maskPositionX = `${-routeContainer.offsetWidth / 2 + relativePosX}`;
    const maskPositionY = `${-routeContainer.offsetHeight / 2 + relativePosY}`;

    carLight.style.transform = `translateY(calc(-50% + ${relativePosY}px))`;

    reflectiveLayers[0].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[0].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
    reflectiveLayers[1].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;

    routeContainer.dataset.clampedMidPointRouteContainer = clampedMidPointRouteContainer;
};

var scrollPositionForLastUpdate = undefined;
const adjustReflectionsWithScrollPosition = () => {
    if (scrollPositionForLastUpdate && scrollPositionForLastUpdate === window.scrollY) return;

    const midPointViewport = window.scrollY + 0.33 * window.innerHeight;

    routeContainers.forEach((routeContainer) => {
        /** @type {HTMLElement} */
        const carLight = routeContainer.querySelector(".route-car-light");
        /** @type {HTMLElement} */
        const reflectiveLayers = routeContainer.querySelectorAll(".reflective-layers");

        const midPointRouteContainer = inverseLerp(
            routeContainer.offsetTop,
            routeContainer.offsetTop + routeContainer.offsetHeight,
            midPointViewport
        );
        const clampedMidPointRouteContainer = clamp(midPointRouteContainer, -0.5, 1.5);

        if (
            routeContainer.dataset.clampedMidPointRouteContainer &&
            routeContainer.dataset.clampedMidPointRouteContainer === clampedMidPointRouteContainer
        )
            return;

        requestAnimationFrame(
            elementUpdater.bind(undefined, routeContainer, carLight, reflectiveLayers, clampedMidPointRouteContainer)
        );
    });

    scrollPositionForLastUpdate = window.scrollY;
};

window.addEventListener("scroll", adjustReflectionsWithScrollPosition);
window.addEventListener("load", adjustReflectionsWithScrollPosition);
