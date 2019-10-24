package store

import (
	"time"
	"strings"
	"strconv"
	"errors"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/calenaur/raidtime/model"
	"github.com/calenaur/raidtime/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DiscordStore struct {
	cfg *config.Config
}

func NewDiscordStore(cfg *config.Config) *DiscordStore {
	return &DiscordStore{
		cfg: cfg,
	}
}

func (ds *DiscordStore) GetCredentialsByCode(code string) (*model.UserCredentials, error) {
	token, err := ds.fetchAccessToken(code)
	if err != nil {
		return nil, err
	}

	credentials, err := ds.fetchCredentials(token)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (ds *DiscordStore) fetchAccessToken(code string) (*model.AccessToken, error) {
	data := url.Values{}
	data.Set("client_id", ds.cfg.Discord.ClientID)
	data.Set("client_secret", ds.cfg.Discord.ClientSecret)
	data.Set("grant_type", ds.cfg.Discord.GrantType)
	data.Set("code", code)
	data.Set("redirect_uri", "http://127.0.0.1:1323/auth")
	data.Set("scope", ds.cfg.Discord.Scope)

    req, _ := http.NewRequest("POST", ds.cfg.Discord.TokenUri, strings.NewReader(data.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
    	return nil, errors.New(string(body))
    }

    a := new(model.AccessToken)
    json.Unmarshal(body, a)
    a.CreationTime = time.Now().Unix()
    return a, nil
}

func (ds *DiscordStore) fetchCredentials(token *model.AccessToken) (*model.UserCredentials, error) {
    req, _ := http.NewRequest("GET", ds.cfg.Discord.UserUri, nil)
    req.Header.Add("Authorization", "Bearer " + token.AccessToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bytes, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
    	return nil, errors.New(string(bytes))
    }

    credentials := new(model.UserCredentials)
    err = json.Unmarshal(bytes, credentials)
	if err != nil {
		return nil, err
	}

    return credentials, nil
}