package goat

import "net/mail"

func IsEmailValid(email string) bool {
	emailAddress, err := mail.ParseAddress(email)
	return err == nil && emailAddress.Address == email
}
