import { toHaveAccessibleDescription } from "@testing-library/jest-dom/dist/matchers";
import React from "react";

import "../css/context-menu.css"

class ContextMenu extends React.PureComponent {
    constructor(props) {
        super()
        // this.setEventListener()
        this.state = {
            posY: 0,
            posX: 0,
            enabled: "hidden",
        }
    }

    static getDerivedStateFromProps(props, state) {
        console.log(props, state)
        return {
            posY: props.posY,
            posX: props.posX,
            enabled: props.enabled,
        }
    }

    _presentationHandler() {
        const animatedElement = document.getElementById("context-menu");
        const animationName = "context-menu-appear"
        const triggerClass = "appearing";

        animatedElement.addEventListener("animationend", function eventHandler(e) {
            if (e.animationName === animationName) {
                animatedElement.classList.remove(triggerClass);
                animatedElement.removeEventListener("animationend", eventHandler);
            }
        });
        animatedElement.classList.add(triggerClass);
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState.enabled === this.state.enabled) {
            return
        }
        if (this.state.enabled) {
            this._presentationHandler()
        }
    }


    render() {
        console.log("render")
        const style = {
            top: `${this.state.posY}px`,
            left: `${this.state.posX}px`,
            visibility: this.enabled ? "visible" : "hidden",
        }
        return <div id="context-menu" style={style}>

            <div className="context-menu-item" title="Mark the task (and its sub-tasks if there are any) complete. "><div className="scale">Complete</div></div>
            <div className="context-menu-item" title="Attact this task to a parent task different from currently attached one." ><div className="scale">Reattach</div></div>

            <div className="seperator"></div>

            <div className="context-menu-item" title="Edit the text inside task. "><div className="scale">Edit task</div></div>
            <div className="context-menu-item" title="Copy the text inside of the task into system clipboard."><div className="scale">Copy text</div></div>

            <div className="seperator"></div>

            <div className="context-menu-item"><div className="scale">Reopen last Take</div></div>
            <div className="context-menu-item"><div className="scale">Start another Take</div></div>

            <div className="seperator"></div>

            <div className="context-menu-item" title="Invite others to complete this task together"><div className="scale">Invite others</div></div>
            <div className="context-menu-item" title="You can share a copy of your resolution for this goal with others. Each individual's progress are isolated and not shared with others."><div className="scale">Share as Blueprint</div></div>

            <div className="seperator"></div>

            <div className="context-menu-item" title="See parent task and all prerequisites of this task detailed view."><div className="scale">Inspect relations</div></div>
            <div className="context-menu-item" title="List of historical actions performed on this task like completion, reopens, shared or invited users."><div className="scale">Historical info</div></div>

            <div className="seperator"></div>

            <div className="context-menu-item destructive-action"><div className="scale">Delete</div></div>
        </div>
    }
}

export default ContextMenu;