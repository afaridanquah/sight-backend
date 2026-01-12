package identificationbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/identificationbus/valueobject"
)

type Identification struct {
	ID                 valueobject.ID
	CustomerID         string
	FirstName          string
	LastName           string
	MiddleName         string
	Sex                string
	Pin                string
	PlaceOfBirth       string
	DateOfBirth        time.Time
	Nationality        valueobject.Country
	IdentificationType valueobject.IdentificationType
	IssuedDate         time.Time
	ExpDate            time.Time
	StateOrRegion      string
	City               string
	Country            valueobject.Country
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NewIdentification struct {
	CustomerID         string
	FirstName          string
	LastName           string
	MiddleName         string
	Sex                string
	Pin                string
	PlaceOfBirth       string
	DateOfBirth        time.Time
	Nationality        valueobject.Country
	IssuedDate         time.Time
	ExpDate            time.Time
	StateOrRegion      string
	City               string
	Country            valueobject.Country
	IdentificationType valueobject.IdentificationType
}

// func (napp *NewIdentification) Validate() error {
// 	switch napp.IdentificationType.String() {
// 	case "PASSPORT":
// 		return RequirementForPassport(*napp)
// 	default:
// 		return errors.New("invalid")
// 	}
// }

// func RequirementForPassport(napp NewIdentification) error {
// 	var err error
// 	if napp.Pin == "" {
// 		errors.Join(err, errors.New("pin is required"))
// 	}

// 	if napp.PlaceOfBirth == "" {
// 		errors.Join(err, errors.New("place of birth is required"))
// 	}

// 	if napp.ExpDate.IsZero() {
// 		errors.Join(err, errors.New("place of birth is required"))
// 	}

// 	if len(err) {
// 		return err
// 	}
// 	return nil
// }
