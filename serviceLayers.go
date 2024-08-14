package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// struct returned by services call
type layerJSON struct {
	LayerId          int32  `json:"id"`
	LayerName        string `json:"name"`
	LayerType        string `json:"type"`
	LayerDescription string `json:"description"`
}

type layersJSON struct {
	Layers []layerJSON `json:"layers"`
}

// Returns layersJSON struct with service layers info
func ServiceLayers(token, serverName, folder, serviceFullName string) (*layersJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/rest/services/"
	baseUrl.Path += folder + "/"
	baseUrl.Path += serviceFullName
	baseUrl.Path += "/layers"

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
	var obj layersJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
