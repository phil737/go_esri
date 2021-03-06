/*
ESRI REST API implementation library.

https://developers.arcgis.com/rest/enterprise-administration/enterprise/service-status.htm

ServiceStatus function.
*/

package go_esri

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type responseJSON struct {
	RealTime  string `json:"realTimeState"`
	ConfState string `json:"configuredState"`
}

// Returns string with service real time state. Configured state can be queried via serviceDetails
// serviceFullName is the service name followed by its type.
// example: "SampleWorldCities.MapServer", serverName can be in the form: https://www.myserver.com/server/ or https://ags.myserver.com:6443/arcgis/
func ServiceState(token, serverName, folder, serviceFullName string) (string, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return "", err
	}

	baseUrl.Path += "/admin/services/" // admin is part of root url
	baseUrl.Path += folder + "/"
	baseUrl.Path += serviceFullName
	baseUrl.Path += "/status"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("f", "json")

	// ----------------------------------------- request
	req := resty.New()

	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())).
		Post(baseUrl.String())

	if err != nil {
		return "", err
	}

	// ----------------------------------------- decode json response
	var obj responseJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return "", err
	}

	if obj.RealTime == "" {
		return "", errors.New(string(resp.Body()))
	}

	return obj.RealTime, nil
}
