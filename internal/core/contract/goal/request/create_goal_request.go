package request

type CreateGoalRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	DueDate     string `json:"due_date" binding:"required,datetime=2006-01-02"`
}
