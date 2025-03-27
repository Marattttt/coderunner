package models

type AuthResponce struct {
	AccessTok string `json:"access"`
	RefreshTok string `json:"refresh"`
}
