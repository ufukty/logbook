package transports

import (
	"regexp"
	"time"
)

var (
	regexp_base64_url = regexp.MustCompile(`[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_]+$`) // defined in std lib [base64.URLEncoding]
	regexp_uuid       = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

var (
	length_anti_csrf_token = 32 // 256 bits
	length_uuid            = len("00000000-0000-0000-0000-000000000000")
)

var (
	min_human_birthday  = time.Now().AddDate(-100, 0, 0) // 100 years ago
	min_length_password = 6
)

var (
	max_human_birthday  = time.Now().AddDate(-18, 0, 1) // 18 years ago
	max_length_password = 2048
)
