package postgres

import (
	"encoding/json"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/verificationbus"
	db "bitbucket.org/msafaridanquah/sight-backend/business/sdk/postgres/out"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/vaulti"
	"github.com/jackc/pgx/v5/pgtype"
)

const vaultKey = "pii_key"

type Customer struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	LastName        string           `json:"last_name"`
	OtherNames      string           `json:"other_names"`
	DateOfBirth     string           `json:"date_of_birth"`
	Email           string           `json:"email"`
	PhoneNumber     string           `json:"phone_number"`
	BirthCountry    string           `json:"birth_country"`
	Identifications []Identification `json:"identifications,omitempty"`
	Address         Address          `json:"address"`
}

type AmlInsight struct {
	Outcome    string        `json:"outcome"`
	EntityHits YentiResponse `json:"responses"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Identification struct {
	IdentificationType string `json:"identification_type"`
	Pin                string `json:"pin"`
	PlaceOfBirth       string `json:"place_of_birth"`
	DateOfBirth        string `json:"date_of_birth"`
	Nationality        string `json:"nationality"`
	IssuedDate         string `json:"issued_date"`
	ExpDate            string `json:"exp_date"`
	IssuedCountry      string `json:"country_issued"`
}

func toDBInsertVerification(bus verificationbus.Verification, vaulti *vaulti.Vaulty) (db.InsertVerificationParams, error) {
	customer := Customer{
		ID:           bus.Customer.ID.String(),
		FirstName:    bus.Customer.Person.FirstName,
		LastName:     bus.Customer.Person.LastName,
		MiddleName:   bus.Customer.Person.MiddleName,
		OtherNames:   bus.Customer.Person.OtherNames,
		DateOfBirth:  bus.Customer.DateOfBirth.String(),
		Email:        bus.Customer.Email.String(),
		PhoneNumber:  bus.Customer.Phone.E164Format,
		BirthCountry: bus.Customer.BirthCountry.Alpha2(),
	}

	if len(bus.Customer.Identifications) > 0 {
		identifications := make([]Identification, len(bus.Customer.Identifications))
		for i, v := range bus.Customer.Identifications {
			nationality := v.Nationality.Alpha2()
			issuedCountry := v.CountryIssued.Alpha2()
			expDate := v.ExpDate.Format(time.DateOnly)
			pin, err := vaulti.TransitEncrypt(v.Pin, vaultKey)
			if err != nil {
				return db.InsertVerificationParams{}, err
			}

			identification := Identification{
				IdentificationType: v.IdentificationType.String(),
				Pin:                pin.Ciphertext,
				Nationality:        nationality,
				IssuedCountry:      issuedCountry,
				ExpDate:            expDate,
			}

			identifications[i] = identification
		}

		customer.Identifications = identifications
	}

	customerJson, err := json.Marshal(customer)
	if err != nil {
		return db.InsertVerificationParams{}, err
	}

	yentiJson, err := json.Marshal(bus.AmlInsight.EntityHits)
	if err != nil {
		return db.InsertVerificationParams{}, err
	}

	var yentiResponse YentiResponse
	if err := json.Unmarshal(yentiJson, &yentiResponse); err != nil {
		return db.InsertVerificationParams{}, err

	}

	amlInsight := AmlInsight{
		Outcome:    bus.AmlInsight.Outcome.String(),
		EntityHits: yentiResponse,
	}

	amlInsightJson, err := json.Marshal(amlInsight)

	return db.InsertVerificationParams{
		ID:         bus.ID,
		CustomerID: bus.CustomerID,
		Customer:   customerJson,
		VerificationType: pgtype.Text{
			String: bus.VerificationType.String(),
			Valid:  true,
		},
		AmlInsight: amlInsightJson,
		CreatedAt: pgtype.Timestamp{
			Time:  bus.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  bus.UpdatedAt,
			Valid: true,
		},
	}, nil
}

// ==================================================================
// Aml insight from yenti api
type ResponsePropertyMatcher struct {
	Description string  `json:"description"`
	Coefficient float32 `json:"coefficient"`
	URL         string  `json:"url"`
}

type ResponseProperty struct {
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

type YentiResponse struct {
	Properties map[string]ResponseProperty        `json:"responses"`
	Matcher    map[string]ResponsePropertyMatcher `json:"matcher"`
	Limit      int                                `json:"limit"`
}
