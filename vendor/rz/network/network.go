package network

import (
	"github.com/sendgrid/rest"
)

func Get(url string, queryParams map[string]string, headers map[string]string) (resp *rest.Response, err error) {

	request := rest.Request{
		Method:      rest.Get,
		BaseURL:     url,
		Headers:     headers,
		QueryParams: queryParams,
	}
	response, err := rest.API(request)
	if err != nil {

		return nil, err
	}

	return response, nil
}