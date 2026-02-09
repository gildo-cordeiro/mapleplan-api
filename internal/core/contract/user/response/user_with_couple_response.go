package response

type UserWithCoupleResponse struct {
	ID                string  `json:"id"`
	Email             string  `json:"email"`
	FirstName         *string `json:"firstName,omitempty"`
	LastName          *string `json:"lastName,omitempty"`
	Phone             *string `json:"phone,omitempty"`
	CoupleDateOfBirth *string `json:"coupleDateOfBirth,omitempty"`
	CoupleID          *string `json:"coupleId,omitempty"`
	PartnerId         *string `json:"partnerId,omitempty"`
	PartnerEmail      *string `json:"partnerEmail,omitempty"`
	PartnerFirstName  *string `json:"partnerFirstName,omitempty"`
	PartnerLastName   *string `json:"partnerLastName,omitempty"`
	Province          *string `json:"province,omitempty"`
	CreatedAt         *string `json:"createdAt,omitempty"`
}
