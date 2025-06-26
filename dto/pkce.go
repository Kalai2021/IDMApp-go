package dto

// PKCE Authorization Request
type PKCEAuthRequest struct {
	ClientID            string `json:"client_id" binding:"required"`
	RedirectURI         string `json:"redirect_uri" binding:"required"`
	Scope               string `json:"scope"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge" binding:"required"`
	CodeChallengeMethod string `json:"code_challenge_method" binding:"required"`
}

// PKCE Authorization Response
type PKCEAuthResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	State            string `json:"state"`
	CodeVerifier     string `json:"code_verifier"`
}

// PKCE Token Exchange Request
type PKCETokenRequest struct {
	GrantType    string `form:"grant_type" json:"grant_type"`
	ClientID     string `form:"client_id" json:"client_id"`
	Code         string `form:"code" json:"code"`
	RedirectURI  string `form:"redirect_uri" json:"redirect_uri"`
	CodeVerifier string `form:"code_verifier" json:"code_verifier"`
	State        string `form:"state" json:"state"`
}

// PKCE Token Response
type PKCETokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// PKCE Error Response
type PKCEErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
}
