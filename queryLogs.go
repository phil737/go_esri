/*


 */

package go_esri

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// struct containing service information
type logMessageJSON struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
	Source  string `json:"source"`
	Code    int32  `json:"code"`
}

type logResponseJSON struct {
	HasMore     bool             `json:"hasMore"`
	StartTime   int64            `json:"startTime"`
	EndTime     int64            `json:"endTime"`
	LogMessages []logMessageJSON `json:"logMessages"`
}

// query ArcGIS server logs, only records with a log level at or more severe than this (SEVERE, WARNING, INFO, FINE, VERBOSE, DEBUG).
func QueryLogs(token, serverName, levelType string, endT int64) (*logResponseJSON, error) {

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
	v.Add("endTime", fmt.Sprint(endT))
	v.Add("level", levelType)
	v.Add("filterType", "json")
	v.Add("pageSize", "100")
	v.Add("filter", "{\"server\": \"*\", \"services\": \"*\", \"machines\":\"*\" }")
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
	var obj logResponseJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
