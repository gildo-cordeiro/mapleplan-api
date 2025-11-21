package user

type Repository interface {
	FindByEmail(email string) (*User, error)
	Save(user *User) (string, error)
}
