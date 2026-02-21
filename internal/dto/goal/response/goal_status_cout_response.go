package response

type GoalStatusCountResponse struct {
	Total      int `json:"total"`
	NotStarted int `json:"notStarted"`
	InProgress int `json:"inProgress"`
	Completed  int `json:"completed"`
}
