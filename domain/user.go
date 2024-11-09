package domain

type User struct {
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	CreateUser(user User) (*User, error)
}
