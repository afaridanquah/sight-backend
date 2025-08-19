package businessbus

import (
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	"github.com/google/uuid"
)

type Business struct {
	ID                     uuid.UUID
	LegalName              string
	DoingBusinessAs        string
	TaxID                  string
	Entity                 valueobject.BusinessEntity
	CountryOfIncorporation valueobject.Country
	Address                valueobject.Address
	EmailAddresses         []valueobject.Email
	PhoneNumbers           []valueobject.Phone
	Website                string
	Owners                 []valueobject.Owner
	Documents              []valueobject.Document
	AdminID                string
	Status                 valueobject.Status
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type NewBusiness struct {
	LegalName       string
	DoingBusinessAs string
	AdminID         string
	Entity          valueobject.BusinessEntity
	EmailAddresses  []valueobject.Email
	PhoneNumbers    []valueobject.Phone
}

type UpdateBusiness struct {
	LegalName       *string
	DoingBusinessAs *string
	EmailAddresses  *[]valueobject.Email
	PhoneNumbers    *[]valueobject.Phone
}
