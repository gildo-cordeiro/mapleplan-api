package request

type CreateGoalRequest struct {
	Title            string  `json:"title" validate:"required,min=1,max=255"`
	Description      *string `json:"description,omitempty" validate:"max=1000"`
	Status           string  `json:"status" validate:"required"`
	Phase            string  `json:"phase" validate:"required"`
	Priority         string  `json:"priority" validate:"required"`
	DueDate          string  `json:"dueDate" validate:"required"`
	Progress         int     `json:"progress" validate:"required,gte=0,lte=100"`
	AssignedToUser   *string `json:"assignedToUser,omitempty"`
	AssignedToCouple *string `json:"assignedToCouple,omitempty"`
}
