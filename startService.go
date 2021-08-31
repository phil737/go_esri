/*
ESRI REST API implementation library.

Start a service.
*/

package go_esri

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type startJSON struct {
	Status string `json:"status"`
}

// Starts service, returns error or nil if success.
func StartService(token, serverName, folder, serviceFullName string) error {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return err
	}

	baseUrl.Path += "admin/services/"
	baseUrl.Path += folder
	baseUrl.Path += serviceFullName
	baseUrl.Path += "/start"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("f", "json")

	// ----------------------------------------- request the token
	req := resty.New()

	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())). // convert url encoding to string first
		Post(baseUrl.String())

	if err != nil {
		return err
	}

	// ----------------------------------------- decode json response and return token
	var obj startJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return err
	}

	if obj.Status != "success" {
		return errors.New(string(resp.Body()))
	}

	return nil
}
