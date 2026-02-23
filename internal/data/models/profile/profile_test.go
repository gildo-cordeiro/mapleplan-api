package profile

import (
	"testing"

	"github.com/gildo-cordeiro/mapleplan-api/internal/data/models"
)

func stringPtr(s string) *string {
	return &s
}

func TestProfileMemberRoleValidation(t *testing.T) {
	tests := []struct {
		role     MemberRole
		expected bool
	}{
		{RolePrimary, true},
		{RoleSpouse, true},
		{MemberRole("invalid"), false},
	}

	for _, test := range tests {
		result := isValidRole(test.role)
		if result != test.expected {
			t.Errorf("isValidRole(%s) = %v, expected %v", test.role, result, test.expected)
		}
	}
}

func TestProfileMemberChangeRole(t *testing.T) {
	member := &ProfileMember{
		ProfileID: "profile-123",
		UserID:    "user-456",
		Role:      RolePrimary,
	}

	err := member.ChangeRole(RoleSpouse)
	if err != nil {
		t.Errorf("ChangeRole failed: %v", err)
	}

	if member.Role != RoleSpouse {
		t.Errorf("Role should be %s, got %s", RoleSpouse, member.Role)
	}
}

func TestProfileMemberChangeRoleInvalid(t *testing.T) {
	member := &ProfileMember{
		ProfileID: "profile-123",
		UserID:    "user-456",
		Role:      RolePrimary,
	}

	err := member.ChangeRole(MemberRole("invalid"))
	if err == nil {
		t.Error("ChangeRole should fail with invalid role")
	}
}

func TestImmigrationProfileBasic(t *testing.T) {
	profile := &ImmigrationProfile{
		Base:   models.Base{ID: "profile-123"},
		UserID: "user-456",
		Name:   "Family Immigration Plan",
	}

	if profile.UserID != "user-456" {
		t.Error("user_id not set correctly")
	}

	if profile.Name != "Family Immigration Plan" {
		t.Error("name not set correctly")
	}
}
