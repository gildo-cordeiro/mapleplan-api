package response

type WidgetGoalResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Status        string  `json:"status"`
	DueDate       string  `json:"dueDate"`
	Progress      int     `json:"progress"`
	Phase         string  `json:"phase"`
	Priority      string  `json:"priority"`
	AssignedTo    string  `json:"assignedTo,omitempty"`
	TargetAmount  float64 `json:"targetAmount"`
	CurrentAmount float64 `json:"currentAmount"`
	Description   *string `json:"description,omitempty"`
}
