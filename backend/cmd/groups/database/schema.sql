CREATE DOMAIN "GroupId" AS uuid;

CREATE DOMAIN "GroupMembershipId" AS uuid;

CREATE DOMAIN "GroupInviteId" AS uuid;

CREATE DOMAIN "UserId" AS uuid;

CREATE TABLE "group"(
    "gid" "GroupId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "name" text NOT NULL,
    "creator" "UserId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    PRIMARY KEY ("gid")
);

CREATE TABLE "group_member_user"(
    "gmid" "GroupMembershipId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "gid" "GroupId",
    "member" "UserId",
    "ginvid" "GroupInviteId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TABLE "group_member_group"(
    "gmid" "GroupMembershipId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "gid" "GroupId",
    "member" "GroupId",
    "ginvid" "GroupInviteId" NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TYPE "GroupInviteStatus" AS ENUM(
    'sent',
    'accepted',
    'rejected'
);

CREATE TABLE "group_invite_user"(
    "ginvid" "GroupInviteId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "inviter" "UserId" NOT NULL,
    "invitee" "UserId" NOT NULL,
    "status" "GroupInviteStatus" NOT NULL DEFAULT 'sent',
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "group_invite_group"(
    "ginvid" "GroupInviteId" NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    "inviter" "UserId" NOT NULL,
    "invitee" "GroupId" NOT NULL,
    "status" "GroupInviteStatus" NOT NULL DEFAULT 'sent',
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

