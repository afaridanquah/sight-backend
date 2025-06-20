package es

type CustomerAdded struct {
	Person       Person  `json:"person"`
	BirthCountry Country `json:"birth_country"`
	// DateOfBirth     DateOfBirth      `json:"date_of_birth"`
	// Email           Email            `json:"email"`
	// Address         Address          `json:"address"`
	// Identifications []Identification `json:"identifications"`
}

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	OtherNames string `json:"other_names"`
}

type Country struct {
	AlphaCode string `json:"alpha_code"`
	Name      string `json:"name"`
}
