package response

type PartnersListResponse struct {
	Partners []Partner `json:"partners"`
}

type Partner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
