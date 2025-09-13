package verificationapp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

type Verification struct {
	ID               string       `json:"id"`
	Customer         Customer     `json:"customer"`
	Outcome          string       `json:"outcome"`
	Business         *Business    `json:"business,omitempty"`
	VerificationType string       `json:"verification_type"`
	AmlInsight       AmlInsight   `json:"aml_insight"`
	PhoneInsight     PhoneInsight `json:"phone_insight"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
}

type AmlInsight struct {
	Outcome    string        `json:"outcome"`
	EntityHits YentiResponse `json:"entity_hits"`
}

type PhoneInsight struct {
	Status              string `json:"status"`
	InternationalFormat string `json:"international_format"`
	PhoneType           string `json:"phone_type"`
	Country             string `json:"country"`
	NationalFormat      string `json:"national_format"`
}

type YentiResponse struct {
	Properties map[string]EntityProperty `json:"responses"`
}

type EntityProperty struct {
	Status  int `json:"status"`
	Results []struct {
		ID         string `json:"id"`
		Caption    string `json:"caption"`
		Schema     string `json:"schema"`
		Properties struct {
			Name []string `json:"name"`
		} `json:"properties"`
		Datasets   []string       `json:"datasets"`
		Referents  []string       `json:"referents"`
		Target     bool           `json:"target"`
		FirstSeen  string         `json:"first_seen"`
		LastSeen   string         `json:"last_seen"`
		LastChange string         `json:"last_change"`
		Score      float64        `json:"score"`
		Features   map[string]any `json:"features"`
		Match      bool           `json:"match"`
		Token      string         `json:"token"`
	} `json:"results"`
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Query struct {
		ID         string `json:"id"`
		Schema     string `json:"schema"`
		Properties struct {
			Name []string `json:"name"`
		} `json:"properties"`
	} `json:"query"`
}

type Customer struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	MiddleName      string           `json:"middle_name"`
	OtherNames      string           `json:"other_names"`
	DateOfBirth     string           `json:"date_of_birth"`
	Identifications []Identification `json:"identifications"`
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
}

type Identification struct {
	IdentificationType string  `json:"identification_type"`
	Pin                string  `json:"pin"`
	Nationality        *string `json:"nationality,omitempty"`
	IssuedDate         *string `json:"issued_date,omitempty"`
	IssuedAt           *string `json:"issued_at,omitempty"`
	ExpDate            *string `json:"exp_date,omitempty"`
	IssuedCountry      string  `json:"issued_country,omitempty"`
}

type Business struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewVerification struct {
	BusinessID       string `json:"business_id" validate:"required_if=CustomerID ''"`
	CustomerID       string `json:"customer_id" validate:"required_if=BusinessID ''"`
	VerificationType string `json:"verification_type" validate:"required"`
}

func (o NewVerification) Validate() error {
	if err := ierr.Check(o); err != nil {
		return fmt.Errorf("validate new verification failed: %w", err)
	}

	return nil
}

func toBusNewCustomerVerification(napp NewVerification, customer valueobject.Customer) (verificationbus.NewVerification, error) {
	verificationType, err := valueobject.ParseVerificationType(napp.VerificationType)
	if err != nil {
		return verificationbus.NewVerification{}, err
	}

	return verificationbus.NewVerification{
		VerificationType: verificationType,
		CustomerID:       customer.ID,
		Customer:         customer,
	}, nil
}

func toBusNewBusinessVerification(napp NewVerification, business valueobject.Business) (verificationbus.NewVerification, error) {
	verificationType, err := valueobject.ParseVerificationType(napp.VerificationType)
	if err != nil {
		return verificationbus.NewVerification{}, err
	}

	return verificationbus.NewVerification{
		VerificationType: verificationType,
		BusinessID:       business.ID,
		Business:         business,
	}, nil
}

func toAppVerification(vbus verificationbus.Verification) Verification {
	identifications := make([]Identification, len(vbus.Customer.Identifications))

	if len(vbus.Customer.Identifications) > 0 {
		for i, v := range vbus.Customer.Identifications {
			identifications[i] = Identification{
				IdentificationType: v.IdentificationType.String(),
				Pin:                v.Pin,
			}
		}
	}

	var amlInsight AmlInsight

	if vbus.VerificationType.String() == "AML_SCREENING" {
		marshHits, err := json.Marshal(vbus.AmlInsight.EntityHits)
		if err != nil {
			fmt.Printf("aml entityies %s", marshHits)
		}

		var yentiResponse YentiResponse

		if err := json.Unmarshal(marshHits, &yentiResponse); err != nil {
			fmt.Printf("aml entityies %v", err)
		}

		amlInsight.EntityHits = yentiResponse
		amlInsight.Outcome = "CLEARED"
	}

	return Verification{
		ID:               vbus.ID.String(),
		VerificationType: vbus.VerificationType.String(),
		Customer: Customer{
			ID:              vbus.Customer.ID.String(),
			FirstName:       vbus.Customer.Person.FirstName,
			LastName:        vbus.Customer.Person.LastName,
			MiddleName:      vbus.Customer.Person.MiddleName,
			OtherNames:      vbus.Customer.Person.OtherNames,
			DateOfBirth:     vbus.Customer.DateOfBirth.String(),
			Identifications: identifications,
			Email:           vbus.Customer.Email.String(),
			Phone:           vbus.Customer.Phone.E164Format,
		},
		PhoneInsight: PhoneInsight{
			NationalFormat:      vbus.Customer.Phone.NationalFormat,
			InternationalFormat: vbus.Customer.Phone.E164Format,
			PhoneType:           vbus.Customer.Phone.PhoneType,
			Country:             vbus.Customer.Phone.Country.Alpha2(),
		},
		AmlInsight: amlInsight,
		CreatedAt:  vbus.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  vbus.UpdatedAt.Format(time.RFC3339),
	}
}

func toBusVoCustomer(customer customerbus.Customer) (valueobject.Customer, error) {
	person, err := valueobject.NewPerson(customer.Person.FirstName, customer.Person.LastName, &customer.Person.MiddleName, &customer.Person.OtherNames)
	if err != nil {
		return valueobject.Customer{}, err
	}

	countryOfBirth, err := valueobject.NewCountry(customer.BirthCountry.Alpha2())
	if err != nil {
		return valueobject.Customer{}, err
	}

	customerIdentifications := make([]valueobject.Identification, len(customer.Identifications))
	if len(customer.Identifications) > 0 {
		for i, v := range customer.Identifications {
			nationality := v.Nationality.String()
			expDate := v.ExpDate.Format(time.DateOnly)
			identification, _ := valueobject.NewIdentification(v.IdentificationType.String(), v.Pin, v.IssedCountry.String(), v.IssuedDate, &nationality, &expDate)
			customerIdentifications[i] = identification
		}
	}

	var dateOfBirth valueobject.DateOfBirth
	var emailAddress valueobject.Email
	var phoneNumber valueobject.Phone

	if customer.DateOfBirth.String() != "" {
		dob, err := valueobject.NewDateOfBirth(customer.DateOfBirth.String())
		if err != nil {
			return valueobject.Customer{}, err
		}
		dateOfBirth = dob
	}

	if customer.Email.String() != "" {
		email, err := valueobject.NewEmail(customer.Email.String())
		if err != nil {
			return valueobject.Customer{}, err
		}
		emailAddress = email
	}

	if !customer.PhoneNumber.IsZero() {
		phone, err := valueobject.ParsePhone(customer.PhoneNumber.Country.Alpha2(), customer.PhoneNumber.E164Format)

		if err != nil {
			fmt.Printf("parse phone: %v", err)
			return valueobject.Customer{}, err
		}
		phoneNumber = phone
	}

	busCustomer, err := valueobject.ParseCustomer(customer.ID, person, &dateOfBirth, &countryOfBirth, &customerIdentifications, &emailAddress, &phoneNumber)
	if err != nil {
		return valueobject.Customer{}, err
	}

	return busCustomer, nil
}

// ===============================================================

type PhoneStatus struct {
	a string
}

var (
	VALID   = PhoneStatus{"VALID"}
	INVALID = PhoneStatus{"INVALID"}
)

func ParsePhoneStatus(s string) PhoneStatus {
	upper := strings.ToUpper(s)

	switch upper {
	case "VALID":
		return VALID
	case "INVALID":
		return INVALID
	default:
		return INVALID
	}
}
