/*


 */

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// struct containing service information
type logMessageJSON struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Source  string `json:"source"`
	Code    string `json:"code"`
}

type logResponseJSON struct {
	HasMore   string           `json:"hasMore"`
	StartTime string           `json:"startTime"`
	EndTime   string           `json:"endTime"`
	Services  []logMessageJSON `json:"services"`
}

// query ArcGIS server logs
func QueryLogs(token, serverName string) (*logResponseJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/logs/query"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("startTime", "") // 2021-12-07T00:55:23
	v.Add("endTime", "")
	v.Add("level", "SEVERE")
	v.Add("filterType", "json")
	v.Add("pageSize", "100")
	v.Add("filter", "{\"server\": \"*\", \"services\": \"*\", \"machines\":\"*\" }")
	v.Add("f", "json")

	// ----------------------------------------- request the token
	req := resty.New()

	// to debug use: req.SetDebug(true).R().
	resp, err := req.R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())). // convert url encoding to string first
		Post(baseUrl.String())

	if err != nil {
		return nil, err
	}

	// ----------------------------------------- decode json response and return token
	var obj logResponseJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
