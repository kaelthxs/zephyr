package sms

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type SmsClient struct {
	ApiID string
}

func NewSmsClient(apiID string) *SmsClient {
	return &SmsClient{ApiID: apiID}
}

func (s *SmsClient) SendSms(to, message string) error {
	baseURL := "https://sms.ru/sms/send"

	params := url.Values{}
	params.Set("api_id", s.ApiID)
	params.Set("to", to)
	params.Set("msg", message)
	params.Set("json", "1")

	url := baseURL + "?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	return err
}
