package model

type User struct {
	ID             int
	Username       string
	Password, Salt []byte
	Balance        int
}
