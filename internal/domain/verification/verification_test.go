package verification_test

import (
	"testing"

	"github.com/afaridanquah/verifylab-service/internal/domain/customer"
	"github.com/afaridanquah/verifylab-service/internal/domain/verification"
	"github.com/afaridanquah/verifylab-service/internal/domain/verification/valueobject"
	vo "github.com/afaridanquverifylab-serviceend/internal/valueobject"
)

func TestAddNewVerification(t *testing.T) {
	t.Parallel()
	country, _ := vo.NewCountry("GH")
	cus, _ := customer.New("John", "Doe", country)

	person, _ := valueobject.NewPerson(cus.ID().String(), cus.FirstName(), cus.MiddleName(), cus.Lastname(), cus.DateOfBirth(), cus.Email(), cus.Country())
	testCases := []struct {
		name             string
		id               valueobject.ID
		person           valueobject.Person
		verificationType valueobject.VerificationType
		expectedErr      error
	}{
		{
			name:             "test passed",
			id:               valueobject.NewID(),
			person:           person,
			verificationType: valueobject.AgeVerification,
			expectedErr:      nil,
		},
	}

	for _, tt := range testCases {
		actual, err := verification.New(tt.id, tt.verificationType, tt.person)
		if err != tt.expectedErr {
			t.Fatalf("expected: %v , got: %v", tt.expectedErr, err)
		}

		t.Logf("actual data: %v", actual)
	}
}
