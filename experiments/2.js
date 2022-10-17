import { clamp, executeWhenDocumentIsReady, inverseLerp, lerp } from "../js/baja.sl/utilities.js";

/** @type {Array.<HTMLElement>} */
const routeContainers = [...document.getElementsByClassName("route-container")];

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
            routeContainer.clampedMidPointRouteContainer &&
            routeContainer.clampedMidPointRouteContainer === clampedMidPointRouteContainer
        ) {
            return;
        }

        const relativePosX = routeContainer.offsetWidth / 2;
        const relativePosY = lerp(0, routeContainer.offsetHeight, clampedMidPointRouteContainer);

        const maskPositionX = `${-routeContainer.offsetWidth / 2 + relativePosX}`;
        const maskPositionY = `${-routeContainer.offsetHeight / 2 + relativePosY}`;

        carLight.style.top = `${relativePosY}px`;
        // carLight.style.left = `${relativePosX}px`;

        reflectiveLayers[0].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
        reflectiveLayers[1].style.webkitMaskPosition = `${maskPositionX}px ${maskPositionY}px`;
        reflectiveLayers[0].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;
        reflectiveLayers[1].style.maskPosition = `${maskPositionX}px ${maskPositionY}px`;

        // console.log(midPointRouteContainer);
        routeContainer.clampedMidPointRouteContainer = clampedMidPointRouteContainer;
    });

    scrollPositionForLastUpdate = window.scrollY;
};

window.addEventListener("scroll", adjustReflectionsWithScrollPosition);
window.addEventListener("load", adjustReflectionsWithScrollPosition);
