package twilio_test

import (
	"testing"

	"bitbucket.org/msafaridanquah/sight-backend/business/sdk/twilio"
)

func TestSend(t *testing.T) {
	tw, err := twilio.New()

	if err != nil {
		t.Fatalf("failed %v", err.Error())
	}

	if err := tw.SendSMS("+17017306525", "Test"); err != nil {
		t.Fatalf("sendsms")
	}

}
