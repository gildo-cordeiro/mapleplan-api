package finance

import (
	"testing"
	"time"
)

func stringPtr(s string) *string {
	return &s
}

func TestFinanceIsPrivate(t *testing.T) {
	finance := &Finance{
		UserID:               stringPtr("user-123"),
		ImmigrationProfileID: nil,
		Amount:               100.0,
	}

	if !finance.IsPrivate() {
		t.Error("finance should be private when profile_id is nil")
	}
}

func TestFinanceShareWithProfile(t *testing.T) {
	finance := &Finance{
		UserID:               stringPtr("user-123"),
		ImmigrationProfileID: nil,
		Amount:               100.0,
	}

	profileID := "profile-456"
	finance.ShareWithProfile(profileID)

	if finance.IsPrivate() {
		t.Error("finance should not be private after sharing")
	}

	if finance.ImmigrationProfileID == nil || *finance.ImmigrationProfileID != profileID {
		t.Errorf("profile_id should be %s", profileID)
	}
}

func TestFinanceMakePrivate(t *testing.T) {
	profileID := "profile-456"
	finance := &Finance{
		UserID:               stringPtr("user-123"),
		ImmigrationProfileID: &profileID,
		Amount:               100.0,
	}

	finance.MakePrivate()

	if !finance.IsPrivate() {
		t.Error("finance should be private after calling MakePrivate")
	}
}

func TestFinanceTypesValidation(t *testing.T) {
	tests := []struct {
		financeType FinanceType
		expected    bool
	}{
		{Income, true},
		{Expense, true},
		{Transfer, true},
		{FinanceType("invalid"), false},
	}

	for _, test := range tests {
		result := IsValidType(test.financeType)
		if result != test.expected {
			t.Errorf("IsValidType(%s) = %v, expected %v", test.financeType, result, test.expected)
		}
	}
}

func TestFinanceCategoryValidation(t *testing.T) {
	tests := []struct {
		category Category
		expected bool
	}{
		{Salary, true},
		{Housing, true},
		{Food, true},
		{Other, true},
		{Category("invalid"), false},
	}

	for _, test := range tests {
		result := IsValidCategory(test.category)
		if result != test.expected {
			t.Errorf("IsValidCategory(%s) = %v, expected %v", test.category, result, test.expected)
		}
	}
}

func TestFinanceBasicFields(t *testing.T) {
	now := time.Now()
	finance := &Finance{
		UserID:   stringPtr("user-123"),
		Category: Salary,
		Type:     Income,
		Amount:   5000.00,
		Currency: "CAD",
		Date:     now,
	}

	if finance.UserID == nil || *finance.UserID != "user-123" {
		t.Error("user_id not set correctly")
	}

	if finance.Amount != 5000.00 {
		t.Error("amount not set correctly")
	}

	if finance.Currency != "CAD" {
		t.Error("currency not set correctly")
	}
}
