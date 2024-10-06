package app

type AuthRequest struct {
	ResponseType string `in:"query=response_type"` // req. [ code [ . | id_token | token | id_token token ] | id_token [ . | token ] ]
	ClientId     string `in:"query=client_id"`     // req. (returned after client registration)
	RedirectURI  string `in:"query=redirect_uri"`  // opt. (should be a URI that is registered at client registration)
	Scope        string `in:"query=scope"`         // req. (no default value. server will reject invalid/empty scope)
	State        string `in:"query=state"`         // recommended against CSRF attacks
	// below is usable when "scope" has "openid"
	Nonce       string `in:"query=nonce"`         // opt.
	Display     string `in:"query=display"`       // opt. (no use)
	Prompt      string `in:"query=prompt"`        // opt. (no use)
	MaxAge      string `in:"query=max_age"`       // opt. (no use, will be decided by OP)
	UILocales   string `in:"query=ui_locales"`    // opt. (no use)
	IDTokenHint string `in:"query=id_token_hint"` // opt. (no use)
	LoginHint   string `in:"query=login_hint"`    // opt. (no use)
	ACRValues   string `in:"query=acr_values"`    // opt. (to do)
}

type AuthResponse struct {
	Code             string `in:"query=code"`              // req. when success
	State            string `in:"query=state"`             // req. always
	Error            string `in:"query=error"`             // req. when failure
	ErrorDescription string `in:"query=error_description"` // opt.
	ErrorURI         string `in:"query=error_uri"`         // opt.
}

type TokenRequest struct {
	GrantType   string `json:""`
	Code        string `json:""`
	RedirectURI string `json:""`
	ClientId    string `json:""`
}

type TokenResponse struct {
	AccessToken  string `json:""`
	TokenType    string `json:""`
	ExpiresIn    string `json:""`
	RefreshToken string `json:""`
}
