/*
ESRI REST API implementation library.

ServiceManifest function.

iteminfo and manifest give more info on a service

call example:
https://geo.geomsf.org/server/admin/services/wrl/wrl_ref_ppl.MapServer/iteminfo?f=pjson
https://geo.geomsf.org/server/admin/services/wrl/wrl_ref_ppl.MapServer/iteminfo/manifest/manifest.json?f=pjson
*/

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type sdataBases struct {
	ByReference                    bool   `json:"byReference"`
	OnServerWorkspaceFactoryProgID string `json:"onServerWorkspaceFactoryProgID"`
	OnServerConnectionString       string `json:"onServerConnectionString"`
	OnPremiseConnectionString      string `json:"onPremiseConnectionString"`
	OnServerName                   string `json:"onServerName"`
	OnPremisePath                  string `json:"onPremisePath"`
}

// struct returned by services call
type manifestJSON struct {
	DataBases []sdataBases `json:"databases"`
}

// Returns DetailsJSON struct with service info
func ServiceManifest(token, serverName, folder, serviceFullName string) (*manifestJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/services/"
	baseUrl.Path += folder + "/"
	baseUrl.Path += serviceFullName + "/"
	baseUrl.Path += "iteminfo/manifest/manifest.json"

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
	var obj manifestJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
