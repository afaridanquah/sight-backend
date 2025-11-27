package twilio

import (
	"context"
	"encoding/json"
	"fmt"

	"bitbucket.org/msafaridanquah/sight-backend/foundation/envvar"
	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Config struct {
	Env *envvar.Configuration
	Log *logger.Logger
}

type Twilio struct {
	params *twilioApi.CreateMessageParams
	client *twilio.RestClient
}

func New(conf Config) (*Twilio, error) {

	get := func(v string) string {
		res, err := conf.Env.Get(v)
		if err != nil {
			conf.Log.Error(context.Background(), "env failed")
		}

		return res
	}

	accountSid := get("TWILIO_ACCOUNT_ID")
	authToken := get("TWILIO_AUTH_TOKEN")
	fromNumber := get("TWILIO_FROM_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(fromNumber)

	return &Twilio{
		params: params,
		client: client,
	}, nil
}

func (tw *Twilio) SendSMS(toNumber, msg string) error {
	tw.params.SetBody(msg)
	resp, err := tw.client.Api.CreateMessage(tw.params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
		return err
	}

	response, _ := json.Marshal(*resp)
	fmt.Println("Response: " + string(response))

	return nil
}
