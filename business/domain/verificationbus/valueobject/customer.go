package valueobject

type Customer struct {
	Person          Person
	DateOfBirth     DateOfBirth
	Identifications []Identification
	Email           Email
	BirthCountry    Country
	Address         Address
	Phone           Phone
}
