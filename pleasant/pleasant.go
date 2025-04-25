package pleasant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

const (
	PathRootFolder   = "/api/v5/rest/folders/root"
	PathEntry        = "/api/v5/rest/entries"
	PathFolders      = "/api/v5/rest/folders"
	PathAccessLevels = "/api/v5/rest/accesslevels"
	PathSearch       = "/api/v5/rest/search"
	PathServerInfo   = "/api/v5/rest/GetServerInfo"
	PathPwStr        = "/api/v5/rest/passwordstrength"
)

type BearerToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetBearerToken(baseUrl, username, password string) (*BearerToken, error) {
	path := "/OAuth2/Token"

	data := url.Values{}
	data.Add("grant_type", "password")
	data.Add("username", username)
	data.Add("password", password)

	res, err := postRequestForm(baseUrl, path, data)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bearerToken := &BearerToken{}

	err = decodeJsonBody(res.Body, bearerToken)
	if err != nil {
		return nil, err
	}

	return bearerToken, nil
}

func GetJsonBody(baseUrl, path, bearerToken string) (string, error) {
	res, err := getRequest(baseUrl, path, bearerToken)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func PostJsonString(baseUrl, path, jsonString, bearerToken string) (string, error) {
	res, err := postRequestJsonString(baseUrl, path, jsonString, bearerToken)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func PatchJsonString(baseUrl, path, jsonString, bearerToken string) (string, error) {
	res, err := patchRequestJsonString(baseUrl, path, jsonString, bearerToken)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func DeleteJsonString(baseUrl, path, jsonString, bearerToken string) (string, error) {
	res, err := deleteRequestJsonString(baseUrl, path, jsonString, bearerToken)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func PostSearch(baseUrl, query, bearerToken string) (string, error) {
	queryJson, err := json.Marshal(map[string]string{"Search": query})
	if err != nil {
		return "", err
	}

	res, err := postRequestJsonString(baseUrl, PathSearch, string(queryJson), bearerToken)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetIdByResourcePath(baseUrl, resourcePath, resourceType, bearerToken string) (string, error) {
	if resourceType != "entry" && resourceType != "folder" {
		return "", ErrInvalidResourceType
	}

	splitPath := strings.Split(resourcePath, "/")

	if splitPath[0] != "Root" {
		return "", ErrPathStartIncorrect
	}

	resourceName := splitPath[len(splitPath)-1]

	result, err := PostSearch(baseUrl, resourceName, bearerToken)
	if err != nil {
		return "", err
	}

	j, err := unmarshalSearchResponse(result)
	if err != nil {
		return "", err
	}

	var count int
	var id string

	if resourceType == "entry" {
		for _, c := range j.Credentials {
			if c.Path+resourceName == resourcePath && c.Name == resourceName {
				count++
				id = c.Id
			}
		}
	} else if resourceType == "folder" {
		for _, c := range j.Groups {
			if strings.TrimSuffix(c.FullPath, "/") == resourcePath && c.Name == resourceName {
				count++
				id = c.Id
			}
		}
	}

	if count > 1 {
		return "", ErrAmbiguousResult
	} else if count == 0 {
		return "", ErrNotFound
	}

	return id, nil
}

func GetParentIdByResourcePath(baseUrl, resourcePath, bearerToken string) (string, error) {
	splitPath := strings.Split(resourcePath, "/")

	if splitPath[0] != "Root" {
		return "", ErrPathStartIncorrect
	}

	parentName := splitPath[len(splitPath)-2]

	result, err := PostSearch(baseUrl, parentName, bearerToken)
	if err != nil {
		return "", err
	}

	j, err := unmarshalSearchResponse(result)
	if err != nil {
		return "", err
	}

	parentPath := strings.Join(splitPath[:len(splitPath)-1], "/")

	var count int
	var id string

	for _, c := range j.Groups {
		if strings.TrimSuffix(c.FullPath, "/") == parentPath && c.Name == parentName {
			count++
			id = c.Id
		}
	}

	if count > 1 {
		return "", ErrAmbiguousResult
	} else if count == 0 {
		return "", ErrParentNotFound
	}

	return id, nil
}

func GetValidPaths(baseUrl, resourcePath, bearerToken string) ([]string, error) {
	splitPath := strings.Split(resourcePath, "/")

	if splitPath[0] != "Root" {
		return nil, ErrPathStartIncorrect
	}

	resourceName := splitPath[len(splitPath)-1]
	if resourceName == "" {
		return nil, ErrLastPathComp
	}

	result, err := PostSearch(baseUrl, resourceName, bearerToken)
	if err != nil {
		return nil, err
	}

	j, err := unmarshalSearchResponse(result)
	if err != nil {
		return nil, err
	}

	var paths []string

	for _, c := range j.Credentials {
		splitCredPath := strings.Split(c.Path, "/")
		if strings.Contains(splitCredPath[len(splitCredPath)-1], resourceName) {
			paths = append(paths, fmt.Sprintf("%s%s", c.Path, c.Name))
		}
	}

	for _, g := range j.Groups {
		paths = append(paths, g.FullPath)
	}

	var validPaths []string

	for _, p := range paths {
		if strings.HasPrefix(p, resourcePath) {
			validPaths = append(validPaths, p)
		}
	}

	return deduplicateStrSlice(validPaths), nil
}

func CompletePathFlag(toComplete string) []cobra.Completion {
	baseUrl, bearerToken := LoadConfig()

	paths, err := GetValidPaths(baseUrl, toComplete, bearerToken)
	if err != nil {
		return nil
	}

	completions := make([]cobra.Completion, len(paths))
	for i, p := range paths {
		completions[i] = cobra.Completion(p)
	}

	return completions
}

func DuplicateEntryExists(baseUrl, jsonString, bearerToken string) (bool, error) {
	input, err := UnmarshalEntry(jsonString)
	if err != nil {
		return false, err
	}

	folder, err := GetJsonBody(baseUrl, PathFolders+"/"+input.GroupId, bearerToken)
	if err != nil {
		return false, err
	}

	contents, err := UnmarshalFolderOutput(folder)
	if err != nil {
		return false, err
	}

	for _, entry := range contents.Credentials {
		if entry.Name == input.Name {
			return true, nil
		}
	}

	return false, nil
}

func DuplicateEntryId(baseUrl, jsonString, bearerToken string) (string, error) {
	input, err := UnmarshalEntry(jsonString)
	if err != nil {
		return "", err
	}

	folder, err := GetJsonBody(baseUrl, PathFolders+"/"+input.GroupId, bearerToken)
	if err != nil {
		return "", err
	}

	contents, err := UnmarshalFolderOutput(folder)
	if err != nil {
		return "", err
	}

	for _, entry := range contents.Credentials {
		if entry.Name == input.Name {
			return entry.Id, nil
		}
	}

	return "", nil
}

func DuplicateFolderExists(baseUrl, jsonString, bearerToken string) (bool, error) {
	input, err := UnmarshalFolder(jsonString)
	if err != nil {
		return false, err
	}

	folder, err := GetJsonBody(baseUrl, PathFolders+"/"+input.ParentId, bearerToken)
	if err != nil {
		return false, err
	}

	contents, err := UnmarshalFolderOutput(folder)
	if err != nil {
		return false, err
	}

	for _, folder := range contents.Children {
		if folder.Name == input.Name {
			return true, nil
		}
	}

	return false, nil
}

func DuplicateFolderId(baseUrl, jsonString, bearerToken string) (string, error) {
	input, err := UnmarshalFolder(jsonString)
	if err != nil {
		return "", err
	}

	folder, err := GetJsonBody(baseUrl, PathFolders+"/"+input.ParentId, bearerToken)
	if err != nil {
		return "", err
	}

	contents, err := UnmarshalFolderOutput(folder)
	if err != nil {
		return "", err
	}

	for _, folder := range contents.Children {
		if folder.Name == input.Name {
			return folder.Id, nil
		}
	}

	return "", nil
}
