package mapper

import (
	"testing"

	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/goal/request"
)

func TestUpdateGoalRequestToGoalDomain_ShouldNotCreatePointerWhenEmpty(t *testing.T) {
	req := request.UpdateGoalRequestBody{
		Title:            "t",
		AssignedToUser:   "",
		AssignedToCouple: "",
	}
	g := UpdateGoalRequestToGoalDomain(&req)
	if g.UserId != nil {
		t.Fatalf("expected UserId to be nil when AssignedToUser is empty, got: %v", *g.UserId)
	}
	if g.CoupleID != nil {
		t.Fatalf("expected CoupleID to be nil when AssignedToCouple is empty, got: %v", *g.CoupleID)
	}
}
