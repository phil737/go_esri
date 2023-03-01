package go_esri

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type machineJSON struct {
	MachineName     string `json:"machineName"`
	MachineAdminUrl string `json:"adminURL"`
}

type machinesJSON struct {
	MachinesList []machineJSON `json:"machines"`
}

// GetMachineNames gets a list of all machines federated in serverName.
//
//	returns machinesJSON struct containing information.
//		getMachinesNames(token, serverName string) (*machinesJSON, error)
func GetMachinesNames(token, serverName string) (*machinesJSON, error) {

	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}
	baseUrl.Path += "/admin/machines/"

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
	var obj machinesJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

type MachineInfoJSON struct {
	MachineNameDomain         string `json:"machineName"`
	AdminURL                  string `json:"adminURL"`
	Platform                  string `json:"platform"`
	ServerStartTime           int64  `json:"ServerStartTime"`
	WebServerMaxHeapSize      int32  `json:"webServerMaxHeapSize"`
	WebServerSSLEnabled       bool   `json:"webServerSSLEnabled"`
	WebServerCertificateAlias string `json:"webServerCertificateAlias"`
	SocMaxHeapSize            int32  `json:"socMaxHeapSize"`
	AppServerMaxHeapSize      int32  `json:"appServerMaxHeapSize"`
	ConfiguredState           string `json:"configuredState"` // <"STARTED"|"STOPPED">
	Synchronize               bool   `json:"synchronize"`
	UnderMaintenance          bool   `json:"underMaintenance"`
}

// GetMachineInformation returns machineInfoJSON struct for given machineName
func GetMachineInformation(token, serverName, machineName string) (*MachineInfoJSON, error) {

	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}
	baseUrl.Path += "/admin/machines/" + machineName

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
	var obj MachineInfoJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// editMachineInformation returns machineInfoJSON struct for given machineName
func EditMachineInformation(token, serverName, machineName string, m MachineInfoJSON) error {
	type statusJSON struct {
		Status string `json:"status"`
	}

	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return err
	}
	baseUrl.Path += "/admin/machines/" + machineName + "/edit"

	// ----------------------------------------- build url encode string to be included in the header body

	v := url.Values{}
	v.Set("token", token)
	v.Add("machineName", machineName) // not documented
	v.Add("adminURL", m.AdminURL)
	v.Add("webServerMaxHeapSize", fmt.Sprintf("%d", m.WebServerMaxHeapSize))
	v.Add("webServerCertificateAlias", m.WebServerCertificateAlias)
	v.Add("socMaxHeapSize", fmt.Sprintf("%d", m.SocMaxHeapSize))
	v.Add("underMaintenance", fmt.Sprintf("%t", m.UnderMaintenance))
	v.Add("f", "json")

	// ----------------------------------------- request
	req := resty.New()

	resp, err := req.SetDebug(true).R().
		SetHeader("Content-type", "application/x-www-form-urlencoded").
		SetBody(string(v.Encode())).
		Post(baseUrl.String())

	if err != nil {
		return err
	}

	// ----------------------------------------- decode json response
	var obj statusJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return err
	}

	if obj.Status != "success" {
		return errors.New("rest/api query error")
	}

	return nil
}
