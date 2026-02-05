package request

type UpdateUserOnboardingRequest struct {
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	ImmigrationGoal string `json:"immigrationGoal" binding:"required"`
	PartnerEmail    string `json:"partnerEmail"`
}
