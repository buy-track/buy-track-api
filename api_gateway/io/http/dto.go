package http

import "my-stocks/domains"

type UserRegisterDto struct {
	Email    string
	Password string
	Name     string
}

type UserLoginDto struct {
	Email    string
	Password string
}

type UserRegisterResponseDto struct {
	Token *domains.Token `json:"token"`
	User  *domains.User  `json:"user"`
}

type UserLoginResponseDto struct {
	Token *domains.Token `json:"token"`
	User  *domains.User  `json:"user"`
}
