package request

type UpdateGoalStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type UpdateGoalRequestBody struct {
	Title            string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description      string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status           string `json:"status,omitempty"`
	Phase            string `json:"phase,omitempty"`
	Priority         string `json:"priority,omitempty"`
	DueDate          string `json:"dueDate,omitempty"`
	Progress         int    `json:"progress,omitempty" validate:"omitempty,gte=0,lte=100"`
	AssignedToUser   string `json:"assignedToUser,omitempty"`
	AssignedToCouple string `json:"assignedToCouple,omitempty"`
}
