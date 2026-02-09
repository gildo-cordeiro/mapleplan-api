package mapper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/request"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/contract/goal/response"
	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
	"github.com/shopspring/decimal"
)

func ToWidgetGoalResponse(g *goal.Goal) response.WidgetGoalResponse {
	return response.WidgetGoalResponse{
		ID:            g.ID,
		Title:         g.Name,
		Status:        goal.StatusToString(g.Status),
		DueDate:       formatDate(g.DueDate),
		Progress:      g.Progress,
		TargetAmount:  decimalToFloat64(g.TargetAmount),
		CurrentAmount: decimalToFloat64(g.CurrentAmount),
		Phase:         goal.PhaseToString(g.Phase),
		Priority:      goal.PriorityToString(g.Priority),
		AssignedTo:    getCorrectAssignedUserID(g),
		Description:   g.Description,
	}
}
func ToGoalResponse(g *goal.Goal) response.GoalResponse {
	return response.GoalResponse{
		ID:            g.ID,
		Title:         g.Name,
		Description:   g.Description,
		DueDate:       formatDate(g.DueDate),
		TargetAmount:  decimalToFloat64(g.TargetAmount),
		CurrentAmount: decimalToFloat64(g.CurrentAmount),
		Status:        goal.StatusToString(g.Status),
		Phase:         goal.PhaseToString(g.Phase),
		Priority:      goal.PriorityToString(g.Priority),
		AssignedTo:    getCorrectAssignedUserID(g),
	}
}

func ToGoalDomain(c *request.CreateGoalRequest) *goal.Goal {
	dueDate, _ := toDate(c.DueDate)

	var current decimal.Decimal
	if c.CurrentAmount != nil {
		current = decimal.NewFromFloat(*c.CurrentAmount)
	} else {
		current = decimal.Zero
	}

	status, _ := goal.StringToStatus(c.Status)
	phase, _ := goal.StringToPhase(c.Phase)
	priority, _ := goal.StringToPriority(c.Priority)

	return &goal.Goal{
		Name:          c.Title,
		Description:   c.Description,
		DueDate:       dueDate,
		TargetAmount:  decimal.NewFromFloat(c.TargetAmount),
		CurrentAmount: current,
		Status:        status,
		Phase:         phase,
		Priority:      priority,
	}
}

func decimalToFloat64(d decimal.Decimal) float64 {
	if d.IsZero() {
		return 0
	}
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0
	}
	return f
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
	if g.CoupleID != nil {
		return g.Couple.Name
	}
	if g.UserId != nil {
		return g.User.FirstName
	}
	return ""
}
