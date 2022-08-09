package parameters

import "net/http"

func EndpointTaskCreate(w http.ResponseWriter, r *http.Request) {
	parameters := TaskCreate{}
	if err := parameters.Sanitize(r); err != nil {
		w
		// return failure
	}

	// check auth

	// check existence of super task

	// check permissions between task and user
	parameters.Request.UserId
}
