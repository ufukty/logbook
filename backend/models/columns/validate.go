package columns

import (
	"logbook/internal/web/validate"
)

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

func (v AccessId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v BookmarkId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v BottomUpPropsId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v CollaborationId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v CollaboratorId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v ControlAreaId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v DelegationId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v Email) Validate() error {
	return validate.StringBasics(string(v), min_length_email, max_length_email, regexp_email)
}

func (v GroupId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v GroupInviteId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v GroupMembershipId) Validate() error {
	return validate.StringBasics(string(v), min_length_group_title, max_length_group_title, regexp_group_title)
}

func (v GroupName) Validate() error {
	return validate.StringBasics(string(v), min_length_group_title, max_length_group_title, regexp_group_title)
}

func (v HumanName) Validate() error {
	return validate.StringBasics(string(v), min_length_human_name, max_length_human_name, regexp_human_name)
}

func (v LinkId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v LoginId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v ObjectiveContent) Validate() error {
	return validate.StringBasics(string(v), min_length_objective_content, max_length_objective_content, regexp_objective_content)
}

func (v ObjectiveId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v OperationId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v Phone) Validate() error {
	return validate.StringBasics(string(v), min_length_phone_number, max_length_phone_number, regexp_phone_number)
}

func (v PropertiesId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v SessionId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v SessionToken) Validate() error {
	return validate.StringBasics(string(v), min_length_session_token, max_length_session_token, regexp_base64_url)
}

func (v TagId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v TagTitle) Validate() error {
	return validate.StringBasics(string(v), min_length_tag_title, max_length_tag_title, regexp_tag_title)
}

func (v UserAgent) Validate() error {
	return validate.StringBasics(string(v), min_length_user_agent, max_length_user_agent, nil)
}

func (v UserId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v Username) Validate() error {
	return validate.StringBasics(string(v), min_length_username, max_length_username, regexp_username)
}

func (v VersionId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}
