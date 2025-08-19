package valueobject_test

import (
	"testing"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/verificationbus/valueobject"
)

func TestParsePhone(t *testing.T) {
	testCases := []struct {
		name        string
		country     string
		phoneNumber string
	}{
		{
			name:        "ok",
			country:     "US",
			phoneNumber: "7017306323",
		},
	}

	for _, tt := range testCases {
		phone, err := valueobject.ParsePhone(tt.country, tt.phoneNumber)
		if err != nil {
			t.Errorf("parse phone: %v ", err)
		}

		t.Logf("parsed phone %v", phone)
	}
}
