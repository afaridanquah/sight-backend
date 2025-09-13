package verificationbus

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/valueobject"
	"bitbucket.org/msafaridanquah/verifylab-service/business/sdk/yenti"
	"github.com/google/uuid"
)

var (
	ErrPhoneNumberRequired = errors.New("a phone number is required for phone verification")
)

type Verification struct {
	ID               uuid.UUID
	CustomerID       uuid.UUID
	Customer         valueobject.Customer
	BusinessID       uuid.UUID
	Business         valueobject.Business
	VerificationType valueobject.VerificationType
	Summary          valueobject.Summary
	CreatorID        uuid.UUID
	DocumentInsight  string
	CallbackUrl      string
	AmlInsight       AmlInsight
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type NewVerification struct {
	VerificationType valueobject.VerificationType
	CreatorID        uuid.UUID
	BusinessID       uuid.UUID
	Business         valueobject.Business
	CustomerID       uuid.UUID
	Customer         valueobject.Customer
}

type PhoneInsight struct {
	Outcome valueobject.Outcome
	Phone   valueobject.Phone
}

type AmlInsight struct {
	Outcome    valueobject.Outcome
	EntityHits yenti.YentiResponse
}

func (ver *Verification) HasPhoneNumber() bool {
	return (ver.Customer.Phone != valueobject.Phone{})
}

func (ver *Verification) HasIdentifications() bool {
	return len(ver.Customer.Identifications) > 0
}

func (ver *Verification) OpenSanctionMatch(ctx context.Context, yentiClient *yenti.Yenti) error {
	queries := make(map[string]yenti.Query)

	switch {
	case ver.CustomerID != uuid.Nil:
		queries["q1"] = yenti.Query{
			Properties: yenti.Properties{
				Name:         []string{ver.Customer.Person.FullName()},
				Nationality:  []string{ver.Customer.BirthCountry.Alpha2()},
				Jurisdiction: []string{ver.Customer.BirthCountry.Alpha2()},
				BirthDate:    []string{ver.Customer.DateOfBirth.String()},
			},
			Schema: "Person",
		}
	case ver.BusinessID != uuid.Nil:
		queries["q1"] = yenti.Query{
			Properties: yenti.Properties{
				Name:               []string{ver.Business.Name},
				Jurisdiction:       []string{ver.Business.Country.Alpha2()},
				RegistrationNumber: []string{ver.Business.RegistrationNumber},
			},
			Schema: "Company",
		}
	}

	yentiResponse, err := yentiClient.Search(ctx, yenti.NewLookup{
		Weights: yenti.Weights{
			NameLiteralMatch: 0.9,
			NameSoundexMatch: 0.9,
		},
		Queries: queries,
	})
	if err != nil {
		outcome, _ := valueobject.ParseOutcome("ATTENTION_NEEDED")
		ver.AmlInsight = AmlInsight{
			Outcome:    outcome,
			EntityHits: yentiResponse,
		}

		return err
	}

	if len(yentiResponse.Properties["q1"].Results) > 0 {
		outcome, _ := valueobject.ParseOutcome("ATTENTION_NEEDED")
		ver.AmlInsight = AmlInsight{
			Outcome:    outcome,
			EntityHits: yentiResponse,
		}
		return nil
	}

	outcome, _ := valueobject.ParseOutcome("CLEARED")
	ver.AmlInsight = AmlInsight{
		Outcome:    outcome,
		EntityHits: yentiResponse,
	}

	return nil
}
