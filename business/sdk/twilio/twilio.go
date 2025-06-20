package twilio

import (
	"encoding/json"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio struct {
	params *twilioApi.CreateMessageParams
	client *twilio.RestClient
}

func New() (*Twilio, error) {

	accountSid := "AC669496069f297202e6526c6428896165"
	authToken := "34132a0483074e37520f3353555a7895"
	fromNumber := "+13864338434"

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
