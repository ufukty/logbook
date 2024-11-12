package cookies

import (
	"fmt"
	"logbook/internal/average"
	"logbook/models/columns"
	"net/http"
	"time"
)

var ErrNoSessionToken = fmt.Errorf("session_token not found")

const sessionTokenKey = "session_token"

func GetSessionToken(r *http.Request) (columns.SessionToken, error) {
	st, err := r.Cookie(sessionTokenKey)
	if err != nil {
		return "", fmt.Errorf("no session_token found")
	}
	return columns.SessionToken(st.Value), nil
}

func SetSessionToken(w http.ResponseWriter, st columns.SessionToken, sessionstart time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionTokenKey,
		Value:    string(st),
		Expires:  sessionstart.Add(average.Week),
		HttpOnly: true,
		Secure:   true,
	})
}

func ExpireSessionToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   sessionTokenKey,
		MaxAge: -1,
	})
}
