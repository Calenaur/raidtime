package model

import (
	"time"
)

type Event struct {
	ID int 					`json:"id"`		
	Name string 			`json:"name"` 	
	Date time.Time 			`json:"date"` 	
	Color *Color 			`json:"tag"` 	
	Creator *User 			`json:"creator"` 	
	Signups []*Signup 		`json:"signups"` 
}

type Color struct {
	ID int 					`json:"-"`
	Name string 			`json:"-"`
	Color string 			`json:"color"`
}

type Signup struct {
	User *User 				`json:"user"`
	SignupType *SignupType 	`json:"reason"`
	Date time.Time 			`json:"date"`
}

type SignupType struct {
	ID int 					`json:"-"`
	WillAttend bool 		`json:"will_attend"`
	Description string 		`json:"description"`
}
