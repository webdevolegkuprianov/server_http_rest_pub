package model

import "time"

//User
type User1 struct {
	ID       int
	Email    string `json:"email"`
	Password string `json:"password"`
}

//for jwt verify
type User2 struct {
	UserID uint64
}

//for token and exp
type Token_exp struct {
	Token string
	Exp   time.Time
}

type AccessDetails struct {
	UserId uint64
	Exp    uint64
}

//response struct
type Response struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

//response struct booking
type ResponseBooking struct {
	StatusMs       string `json:"status_ms"`
	ResponseMs     string `json:"response_ms"`
	StatusGazCrm   string `json:"status_gcrm"`
	ResponseGazCrm string `json:"response_gcrm"`
}
