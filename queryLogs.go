/*
query ArcGIS server logs.

10.12.2021 initial release
01.09.2023 added process, user, elapsed, thread, methodname and requestid (ProcID) in logMessageJSON struct
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
	Type       string `json:"type"`
	Message    string `json:"message"`
	Time       int64  `json:"time"`
	Source     string `json:"source"`
	Code       int32  `json:"code"`
	ProcessID  string `json:"process"`
	UserName   string `json:"user"`
	Elapsed    string `json:"elapsed"`
	Thread     string `json:"thread"`
	MethodName string `json:"methodName"`
	RequestID  string `json:"requestID"`
}

type logResponseJSON struct {
	HasMore     bool             `json:"hasMore"`
	StartTime   int64            `json:"startTime"`
	EndTime     int64            `json:"endTime"`
	LogMessages []logMessageJSON `json:"logMessages"`
}

// query ArcGIS server logs, returns records with a log levelType at or more severe than given (SEVERE, WARNING, INFO, FINE, VERBOSE, DEBUG).
// returns logs between startTime and endTime (values in milliseconds)
func QueryLogs(token, serverName, levelType string, startTime, endTime int64) (*logResponseJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/logs/query"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("startTime", fmt.Sprint(startTime)) // 2021-12-07T00:55:23 or milliseconds
	v.Add("endTime", fmt.Sprint(endTime))
	v.Add("level", levelType)
	v.Add("filterType", "json")
	v.Add("pageSize", "1000")
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
