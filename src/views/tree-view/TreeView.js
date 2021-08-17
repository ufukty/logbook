import React from "react";

import "./TreeView.css";

import Task from "./../../ui-components/task-group/task-list/task/Task";

class IndentationGroup extends React.Component {
    constructor(props) {
        super();
        this.state = {
            root: props.root,
            dataset: props.dataset,
            children: props.dataset
                .filter((child) => child.parent === props.root.id)
                .map((child) => (
                    <IndentationGroup dataset={props.dataset} root={child} />
                )),
        };
        console.log(this.state.root.text);
        console.log(this.state.dataset);
        console.log(this.state.children.length);
    }

    render() {
        return (
            <div className="task-indentation-group">
                <details open="true">
                    <summary>
                        <Task key={this.state.root.id} task={this.state.root} />
                    </summary>
                    <div className="task-indentation-margin">
                        {this.state.children}
                    </div>
                </details>
            </div>
        );
    }
}

class TreeView extends React.Component {
    constructor(props) {
        super();
        this.state = {
            dataset: props.dataset,
        };
    }
    render() {
        return (
            <div>
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[0]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[2]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[3]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[4]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[5]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[6]}
                />
                <IndentationGroup
                    dataset={this.state.dataset}
                    root={this.state.dataset[7]}
                />
            </div>
        );
    }
}

export default TreeView;
