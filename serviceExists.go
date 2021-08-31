/*
ESRI REST API implementation library.

ServiceExists function.
*/

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// JSON fields in response from getToken request
type existsJSON struct {
	Exists bool `json:"exists"`
}

// Returns true if service exists, false otherwise, for root services folder should be empty.
// serverName can be in the form: https://www.myserver.com/server/ or https://ags.myserver.com:6443/arcgis/
func ServiceExists(token, serverName, folder, serviceName, serviceType string) (bool, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return false, err
	}
	baseUrl.Path += "admin/services/exists"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("f", "json")
	v.Add("folderName", folder)
	v.Add("serviceName", serviceName)
	v.Add("type", serviceType)

	// ----------------------------------------- request the token
	req := resty.New()

	// to debug use: req.SetDebug(true).R().
	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())). // convert url encoding to string first
		Post(baseUrl.String())
	if err != nil {
		return false, err
	}

	// ----------------------------------------- decode json response and return token
	var obj existsJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return false, err
	}

	return obj.Exists, nil
}
