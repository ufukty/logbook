import TaskDay from "./ui-components/task-day/TaskDay";

import "./App.css";
import React from "react";

var task_events = [
    {
        log_id: 1,
        event_type: "new child",
        event_time: "198821988219882",
    },
];

var dataset = {
    days: [
        {
            day: "13 January",
            tasks: [
                {
                    text: "Revoke passwordless sudo rights after provision at cluster",
                    id: 0,
                    parent: 0,
                    created_at: "445884458844588",
                },
                {
                    text: "PAM for SSH",
                    id: 1,
                    parent: 4,
                    created_at: "337543375433754",
                },
                {
                    text: "iptables for redis",
                    id: 2,
                    parent: 1,
                    created_at: "425042504250",
                },
                {
                    text: "terraform for redis",
                    id: 3,
                    parent: 5,
                    created_at: "391839183918",
                },
                {
                    text: "ACL - Redis",
                    id: 4,
                    parent: 5,
                    created_at: "324363243632436",
                },
                {
                    text: "Update redis/tf file according to prod.tfvars file",
                    id: 5,
                    parent: 5,
                    created_at: "227322273222732",
                },
                {
                    text: "Redis security",
                    id: 6,
                    parent: 3,
                    created_at: "334063340633406",
                },
                {
                    text: "TOTP for SSH",
                    id: 7,
                    parent: 2,
                    created_at: "880588058805",
                },
                {
                    text: "API gateway without redis",
                    id: 8,
                    parent: 6,
                    created_at: "582358235823",
                },
                {
                    text: "Golden image interitance re-organize",
                    id: 9,
                    parent: 5,
                    created_at: "360893608936089",
                },
                {
                    text: "Postgres",
                    id: 10,
                    parent: 7,
                    created_at: "607006070060700",
                },
                {
                    text: "Auth service",
                    id: 11,
                    parent: 4,
                    created_at: "359643596435964",
                },
                {
                    text: "MQ",
                    id: 12,
                    parent: 0,
                    created_at: "996499649964",
                },
                {
                    text: "Federated learning",
                    id: 13,
                    parent: 7,
                    created_at: "649286492864928",
                },
                {
                    text: "Bluetooth transmission test",
                    id: 14,
                    parent: 3,
                    created_at: "475904759047590",
                },
                {
                    text: "Intrusion detection system (centralised) (OSSEC",
                    id: 15,
                    parent: 6,
                    created_at: "450134501345013",
                },
                {
                    text: "Envoy - HAProxy - NGiNX",
                    id: 16,
                    parent: 6,
                    created_at: "339853398533985",
                },
            ],
        },
    ],

    todo_drawer: [],
};

class DocumentSheet extends React.Component {
    constructor() {
        super();
        this.state = {
            children: dataset.days.map((data) => (
                <TaskDay key={data.day} data={data} />
            )),
        };
    }

    render() {
        return <div>{this.state.children}</div>;
    }
}
function App() {
    return (
        <div className="document-sheet">
            <a
                id="home-button"
                class="floating-corner left top"
                href="index.html"
            >
                Logbook
            </a>

            <div id="sheet-settings" class="floating-corner right top dark">
                <div>Share</div>

                <div>Sync</div>
            </div>

            <div id="active-task-details" class="floating-corner left bottom">
                History for active task
                <div class="task">PAM for SSH</div>
                <div class="task">ACL - Redis</div>
                <div class="task">TOTP for SSH</div>
            </div>

            <div id="date-anchors" class="floating-corner right bottom">
                <a href="#august-13-2021">13th August</a>
                <a href="#august-14-2021">14th August</a>
            </div>

            <DocumentSheet />
        </div>
    );
}

export default App;
