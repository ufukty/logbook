package serialize

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidationIssues(w http.ResponseWriter, issues map[string]any) error {
	w.WriteHeader(400)
	err := json.NewEncoder(w).Encode(issues)
	if err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}
	return nil
}
