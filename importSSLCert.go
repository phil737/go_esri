/*
https://developers.arcgis.com/rest/enterprise-administration/server/importexistingservercertificate.htm
*/
package go_esri

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/go-resty/resty/v2"
)

//	   ImportSSLCert imports pfx certificate to ArcGIS server.
//
//		   machineName is the ArcGIS server machine name where the certificate will be uploaded.
//			certFileName must be a valid PathName to the certificate file with its password in certPassword.
//			alias is the unique name used by the server to identify the certifcate.
func ImportSSLCert(token, serverName, machineName, alias, certPassword, certFileName string) error {

	type statusJSON struct {
		Status string `json:"status"`
	}

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return err
	}

	baseUrl.Path += "/admin/machines/" + machineName + "/sslcertificates/importExistingServerCertificate"

	// ----------------------------------------- request
	req := resty.New()

	resp, err := req.R().
		SetFile("certFile", certFileName).
		ForceContentType("application/octet-stream").
		SetFormData(map[string]string{
			"alias":        alias,
			"certPassword": certPassword,
			"f":            "json",
			"token":        token,
		}).
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
