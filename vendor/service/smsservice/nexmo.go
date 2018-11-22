package smsservice

import (
	"fmt"
	"errors"
	"github.com/ammario/nexmo"
)

var nexmoUrl string = "https://rest.nexmo.com/sms/json"
var apiKey string = "b0ab5e93"
var apiSecret string = "edc22192add8af46"

type NexmoVerifyCode struct {

}

func (self *NexmoVerifyCode) SendVerifyCode(zone, mobile, code string) error  {
	msg := fmt.Sprintf("Your verify code is [%s]", code)
	return self.Send(combinePhoneFormat(zone, mobile), msg)
}
func (self *NexmoVerifyCode) Send(phone, content string) error {

	client := nexmo.New(apiKey, apiSecret)
	if _, err := client.SMS("VVChat", phone, content); err != nil {
		return errors.New(ERROR_SEND_STRING)
	}
	return nil

}