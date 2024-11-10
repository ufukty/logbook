package columns

import (
	"regexp"
)

var (
	regexp_base64_url   = regexp.MustCompile(`[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_]+$`) // defined in std lib base64.URLEncoding
	regexp_credit_card  = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`)
	regexp_date         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // FIXME:
	regexp_email        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	regexp_human_name   = regexp.MustCompile(`^\p{L}+([ '-]\p{L}+)*$`)
	regexp_numeric      = regexp.MustCompile(`^[1-9][0-9]*$`)
	regexp_phone_number = regexp.MustCompile(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`)
	regexp_group_title  = regexp.MustCompile(`[\p{L} ]+`)
	regexp_tag_title    = regexp.MustCompile(`[\p{L} ]+`)
	regexp_text         = regexp.MustCompile(`^[\p{L}0-9 ,.?!'’“”-]+$`)
	regexp_url          = regexp.MustCompile(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`)
	regexp_username     = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`)
	regexp_uuid         = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

var (
	max_length_email         = 150
	max_length_group_title   = 100
	max_length_tag_title     = 40
	max_length_human_name    = 100
	max_length_session_token = 256
	max_length_username      = 20
	max_length_uuid          = len("00000000-0000-0000-0000-000000000000")
)

var (
	min_length_email         = 6
	min_length_group_title   = 2
	min_length_tag_title     = 2
	min_length_human_name    = 6
	min_length_session_token = 256
	min_length_username      = 6
	min_length_uuid          = len("00000000-0000-0000-0000-000000000000")
)
