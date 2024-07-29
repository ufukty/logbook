package columns

import (
	"logbook/internal/web/validate"
	"regexp"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	NonNegativeNumber int
)

var (
	regexp_uuid       = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	regexp_text       = regexp.MustCompile(`^[\p{L}0-9 ,.?!'’“”-]+$`)
	regexp_url        = regexp.MustCompile(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`)
	regexp_date       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // FIXME:
	regexp_numeric    = regexp.MustCompile(`^[1-9][0-9]*$`)
	regexp_base64_url = regexp.MustCompile(`[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_]+$`) // defined in std lib base64.URLEncoding
)

var (
	max_length_uuid = len("00000000-0000-0000-0000-000000000000")
)

var (
	min_length_uuid = len("00000000-0000-0000-0000-000000000000")
)

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

var ZeroTimestamp = pgtype.Timestamp{}
