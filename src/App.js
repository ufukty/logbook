import "./App.css";

import ChronologicalView from "./views/chronological-view/Chronological";
import TreeView from "./views/tree-view/TreeView";

import React from "react";

// var task_events = [
//     {
//         log_id: 1,
//         event_type: "new child",
//         event_time: "198821988219882",
//     },
// ];

var dataset = {
    days: [
        {
            day: "10 August",
            tasks: [
                {
                    text: "Revoke passwordless sudo rights after provision at cluster",
                    id: 0,
                    parent: 0,
                    created_at: "445884458844588",
                    active: false,
                },
                {
                    text: "iptables for redis",
                    id: 2,
                    parent: 1,
                    created_at: "425042504250",
                    active: false,
                },
                {
                    text: "terraform for redis",
                    id: 3,
                    parent: 5,
                    created_at: "391839183918",
                    active: false,
                },
                {
                    text: "Update redis/tf file according to prod.tfvars file",
                    id: 5,
                    parent: 5,
                    created_at: "227322273222732",
                    active: false,
                },
            ],
        },
        {
            day: "12 August",
            tasks: [
                {
                    text: "Remove: seperator from ovpn-auth",
                    id: 15,
                    parent: 3,
                    created_at: "475904759047590",
                    active: false,
                },
                {
                    text: "Write tests for ovpn-auth",
                    id: 16,
                    parent: 6,
                    created_at: "450134501345013",
                    active: false,
                },
                {
                    text: "Decrease timing gap of ovpn-auth under 1ms",
                    id: 17,
                    parent: 6,
                    created_at: "339853398533985",
                    active: false,
                },
                {
                    text: "Prepare releases for ovpn-auth",
                    id: 18,
                    parent: 6,
                    created_at: "339853398533985",
                    active: false,
                },
                {
                    text: "Provision golden-image for gitlab-runner",
                    id: 19,
                    parent: 6,
                    created_at: "339853398533985",
                    active: false,
                },
            ],
        },
        {
            day: "13 August",
            tasks: [
                {
                    text: "gitlab-runner --(vpn)--> DNS ----> gitlab",
                    id: 13,
                    parent: 0,
                    created_at: "996499649964",
                    active: false,
                },
                {
                    text: "Firewall & unbound rules update from prov script (VPN)",
                    id: 14,
                    parent: 7,
                    created_at: "649286492864928",
                    active: false,
                },
                {
                    text: "Script pic_gitlab_runner_post_creation",
                    id: 15,
                    parent: 3,
                    created_at: "475904759047590",
                    active: false,
                },
                {
                    text: "Execute 1 CI/CD pipeline on gitlab-runner",
                    id: 16,
                    parent: 6,
                    created_at: "450134501345013",
                    active: false,
                },
                {
                    text: "gitlab-runner provisioner with resolv.conf/docker/runner-register",
                    id: 17,
                    parent: 6,
                    created_at: "339853398533985",
                    active: false,
                },
                {
                    text: "prepare gitlab-ci for ovpn-auth repo",
                    id: 19,
                    parent: 6,
                    created_at: "339853398533985",
                    active: false,
                },
            ],
        },
    ],
    active_tasks: [
        {
            text: "PAM for SSH",
            id: 1,
            parent: 4,
            created_at: "337543375433754",
            active: true,
        },
        {
            text: "ACL - Redis",
            id: 4,
            parent: 5,
            created_at: "324363243632436",
            active: true,
        },
    ],
    todo_drawer: [
        {
            text: "Redis security",
            id: 6,
            parent: 3,
            created_at: "334063340633406",
            active: false,
        },
        {
            text: "TOTP for SSH",
            id: 7,
            parent: 2,
            created_at: "880588058805",
            active: false,
        },
        {
            text: "API gateway without redis",
            id: 8,
            parent: 6,
            created_at: "582358235823",
            active: false,
        },
        {
            text: "Golden image interitance re-organize",
            id: 9,
            parent: 5,
            created_at: "360893608936089",
            active: false,
        },
        {
            text: "Postgres",
            id: 10,
            parent: 7,
            created_at: "607006070060700",
            active: false,
        },
        {
            text: "Auth service",
            id: 11,
            parent: 4,
            created_at: "359643596435964",
            active: false,
        },
        {
            text: "MQ",
            id: 12,
            parent: 0,
            created_at: "996499649964",
            active: false,
        },
        {
            text: "Federated learning",
            id: 13,
            parent: 7,
            created_at: "649286492864928",
            active: false,
        },
        {
            text: "Bluetooth transmission test",
            id: 14,
            parent: 3,
            created_at: "475904759047590",
            active: false,
        },
        {
            text: "Intrusion detection system (centralised) (OSSEC",
            id: 15,
            parent: 6,
            created_at: "450134501345013",
            active: false,
        },
        {
            text: "Envoy - HAProxy - NGiNX",
            id: 16,
            parent: 6,
            created_at: "339853398533985",
            active: false,
        },
        {
            text: "web-front/Privacy against [friend/pubic/company/attackers]",
            id: 13,
            parent: 0,
            created_at: "996499649964",
            active: false,
        },
        {
            text: "Redis/cluster script test for multi datacenter",
            id: 14,
            parent: 7,
            created_at: "649286492864928",
            active: false,
        },
        {
            text: "gitlab-runner firewall rules: close public internet",
            id: 18,
            parent: 6,
            created_at: "339853398533985",
            active: false,
        },
        {
            text: "static-challange for ovpn-auth",
            id: 18,
            parent: 6,
            created_at: "339853398533985",
            active: false,
        },
        {
            text: "Golden image for vpn server",
            id: 19,
            parent: 6,
            created_at: "339853398533985",
            active: false,
        },
    ],
};

