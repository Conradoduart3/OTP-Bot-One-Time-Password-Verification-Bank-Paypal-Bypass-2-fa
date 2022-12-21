package entity

import "strconv"

type ResponseData struct {
	Status   int     `json:"status"`
	Response string  `json:"Response"`
	Origins  []Origs `json:"origins,omitempty"`
	Users    []User  `json:"users,omitempty"`
}

type Payload struct {
	Email  string `json:"email"`
	ChatID string `json:"chat-id"`
	Origin string `json:"origin"`
}

type Origs struct {
	Origin string `json:"origin,omitempty"`
}

type User struct {
	User string `json:"user,omitempty"`
}

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}
