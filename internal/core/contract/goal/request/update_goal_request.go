package request

type UpdateGoalRequest struct{}

type UpdateGoalStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
