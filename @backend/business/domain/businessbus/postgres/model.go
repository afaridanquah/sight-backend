package postgres

import (
	"encoding/json"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus"
	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/businessbus/valueobject"
	db "bitbucket.org/msafaridanquah/verifylab-service/business/sdk/postgres/out"
)

type Owner struct {
	Person struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		MiddleName string `json:"middle_name"`
		OtherNames string `json:"other_names"`
	} `json:"person"`
	Percentage        float32 `json:"percentage"`
	Address           Address `json:"address"`
	CountryOfResident string  `json:"country_of_resident"`
}

type Address struct {
	Line1         string `json:"line_1"`
	Line2         string `json:"line_2"`
	City          string `json:"city"`
	StateOrRegion string `json:"state_or_region"`
	Country       string `json:"country"`
}

func toDBAddress(a valueobject.Address) ([]byte, error) {
	var addr Address
	if !a.IsEmpty() {
		addr.Line1 = a.Line1
		addr.Line2 = a.Line2
		addr.City = a.City
		addr.StateOrRegion = a.StateOrRegion
		addr.Country = a.Country.Alpha2()
	}
	address, err := json.Marshal(addr)
	if err != nil {
		return []byte{}, err
	}
	return address, nil
}

func toBusBusiness(res db.Businesses) (businessbus.Business, error) {
	entity, err := valueobject.ParseBusinessEntity(string(res.Entity))
	if err != nil {
		return businessbus.Business{}, err
	}

	countryOfCorp, err := valueobject.NewCountry(res.Jurisdiction)
	if err != nil {
		return businessbus.Business{}, err
	}

	var address Address
	if err := json.Unmarshal(res.Address, &address); err != nil {
		return businessbus.Business{}, err
	}

	return businessbus.Business{
		ID:                     res.ID,
		LegalName:              res.LegalName,
		DoingBusinessAs:        res.Dba,
		TaxID:                  res.TaxID.String,
		Entity:                 entity,
		CountryOfIncorporation: countryOfCorp,
		CreatedAt:              res.CreatedAt.Time,
		UpdatedAt:              res.UpdatedAt.Time,
	}, nil

}
