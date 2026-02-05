package response

type UserWithCoupleResponse struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	FirstName            *string `json:"FirstName,omitempty"`
	LastName             *string `json:"lastName,omitempty"`
	Phone                *string `json:"phone,omitempty"`
	CoupleDateOfBirth    *string `json:"coupleDateOfBirth,omitempty"`
	PartnerEmail         *string `json:"partnerEmail,omitempty"`
	PartnerNameFirstName *string `json:"partnerNameFirstName,omitempty"`
	PartnerNameLastName  *string `json:"partnerNameLastName,omitempty"`
	Province             *string `json:"province,omitempty"`
	CreatedAt            *string `json:"createdAt,omitempty"`
}
