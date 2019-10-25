package model

import (
	"time"
	"strconv"
	"encoding/hex"
	"crypto/sha512"
)

type User struct {
	ID int64 					`json:"id"`
	Username string 			`json:"username"`
	Discriminator int 			`json:"-"`
	Avatar string 				`json:"-"`
	Class *Class 				`json:"class"`
	GuildRank *GuildRank 		`json:"guild_rank"`
	Permissions *Permissions 	`json:"-"`
}

type Class struct {
	ID int 						`json:"-"`
	Name string 				`json:"name"`
	Color string 				`json:"color"`
}

type Permissions struct {
	ID int 						`json:"-"`
	Name string 				`json:"-"`
	ManageUsers bool 			`json:"-"`
	ManageEvents bool 			`json:"-"`
}

type GuildRank struct {
	ID int 						`json:"-"`
	Name string 				`json:"name"`
	Protected bool 				`json:"-"`
}

func (user *User) GenerateSession(secret string) string {
	timestamp := time.Now().Unix()
	id := ^user.ID
	t := ^timestamp
	sha := sha512.New()
	sha.Write([]byte(strconv.FormatInt(id, 9) + secret + strconv.FormatInt(t, 8)))
	return hex.EncodeToString(sha.Sum(nil))
}