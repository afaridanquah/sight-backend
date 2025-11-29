package elasticsearch

type Customer struct {
	ID              string         `json:"id"`
	FirstName       string         `json:"first_name"`
	MiddleName      string         `json:"middle_name"`
	LastName        string         `json:"last_name"`
	DateOfBirth     string         `json:"date_of_birth"`
	Email           string         `json:"email"`
	PhoneNumber     string         `json:"phone_number"`
	BirthCountry    string         `json:"birth_country"`
	Identifications map[string]any `json:"identifications"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}
