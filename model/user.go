package model

import (
	"time"
	"strconv"
	"encoding/hex"
	"crypto/sha512"
)

type User struct {
	ID int64
	Username string
	Discriminator int
	Avatar string
	GuildRank string
	Class *Class
	Permissions *Permissions
	Session string
	SessionCreationTime time.Time
}

type Class struct {
	ID int
	Name string
	Color string
}

type Permissions struct {
	ID int
	Name string
	ManageUsers bool
	ManageEvents bool
}

func (user *User) GenerateSession(secret string) string {
	timestamp := time.Now().Unix()
	id := ^user.ID
	t := ^timestamp
	sha := sha512.New()
	sha.Write([]byte(strconv.FormatInt(id, 9) + secret + strconv.FormatInt(t, 8)))
	return hex.EncodeToString(sha.Sum(nil))
}