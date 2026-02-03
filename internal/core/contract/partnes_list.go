package contract

type PartnersListDto struct {
	Partners []Partner `json:"partners"`
}

type Partner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
