package response

type GoalResponse struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	Status           string  `json:"status"`
	Phase            string  `json:"phase"`
	Priority         string  `json:"priority"`
	DueDate          string  `json:"dueDate"`
	Progress         int     `json:"progress"`
	AssignedTo       string  `json:"assignedTo,omitempty"`
	AssignedToUser   *string `json:"assignedToUser,omitempty"`
	AssignedToCouple *string `json:"assignedToCouple,omitempty"`
}
