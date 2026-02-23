package mapper

import (
	"testing"

	"github.com/gildo-cordeiro/mapleplan-api/internal/dto/goal/request"
)

func TestUpdateGoalRequestToGoalDomain_ShouldNotCreatePointerWhenEmpty(t *testing.T) {
	req := request.UpdateGoalRequestBody{
		Title:             "t",
		AssignedToUser:    "",
		AssignedToProfile: "",
	}
	g := UpdateGoalRequestToGoalDomain(&req)
	if g.UserID != nil {
		t.Fatalf("expected UserID to be nil when AssignedToUser is empty, got: %v", *g.UserID)
	}
	if g.ImmigrationProfileID != nil {
		t.Fatalf("expected ProfileID to be nil when AssignedToProfile is empty, got: %v", *g.ImmigrationProfileID)
	}
}
