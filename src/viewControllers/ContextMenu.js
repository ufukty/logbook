import React from "react";

import "../css/context-menu.css"

class ContextMenu extends React.Component {
    constructor() {
        super()
        this.setEventListener()
        this.state = {
            style: {
                top: 0,
                left: 0,
                visibility: "hidden",
            }
        }
    }

    setEventListener() {
        document.addEventListener("contextmenu", this.contextMenuEventHandler.bind(this));
        document.addEventListener("click", this.clickEventListener.bind(this));
    }

    contextMenuEventHandler(e) {
        if (e.type !== "contextmenu" || !e.target.classList.contains("task")) {
            return
        }
        e.preventDefault();
        this.setState(prevState => {
            return {
                style: {
                    top: e.pageY,
                    left: e.pageX,
                    visibility: "visible",
                }
            }
        })
        this.presentationHandler()
        console.log(e.target.getAttribute("task_id"))
    }

    presentationHandler() {
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

    clickEventListener(e) {
        if (e.type !== "click") {
            return
        }
        this.setState(prevState => {
            return {
                style: {
                    visibility: "hidden",
                }
            }
        })
    }


    render() {
        return <div id="context-menu" style={this.state.style}>

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