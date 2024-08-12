package params_test

import (
	"errors"
	"testing"

	"github.com/afaridanquah/verifylab-backend/internal/params"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func TestCreateVerficationValidate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		input   params.CreateVerificationRequest
		withErr bool
	}{
		{
			name: "OK",
			input: params.CreateVerificationRequest{
				Customer: params.CreateVerificationRequestCustomer{
					FirstName: "Jane",
					LastName:  "Doe",
				},
				VerificationType: "Passport",
			},
			withErr: false,
		},
		{
			name: "Err",
			input: params.CreateVerificationRequest{
				Customer: params.CreateVerificationRequestCustomer{
					LastName: "Jane",
				},
				VerificationType: "Passport",
			},
			withErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualErr := tt.input.Validate()
			if (actualErr != nil) != tt.withErr {
				t.Fatalf("expected error %t, got %s", tt.withErr, actualErr)
			}

			var ierr validation.Errors
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}

			t.Logf("err: %v", actualErr)
		})
	}

}
