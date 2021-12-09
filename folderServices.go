/*
ESRI REST API implementation library.

FolderServices function.
*/

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// struct containing service information
type ServicesJSON struct {
	FolderName  string `json:"folderName"`
	ServiceName string `json:"serviceName"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type folderJSON struct {
	FolderName  string         `json:"folderName"`
	Description string         `json:"description"`
	Services    []ServicesJSON `json:"services"`
}

// Returns struct list of services in given folder. For root folder string should be empty.
func FolderServices(token, serverName, folder string) ([]ServicesJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/services/"
	baseUrl.Path += folder

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("f", "json")

	// ----------------------------------------- request
	req := resty.New()

	// to debug use: req.SetDebug(true).R().
	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())). // convert url encoding to string first
		Post(baseUrl.String())

	if err != nil {
		return nil, err
	}

	// ----------------------------------------- decode json response
	var obj folderJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return obj.Services, nil
}
