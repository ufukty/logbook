import React from "react";

import "./TaskList.css";

import Task from "./task/Task";

var data_set = [
    {
        text: "Revoke passwordless sudo rights after provision at cluster",
        id: 0,
        parent: 0,
        created_at: "445884458844588",
    },
    {
        text: "PAM for SSH",
        id: 4,
        parent: 4,
        created_at: "337543375433754",
    },
    {
        text: "iptables for redis",
        id: 4,
        parent: 1,
        created_at: "425042504250",
    },
    {
        text: "terraform for redis",
        id: 0,
        parent: 5,
        created_at: "391839183918",
    },
    {
        text: "ACL - Redis",
        id: 0,
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
        id: 3,
        parent: 6,
        created_at: "334063340633406",
    },
    {
        text: "TOTP for SSH",
        id: 5,
        parent: 7,
        created_at: "880588058805",
    },
    {
        text: "API gateway without redis",
        id: 17,
        parent: 6,
        created_at: "582358235823",
    },
    {
        text: "Golden image interitance re-organize",
        id: 15,
        parent: 5,
        created_at: "360893608936089",
    },
    {
        text: "Postgres",
        id: 16,
        parent: 7,
        created_at: "607006070060700",
    },
    {
        text: "Auth service",
        id: 16,
        parent: 4,
        created_at: "359643596435964",
    },
    {
        text: "MQ",
        id: 18,
        parent: 0,
        created_at: "996499649964",
    },
    {
        text: "Federated learning",
        id: 14,
        parent: 7,
        created_at: "649286492864928",
    },
    {
        text: "Bluetooth transmission test",
        id: 17,
        parent: 3,
        created_at: "475904759047590",
    },
    {
        text: "Intrusion detection system (centralised) (OSSEC",
        id: 18,
        parent: 6,
        created_at: "450134501345013",
    },
    {
        text: "Envoy - HAProxy - NGiNX",
        id: 15,
        parent: 6,
        created_at: "339853398533985",
    },
];

var task_events = [
    {
        log_id: 1,
        event_type: "new child",
        event_time: "198821988219882",
    },
];

function TaskList() {
    const task_items = data_set.map((data) => (
        <Task key={data.id} item={data} />
    ));
    return <div className="task-list">{task_items}</div>;
}

export default TaskList;
