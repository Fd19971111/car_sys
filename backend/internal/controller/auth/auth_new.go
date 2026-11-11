package auth

import (
	"engineServ/api/auth"
	authLogic "engineServ/internal/logic/api/auth"
)

type ControllerV1 struct {
	auth *authLogic.Auth
}

func NewV1() auth.IAuthV1 {
	return &ControllerV1{auth: authLogic.New()}
}

type ControllerV1Plus struct {
	auth *authLogic.Auth
}

func NewV1Plus() auth.IAuthV1Plus {
	return &ControllerV1Plus{auth: authLogic.New()}
}
