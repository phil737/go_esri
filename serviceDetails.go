/*
ESRI REST API implementation library.

ServiceDetails function.

*/

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type DetailsJSON struct {
	ServiceType           string `json:"type"`
	ServiceDescription    string `json:"description"`
	ServiceCapabilities   string `json:"capabilities"`
	ServiceClusterName    string `json:"clusterName"`
	ServiceMinInstPerNode int32  `json:"minInstancesPerNode"`
	ServiceMaxInstPerNode int32  `json:"maxInstancesPerNode"`
	ServiceMaxWaitTime    int32  `json:"maxWaitTime"`
	ServiceMaxIdelTime    int32  `json:"maxIdleTime"`
	ServiceMaxUsageTime   int32  `json:"maxUsageTime"`
}

// Returns DetailsJSON struct with service info
func ServiceDetails(token, serverName, folder, serviceFullName string) (*DetailsJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/services/"
	baseUrl.Path += folder + "/"
	baseUrl.Path += serviceFullName

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
		return nil, err
	}

	// ----------------------------------------- decode json response
	var obj DetailsJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}