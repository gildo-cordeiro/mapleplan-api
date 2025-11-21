package contract

type CreateNewUserDto struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	LastName string `json:"last_name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
