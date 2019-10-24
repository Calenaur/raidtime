package config

import (
	"os"
	"io/ioutil"
    "encoding/json"
)

type Config struct {
	Discord *Discord 		`json:"discord"`
	Database *Database 		`json:"database"`
	Session *Session 		`json:"session"`
}

type Discord struct {
	ClientID string 		`json:"client_id"`
	ClientSecret string 	`json:"client_secret"`
	GrantType string 		`json:"grant_type"`
	Scope string			`json:"scope"` 				
	APIUri string 			`json:"api_uri"`
	UserUri string 			`json:"user_uri"`
	RedirectUri string 		`json:"redirect_uri"`
	AuthorizeUri string 	`json:"authorize_uri"`
	TokenUri string 		`json:"token_uri"`
	TokenRevokeUri string 	`json:"token_revoke_uri"`
}

type Database struct {
	Database string 		`json:"database"`
	Username string 		`json:"username"`
	Password string 		`json:"password"`
}

type Session struct {
	SessionDuration int 	`json:"session_duration"`
	SessionSecret string 	`json:"session_secret"`
}

func Load(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	var cfg Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	d := cfg.Discord
	d.UserUri = d.APIUri + d.UserUri
	d.RedirectUri = d.APIUri + d.RedirectUri
	d.AuthorizeUri = d.APIUri + d.AuthorizeUri
	d.TokenUri = d.APIUri + d.TokenUri
	d.TokenRevokeUri = d.APIUri + d.TokenRevokeUri
	
	return &cfg, nil
}
