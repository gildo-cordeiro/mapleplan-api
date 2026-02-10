package repository

import (
	"testing"

	"github.com/gildo-cordeiro/mapleplan-api/internal/core/domain/goal"
)

func TestSanitizeUUIDPointers(t *testing.T) {
	empty := ""
	g := &goal.Goal{
		UserId:   &empty,
		CoupleID: &empty,
	}

	sanitizeUUIDPointers(g)

	if g.UserId != nil {
		t.Fatalf("expected UserId to be nil after sanitization, got: %v", *g.UserId)
	}
	if g.CoupleID != nil {
		t.Fatalf("expected CoupleID to be nil after sanitization, got: %v", *g.CoupleID)
	}
}
