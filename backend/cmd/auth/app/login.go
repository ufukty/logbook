package app

import "io"

// TODO: embed anti-CSRF token which is also stored in cache
func (a *App) GetLoginPage() (io.Reader, error) {
	return nil, nil
}
