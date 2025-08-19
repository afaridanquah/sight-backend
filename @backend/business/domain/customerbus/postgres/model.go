package postgres

import (
	"encoding/json"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/customerbus/valueobject"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/vaulti"
)

const vaultKey = "pii_key"

type Identification struct {
	IdentificationType string  `json:"identification_type"`
	Pin                string  `json:"pin"`
	Nationality        *string `json:"nationality,omitempty"`
	IssuedDate         *string `json:"issued_date,omitempty"`
	IssuedAt           *string `json:"issued_at,omitempty"`
	ExpDate            *string `json:"exp_date,omitempty"`
	IssuedCountry      string  `json:"issued_country,omitempty"`
}

func toDBIdentifications(ids []valueobject.Identification, vaulti *vaulti.Vaulty) ([]byte, error) {
	if len(ids) > 0 {
		identifications := make([]Identification, len(ids))
		for i, v := range ids {
			nationality := v.Nationality.Alpha2()
			issuedCountry := v.IssedCountry.Alpha2()
			expDate := v.ExpDate.Format(time.DateOnly)
			pin, err := vaulti.TransitEncrypt(v.Pin, vaultKey)
			if err != nil {
				return nil, err
			}

			identification := Identification{
				IdentificationType: v.IdentificationType.String(),
				Pin:                pin.Ciphertext,
				Nationality:        &nationality,
				IssuedCountry:      issuedCountry,
				ExpDate:            &expDate,
			}

			identifications[i] = identification
		}
		identificationsJson, err := json.Marshal(identifications)
		if err != nil {
			return nil, err
		}
		return identificationsJson, nil
	}

	return []byte{}, nil
}

func toBusCustomer(resp db.Customers, vaulti *vaulti.Vaulty) (customerbus.Customer, error) {
	person, err := valueobject.NewPerson(resp.FirstName, resp.LastName, &resp.MiddleName.String, nil)
	if err != nil {
		return customerbus.Customer{}, err
	}

	email, err := valueobject.NewEmail(resp.Email.String)
	if err != nil {
		return customerbus.Customer{}, err
	}

	var identifications []Identification
	if err := json.Unmarshal(resp.Identifications, &identifications); err != nil {
		return customerbus.Customer{}, err
	}

	busidentifications := make([]valueobject.Identification, len(identifications))

	if len(resp.Identifications) > 0 {

		for i, v := range identifications {
			decrypted, err := vaulti.TransitDecrypt(v.Pin, vaultKey)
			if err != nil {
				return customerbus.Customer{}, err
			}

			identity, err := valueobject.NewIdentification(v.IdentificationType, decrypted.Plaintext, v.IssuedCountry, v.IssuedDate, v.Nationality, v.ExpDate)
			if err != nil {
				return customerbus.Customer{}, err
			}
			busidentifications[i] = identity
		}
	}

	birthCountry, err := valueobject.NewCountry(resp.BirthCountry.String)
	if err != nil {
		return customerbus.Customer{}, err
	}

	dob, err := valueobject.ParseDateOfBirth(resp.DateOfBirth.Time)
	if err != nil {
		return customerbus.Customer{}, err
	}

	customer := customerbus.Customer{
		ID:              resp.ID,
		Person:          person,
		UserID:          resp.CreatorID.UUID,
		BusinessID:      resp.BusinessID.UUID,
		Email:           email,
		BirthCountry:    birthCountry,
		DateOfBirth:     dob,
		Identifications: busidentifications,
		CreatedAt:       resp.CreatedAt.Time,
		UpdatedAt:       resp.UpdatedAt.Time,
	}

	if resp.PhoneNumber.String != "" {
		phone, err := valueobject.ParseIntlPhone(resp.PhoneNumber.String)
		if err != nil {
			return customerbus.Customer{}, err
		}
		customer.PhoneNumber = phone
	}

	return customer, nil
}
