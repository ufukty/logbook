DROP DATABASE IF EXISTS testdatabase;

CREATE DATABASE testdatabase;

\c testdatabase;

CREATE TABLE "DOCUMENT" (
    "document_id"       UUID UNIQUE DEFAULT gen_random_UUID(), 
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ACCESS"(
    "document_id"       UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "user-agent"        VARCHAR(256),
    "ip-address"        INET NOT NULL
);

CREATE TYPE GROUP_TYPE AS ENUM('archive', 'drawer', 'ready-to-start', 'active', 'paused', 'dropped');

CREATE TABLE "TASK_GROUP" (
    "task_group_id"     UUID UNIQUE DEFAULT gen_random_UUID(),
    "document_id"       UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
    "task_group_type"   GROUP_TYPE NOT NULL,
    "created_at"        DATE DEFAULT CURRENT_DATE
);

CREATE TYPE TASK_STATUS AS ENUM('archive', 'drawer', 'ready-to-start', 'active', 'paused', 'dropped');

CREATE TABLE "TASK" (
    "task_id"           UUID UNIQUE DEFAULT gen_random_UUID(),
    "task_group_id"     UUID NOT NULL REFERENCES "TASK_GROUP" ("task_group_id"),
    "parent_id"         UUID DEFAULT '00000000-0000-0000-0000-000000000000',
    "content"           TEXT NOT NULL,
    "task_status"       TASK_STATUS NOT NULL,
    "degree"            INT NOT NULL DEFAULT 1,
    "depth"             INT NOT NULL DEFAULT 1,
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE VIEW tasks_linearized AS SELECT * FROM tasks;

-- INSERT INTO "USER"("user_id", "password") VALUES ('0842c266-af1b-41bc-b180-653ca42dff82', '123456789');
-- INSERT INTO "USER"("password") VALUES ('123456789');
-- INSERT INTO "USER"("password") VALUES ('123456789');

-- INSERT INTO "TASK_GROUP"(task_group_id, document_id, group_name) VALUES('$1', '$2', '$3')

-- INSERT INTO "DOCUMENT"("user_id", "document_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', '7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Test Document');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'Donec eleifend est ac facilisis malesuada.');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'Cras a lorem sed arcu pretium congue sed sit amet mi.');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'Fusce lacinia quam sed maximus venenatis.');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'In nec tellus viverra');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'Cras porttitor nisl et urna viverra');
-- INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES('0842c266-af1b-41bc-b180-653ca42dff82', 'Praesent finibus lorem a ornare dapibus.');

-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'deploy redis cluster on multi DC', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'deploy redis cluster on 1 DC', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Revoke passwordless sudo rights after provision at cluster', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'iptables for redis', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'terraform for redis', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Update redis/tf file according to prod.tfvars file', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Remove: seperator from ovpn-auth', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Write tests for ovpn-auth', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Decrease timing gap of ovpn-auth under 1ms', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Prepare releases for ovpn-auth', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Provision golden-image for gitlab-runner', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'gitlab-runner --(vpn)--> DNS ----> gitlab', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Firewall & unbound rules update from prov script (VPN)', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Script pic_gitlab_runner_post_creation', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Execute 1 CI/CD pipeline on gitlab-runner', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'gitlab-runner provisioner with resolv.conf/docker/runner-register', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'prepare gitlab-ci for ovpn-auth repo', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'PAM for SSH', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'ACL - Redis', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Redis security', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'TOTP for SSH', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'API gateway without redis', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Golden image interitance re-organize', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Postgres', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Auth service', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'MQ', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Federated learning', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Bluetooth transmission test', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Intrusion detection system (centralised) (OSSEC', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Envoy - HAProxy - NGiNX', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'web-front/Privacy against [friend/pubic/company/attackers]', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Redis/cluster script test for multi datacenter', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'gitlab-runner firewall rules: close public internet', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'static-challange for ovpn-auth', 'initial');
-- INSERT INTO "TASK"("document_id", "content", "task_status") VALUES ('7baf5a59-b6fd-554b-9b4f-d694dc6f6d36', 'Golden image for vpn server', 'initial');