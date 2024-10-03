package app

import (
	"fmt"
	"logbook/models/columns"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type GenerateTotpKeyRequest struct {
	User columns.UserId
}

type GenerateTotpKeyResponse struct {
	QrMatrix []string
	Text     string
}

func ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

// booleans take too much characters in JSON
func encode(bitmap [][]bool) []string {
	ss := []string{}
	for i := 0; i < len(bitmap); i++ {
		ss = append(ss, "")
		for j := 0; j < len(bitmap[i]); j++ {
			ss[len(ss)-1] = fmt.Sprintf("%s%s", ss[len(ss)-1], ternary(bitmap[i][j], "X", " "))
		}
	}
	return ss
}

// TODO: store the pending totp secret for validation
func (a *App) GenerateTotpKey(params GenerateTotpKeyRequest) (*GenerateTotpKeyResponse, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Logbook",
		AccountName: "TestUser",
	})
	if err != nil {
		return nil, fmt.Errorf("totp.Generate: %w", err)
	}

	qr, err := qrcode.New(key.String(), qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("qrcode.: %v", err)
	}

	qr.DisableBorder = true
	r := &GenerateTotpKeyResponse{
		QrMatrix: encode(qr.Bitmap()),
		Text:     key.Secret(),
	}

	return r, nil
}

type ValidateTotpKeyRequest struct {
	Nonce string
}

// TODO: persist secret
func (a *App) ValidateTotpKey(params ValidateTotpKeyRequest) error {
	secret := ""

	valid := totp.Validate(params.Nonce, secret)
	if !valid {
		return fmt.Errorf("invalid nonce")
	}

	return nil
}
