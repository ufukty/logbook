package controllers

import "net/http"

const (
	OFFSET_MIN = 0
	OFFSET_MAX = 10000000
	LIMIT_MIN  = 1
	LIMIT_MAX  = 100000
)

func PlacementArrayHierarchical(w http.ResponseWriter, r *http.Request) {
	// params := parameters.PlacementArrayHierarchical{}

	// offset, err := strconv.Atoi(vars["offset"])
	// if err != nil {
	// 	return []error{err, c.ErrOffsetInputIsNotValidInteger}
	// }
	// if offset < OFFSET_MIN || OFFSET_MAX < offset {
	// 	return []error{c.ErrOffsetInputIsNotInAllowedRange}
	// }

	// limit, err := strconv.Atoi(vars["limit"])
	// if err != nil {
	// 	return []error{err, c.ErrLimitInputIsNotValidInteger}
	// }
	// if limit < LIMIT_MIN || LIMIT_MAX < limit {
	// 	return []error{c.ErrLimitInputIsNotInAllowedRange}
	// }
}