var dataset_tree = [
    {
        text: "Revoke passwordless sudo rights after provision at cluster",
        id: 0,
        parent: -1,
        created_at: "445884458844588",
        active: true,
    },
    {
        text: "iptables for redis",
        id: 2,
        parent: 0,
        created_at: "425042504250",
        active: true,
    },
    {
        text: "terraform for redis",
        id: 3,
        parent: 2,
        created_at: "391839183918",
        active: true,
    },
    {
        text: "Update redis/tf file according to prod.tfvars file",
        id: 5,
        parent: 2,
        created_at: "227322273222732",
        active: true,
    },
    {
        text: "gitlab-runner --(vpn)--> DNS ----> gitlab",
        id: 13,
        parent: 6,
        created_at: "996499649964",
        active: true,
    },
    {
        text: "Firewall & unbound rules update from prov script (VPN)",
        id: 14,
        parent: 2,
        created_at: "649286492864928",
        active: true,
    },
    {
        text: "Script pic_gitlab_runner_post_creation",
        id: 15,
        parent: 3,
        created_at: "475904759047590",
        active: true,
    },
    {
        text: "Execute 1 CI/CD pipeline on gitlab-runner",
        id: 16,
        parent: 4,
        created_at: "450134501345013",
        active: true,
    },
    {
        text: "gitlab-runner provisioner with resolv.conf/docker/runner-register",
        id: 17,
        parent: 1,
        created_at: "339853398533985",
        active: true,
    },
    {
        text: "Remove: seperator from ovpn-auth",
        id: 18,
        parent: 2,
        created_at: "475904759047590",
        active: true,
    },
    {
        text: "Write tests for ovpn-auth",
        id: 19,
        parent: 1,
        created_at: "450134501345013",
        active: true,
    },
    {
        text: "Decrease timing gap of ovpn-auth under 1ms",
        id: 20,
        parent: 0,
        created_at: "339853398533985",
        active: true,
    },
    {
        text: "Prepare releases for ovpn-auth",
        id: 21,
        parent: 3,
        created_at: "339853398533985",
        active: true,
    },
    {
        text: "Provision golden-image for gitlab-runner",
        id: 22,
        parent: 5,
        created_at: "339853398533985",
        active: true,
    },
    {
        text: "prepare gitlab-ci for ovpn-auth repo",
        id: 23,
        parent: 5,
        created_at: "339853398533985",
        active: false,
    },
    {
        text: "PAM for SSH",
        id: 24,
        parent: 3,
        created_at: "337543375433754",
        active: true,
    },
    {
        text: "ACL - Redis",
        id: 25,
        parent: 7,
        created_at: "324363243632436",
        active: true,
    },
    {
        text: "Redis security",
        id: 26,
        parent: 2,
        created_at: "334063340633406",
        active: false,
    },
    {
        text: "TOTP for SSH",
        id: 27,
        parent: 3,
        created_at: "880588058805",
        active: false,
    },
    {
        text: "API gateway without redis",
        id: 28,
        parent: 0,
        created_at: "582358235823",
        active: false,
    },
    {
        text: "Golden image interitance re-organize",
        id: 29,
        parent: 3,
        created_at: "360893608936089",
        active: false,
    },
    {
        text: "Postgres",
        id: 30,
        parent: 3,
        created_at: "607006070060700",
        active: false,
    },
    {
        text: "Auth service",
        id: 31,
        parent: 4,
        created_at: "359643596435964",
        active: false,
    },
    {
        text: "MQ",
        id: 32,
        parent: 2,
        created_at: "996499649964",
        active: false,
    },
    {
        text: "Federated learning",
        id: 33,
        parent: 3,
        created_at: "649286492864928",
        active: false,
    },
    {
        text: "Bluetooth transmission test",
        id: 34,
        parent: 6,
        created_at: "475904759047590",
        active: false,
    },
    {
        text: "Intrusion detection system (centralised) (OSSEC",
        id: 35,
        parent: 28,
        created_at: "450134501345013",
        active: false,
    },
    {
        text: "Envoy - HAProxy - NGiNX",
        id: 36,
        parent: 7,
        created_at: "339853398533985",
        active: false,
    },
    {
        text: "web-front/Privacy against [friend/pubic/company/attackers]",
        id: 37,
        parent: 7,
        created_at: "996499649964",
        active: false,
    },
    {
        text: "Redis/cluster script test for multi datacenter",
        id: 38,
        parent: 7,
        created_at: "649286492864928",
        active: false,
    },
    {
        text: "gitlab-runner firewall rules: close public internet",
        id: 39,
        parent: 5,
        created_at: "339853398533985",
        active: false,
    },
    {
        text: "static-challange for ovpn-auth",
        id: 40,
        parent: 6,
        created_at: "339853398533985",
        active: false,
    },
    {
        text: "Golden image for vpn server",
        id: 41,
        parent: 3,
        created_at: "339853398533985",
        active: false,
    },
];

class App extends React.Component {
    constructor() {
        super();
        this.state = {
            mode: "tree",
        };
    }

    render() {
        var content;
        if (this.state.mode === "chronological") {
            content = <ChronologicalView dataset={dataset} />;
        } else if (this.state.mode === "tree") {
            content = <TreeView dataset={dataset_tree} />;
        }
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

                <div
                    id="active-task-details"
                    class="floating-corner left bottom"
                >
                    History for active task
                    <div class="task">PAM for SSH</div>
                    <div class="task">ACL - Redis</div>
                    <div class="task">TOTP for SSH</div>
                </div>

                <div id="date-anchors" class="floating-corner right bottom">
                    <a href="#august-10-2021">10th August</a>
                    <a href="#august-12-2021">12th August</a>
                    <a href="#august-13-2021">13th August</a>
                    <a href="#august-14-2021">Active Tasks</a>
                    <a href="#august-14-2021">To-do Drawer</a>
                </div>

                {content}
            </div>
        );
    }
}

export default App;
