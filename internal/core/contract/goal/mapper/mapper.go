package mapper

import (
	"strconv"
	"time"

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

func getCorrectAssignedUserID(g *goal.Goal) string {
	if g.CoupleID != nil {
		return g.Couple.Name
	}
	if g.UserId != nil {
		return g.User.FirstName
	}
	return ""
}
