package customer_test

import (
	"testing"

	"github.com/afaridanquah/verifylab-service/internal/domain/customer"
	cvo "github.com/afaridanquah/verifylab-service/internal/domain/customer/valueobject"
	"github.com/afaridanquah/verifylab-service/internal/valueobject"
)

func TestCustomerNewCustomer(t *testing.T) {
	t.Parallel()
	country, _ := valueobject.NewCountry("US")
	testCases := []struct {
		test        string
		firstName   string
		lastName    string
		middleName  string
		country     valueobject.Country
		expextedErr error
	}{
		{
			test:        "Empty First Name Validation",
			firstName:   "",
			lastName:    "Afari",
			country:     country,
			expextedErr: customer.ErrInvalidPerson,
		},
		{
			test:        "Ok",
			firstName:   "John",
			lastName:    "Doe",
			country:     country,
			expextedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			customerID := cvo.NewID()
			_, err := customer.New(customerID, tc.firstName, tc.lastName, tc.country)
			if err != tc.expextedErr {
				t.Errorf("Expected error %v, got %v", tc.expextedErr, err)
			}
		})
	}

}
