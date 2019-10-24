package model

type AccessToken struct {
	AccessToken string 			`json:"access_token"`
	Scope string 				`json:"scope"`
	TokenType string 			`json:"token_type"`
	ExpiresIn int  				`json:"expires_in"`
	RefreshToken string 		`json:"refresh_token"`
	CreationTime int64 			`json:"creation_time"`
}

type UserCredentials struct {
	ID string 					`json:"id"`
	Username string 			`json:"username"`
	Locale string 				`json:"locale"`
	Has2FA bool 				`json:"mfa_enabled"`
	Flags int 					`json:"flags"`
	Avatar string 				`json:"avatar"`
	Discriminator string 		`json:"discriminator"`
}