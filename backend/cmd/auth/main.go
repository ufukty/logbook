package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5/middleware"
)

// TODO: Anti-CSRF token
// TODO: Check user input against XSS
// TODO: Measure against Clickjacking concern [RFC6749:OAuth 2.0/10.13]: w.Header().Set("X-Frame-Options", "DENY")
// TODO: Check if [http.FileServer] sets Timeouts against Slowloris attack
// TODO: CORS
func Main() error {
	crt := filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/auth/certificates/localhost.crt") // TODO: read flags
	key := filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/auth/certificates/localhost.key") // TODO: read flags
	build := filepath.Join(os.Getenv("WORKSPACE"), "backend/cmd/auth/web/build")

	err := http.ListenAndServeTLS(":8082", crt, key, middleware.Logger(http.FileServer(http.Dir(build))))
	if err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Println(err)
	}
}
