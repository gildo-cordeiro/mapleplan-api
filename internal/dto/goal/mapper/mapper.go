package mapper

import (
	"fmt"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/goal/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models/goal"
)

func ToWidgetGoalResponse(g *goal.Goal) response.WidgetGoalResponse {
	return response.WidgetGoalResponse{
		ID:          g.ID,
		Title:       g.Name,
		Status:      goal.StatusToString(g.Status),
		DueDate:     formatDate(g.DueDate),
		Progress:    g.Progress,
		Phase:       goal.PhaseToString(g.Phase),
		Priority:    goal.PriorityToString(g.Priority),
		AssignedTo:  getCorrectAssignedUserID(g),
		Description: g.Description,
	}
}
func ToGoalResponse(g *goal.Goal) response.GoalResponse {
	return response.GoalResponse{
		ID:          g.ID,
		Title:       g.Name,
		Description: g.Description,
		DueDate:     formatDate(g.DueDate),
		Status:      goal.StatusToString(g.Status),
		Phase:       goal.PhaseToString(g.Phase),
		Priority:    goal.PriorityToString(g.Priority),
		AssignedTo:  getCorrectAssignedUserID(g),
	}
}

func ToGoalUpdateResponse(g *goal.Goal) response.GoalResponse {
	return response.GoalResponse{
		ID:               g.ID,
		Title:            g.Name,
		Description:      g.Description,
		DueDate:          formatDate(g.DueDate),
		Status:           goal.StatusToString(g.Status),
		Phase:            goal.PhaseToString(g.Phase),
		Priority:         goal.PriorityToString(g.Priority),
		AssignedToUser:   g.UserId,
		AssignedToCouple: g.CoupleID,
	}
}

func CreateGoalRequestToGoalDomain(c *request.CreateGoalRequest) *goal.Goal {
	dueDate, _ := toDate(c.DueDate)

	status, _ := goal.StringToStatus(c.Status)
	phase, _ := goal.StringToPhase(c.Phase)
	priority, _ := goal.StringToPriority(c.Priority)

	return &goal.Goal{
		Name:        c.Title,
		Description: c.Description,
		DueDate:     dueDate,
		Status:      status,
		Phase:       phase,
		Priority:    priority,
	}
}

func UpdateGoalRequestToGoalDomain(u *request.UpdateGoalRequestBody) *goal.Goal {
	dueDate, _ := toDate(u.DueDate)

	status, _ := goal.StringToStatus(u.Status)
	phase, _ := goal.StringToPhase(u.Phase)
	priority, _ := goal.StringToPriority(u.Priority)

	// Only set pointer fields if incoming value is non-empty to avoid writing empty string as UUID
	var userID *string
	if u.AssignedToUser != "" {
		userID = &u.AssignedToUser
	}
	var coupleID *string
	if u.AssignedToCouple != "" {
		coupleID = &u.AssignedToCouple
	}

	return &goal.Goal{
		Name:        u.Title,
		Description: &u.Description,
		DueDate:     dueDate,
		Status:      status,
		Phase:       phase,
		Priority:    priority,
		UserId:      userID,
		CoupleID:    coupleID,
	}
}

func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(time.RFC3339)
}

func toDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

func getCorrectAssignedUserID(g *goal.Goal) string {
	// Be defensive: check both the ID pointer and the loaded relation to avoid nil dereference
	if g.CoupleID != nil && g.Couple != nil {
		return g.Couple.Name
	}
	if g.UserId != nil && g.User != nil {
		return g.User.FirstName
	}
	return ""
}
