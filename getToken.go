/*
ESRI REST API implementation library.

GetToken function.
*/

package go_esri

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// JSON fields in response from getToken request
type token struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// Queries an ESRI server to obtain an authentication token, returns this token as a string.
func GetToken(username, password, serverName string) (string, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return "", err
	}

	baseUrl.Path += "admin/generateToken" //	portal: sharing/rest/generateToken

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("username", username)
	v.Add("password", password)
	v.Add("client", "requestip")
	v.Add("f", "json")

	// ----------------------------------------- request the token
	req := resty.New()

	// to debug use: req.SetDebug(true).R().
	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())). // convert url encoding to string first
		Post(baseUrl.String())

	if err != nil {
		return "", err
	}

	// ----------------------------------------- decode json response and return token
	var obj token
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return "", err
	}
	// empty token, something went wrong, return body which contains ESRI error message
	if obj.Token == "" {
		return "", errors.New(string(resp.Body()))
	}

	return obj.Token, nil
}
