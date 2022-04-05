/*
ESRI REST API implementation library.

ServiceDetails function.

DetailsJSON and properties structs can be enriched with values found in server doc at
https://developers.arcgis.com/rest/enterprise-administration/server/service.htm, when with Enterprise a more complete description can be found here:
https://developers.arcgis.com/rest/enterprise-administration/enterprise/service.htm

*/

package go_esri

import (
	"encoding/json"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type sProperties struct {
	PortalURL        string `json:"portalURL"`
	VirtualOutputDir string `json:"virtualOutputDir"`
	MaxImageHeight   string `json:"maxImageHeight"`
	MaxRecordCount   string `json:"maxRecordCount"`
	MaxScale         string `json:"maxScale"`
	MinScale         string `json:"minScale"`
	FilePath         string `json:"filePath"`
	CacheOnDemand    string `json:"cacheOnDemand"`
}

type javaHeapSize struct {
	ServiceHeapSize string `json:"javaHeapSize"`
}

// struct returned by services call
type detailsJSON struct {
	ServiceType                  string       `json:"type"`
	ServiceDescription           string       `json:"description"`
	ServiceCapabilities          string       `json:"capabilities"`
	ServiceClusterName           string       `json:"clusterName"`
	ServiceMinInstPerNode        int32        `json:"minInstancesPerNode"`
	ServiceMaxInstPerNode        int32        `json:"maxInstancesPerNode"`
	ServiceMaxWaitTime           int32        `json:"maxWaitTime"`
	ServiceMaxIdelTime           int32        `json:"maxIdleTime"`
	ServiceMaxUsageTime          int32        `json:"maxUsageTime"`
	ServiceRecycleInterval       int32        `json:"recycleInterval"`
	ServiceProvider              string       `json:"provider"`
	ServiceLoadBalancing         string       `json:"loadBalancing"`
	ServiceKeepAliveInterval     int32        `json:"keepAliveInterval"`
	ServiceIsolationLevel        string       `json:"isolationLevel"`
	ServicehasVersionedData      bool         `json:"hasVersionedData"`
	ServiceInstancesPerContainer int32        `json:"instancesPerContainer"`
	ServiceMaxUploadFileSize     int32        `json:"maxUploadFileSize"`
	ServiceRecycleStartTime      string       `json:"recycleStartTime"`
	ServiceProperties            sProperties  `json:"properties"`
	ServiceFramework             javaHeapSize `json:"frameworkProperties"`
}

// Returns DetailsJSON struct with service info
func ServiceDetails(token, serverName, folder, serviceFullName string) (*detailsJSON, error) {

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
	var obj detailsJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
