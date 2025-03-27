package models

type OAuthRedirect struct {
	LoginURL string `json:"loginURL"`
	TokenAccessCode    string `json:"tokenAccessCode"`
}
