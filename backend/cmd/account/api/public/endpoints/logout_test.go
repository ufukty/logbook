package endpoints

import (
	"fmt"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutWithInvalidToken(t *testing.T) {
	tcs := map[string]string{
		"short token":   "",
		"invalid token": "wR5nf0MmDuysxpjGUmfYdhlAg81_OVQtEnRhlNXLgK2Y_EWH-vIFNv_vsxLNTFl-MDnfNJnRm7q-tkkpFPu7bsW7hOK3XshC2-95NPssbkI30rMAsble64FH4H_F1tootwQApVk6HnVK21fH355yqBmqkB8C5VWgCG6lVKmFaaYPrpuu7CEJ_BRml_njyyyehUbn5yUNOuLGxOgGTZHQxLawwlkCh2f6pijYkumVcM80WnpJ-T1ZScdmtb8Qj3T3",
	}
	for tc, token := range tcs {
		t.Run(tc, func(t *testing.T) {
			ep, err := getTestDependencies()
			if err != nil {
				t.Fatal(fmt.Errorf("getting dependencies: %w", err))
			}

			r := httptest.NewRequest("GET", "/", nil)
			r.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: token,
			})
			w := httptest.NewRecorder()

			ep.Logout(receptionist.ZeroRequestId, &middlewares.Store{}, w, r)

			if w.Code == 200 {
				t.Error("Failure expected.")
			}
		})
	}
}
