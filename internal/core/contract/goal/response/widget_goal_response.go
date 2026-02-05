package response

type WidgetGoalResponse struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Status        string  `json:"status"`
	DueDate       string  `json:"due_date"`
	Progress      int     `json:"progress"`
	TargetAmount  float64 `json:"target_amount"`
	CurrentAmount float64 `json:"current_amount"`
	Description   *string `json:"description,omitempty"`
}
