package endpoints

import (
	"logbook/internal/web/validate"
	"regexp"
)

type ObjectiveContent string

var (
	regexp_basic_test = regexp.MustCompile(`^[\p{L} .,-_]*$`)
)

var (
	max_length_objective_content = len("Ensure that the input conforms to your business logic. For example, check if a user is allowed to perform a certain action, if the data respects certain business rules, or if relationships between data entities are maintained.")
)

var (
	min_length_objective_content = 1
)

func (v ObjectiveContent) Validate() error {
	return validate.StringBasics(string(v), min_length_objective_content, max_length_objective_content, regexp_basic_test)
}
