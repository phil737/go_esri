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
type FoldersJSON struct {
	FolderName   string `json:"folderName"`
	Description  string `json:"description"`
	WebEncrypted bool   `json:"webEncrypted"`
	IsDefault    bool   `json:"isDefault"`
}

type rootItemsJSON struct {
	FolderName    string        `json:"folderName"`
	Description   string        `json:"description"`
	FoldersDetail []FoldersJSON `json:"foldersDetail"`
}

// Returns struct list of folders by listing items at the root folder level.
func RootFolders(token, serverName string) ([]FoldersJSON, error) {

	// ----------------------------------------- build and validate url
	baseUrl, err := url.Parse(serverName)
	if err != nil {
		return nil, err
	}

	baseUrl.Path += "/admin/services/"

	// ----------------------------------------- build url encode string to be included in the header body
	v := url.Values{}
	v.Set("token", token)
	v.Add("f", "json")
	v.Add("detail", "true")

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
	var obj rootItemsJSON
	err = json.Unmarshal(resp.Body(), &obj)
	if err != nil {
		return nil, err
	}

	return obj.FoldersDetail, nil
}
