package dto

import (
	"github.com/asaskevich/govalidator"
	"github.com/ilyushkaaa/Filmoteka/pkg/validator"
)

type (
	AuthRequest struct {
		Password string `json:"password" valid:"required,length(8|255)"`
		Username string `json:"username" valid:"required,matches(^[a-zA-Z0-9_]+$)"`
	}
	AuthResponse struct {
		SessionID string `json:"session_id"`
	}
)

func (authReqDTO *AuthRequest) Validate() []string {
	_, err := govalidator.ValidateStruct(authReqDTO)
	return validator.CollectErrors(err)
}
