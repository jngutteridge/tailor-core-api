package core

import (
	"fmt"
	_ "github.com/lib/pq"
)

type ErrorMessage struct {
	ErrorMessage string
}

type User struct {
	ID               string `db:"id"`
	Provider         string `db:"provider"`
	ProviderID       string `db:"provider_id"`
	Email            string `db:"email"`
	Name             string `db:"name"`
	ImageURL         string `db:"image_url"`
	VerificationCode string `db:"verification_code"`
}

type RegistrationInfoResponse struct {
	ErrorMessage
	User
}

type RegistrationResponse struct {
	ErrorMessage
	VerificationCode string
}
type UserIndexResponse struct {
	Users []User
}

func GetImageURL(provider string, clientId string) string {
	switch provider {
	case "facebook":
		return fmt.Sprintf("https://graph.facebook.com/%s/picture?type=large&width=720&height=720", clientId)
	default:
		return ""
	}
}
