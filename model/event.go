package model

import (
	"time"
)

type Event struct {
	ID int 			
	Name string 	
	Date time.Time 	
	Color *Color 	
	Creator *User 	
	Signups []*Signup 
}

type Color struct {
	ID int
	Name string
	Color string
}

type Signup struct {
	Event *Event
	User *User
	SignupType *SignupType
	Date time.Time
}

type SignupType struct {
	ID int
	WillAttend bool
	Description string
}
