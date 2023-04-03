/*
Copyright © 2023 Martijn Evers

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pleasant

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

func CheckPrerequisites(checks ...func() (string, bool)) bool {
	for _, c := range checks {
		s, ok := c()

		if !ok {
			fmt.Println(s)

			return false
		}
	}

	return true
}

func IsTokenValid() (string, bool) {
	s := "Token is expired or not present. Please log in (again) with 'pleasant-cli login'."
	b := time.Now().Unix() <= viper.GetInt64("bearertoken.expiresat")
	return s, b
}

func IsServerUrlSet() (string, bool) {
	s := "Server URL is not set. Please set it with 'pleasant-cli config serverurl <SERVER URL>'."
	b := viper.IsSet("serverurl")
	return s, b
}

func StringPrompt(label string) string {
	var s string

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stderr, label+" ")

		s, _ = r.ReadString('\n')

		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}

func PasswordPrompt(label string) string {
	var s string

	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))

		s = string(b)

		if s != "" {
			break
		}
	}

	return s
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: 20 * time.Second,
	}
}

func getRequest(baseUrl, path, bearerToken string) (*http.Response, error) {
	method := "GET"

	req, err := http.NewRequest(method, baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}

	client := newHttpClient()

	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func decodeJsonBody(body io.ReadCloser, target any) error {
	err := json.NewDecoder(body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func postRequestForm(baseUrl, path string, urlValues url.Values) (*http.Response, error) {
	method := "POST"

	payload := strings.NewReader(urlValues.Encode())

	req, err := http.NewRequest(method, baseUrl+path, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := newHttpClient()

	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func postRequestJsonString(baseUrl, path, jsonString, bearerToken string) (*http.Response, error) {
	method := "POST"

	payload := []byte(jsonString)

	req, err := http.NewRequest(method, baseUrl+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := newHttpClient()

	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

type SearchResponse struct {
	Credentials []Credential
	Groups      []Group
}

type Credential struct {
	Id       string
	Name     string
	Username string
	Url      string
	Notes    string
	Tags     string
	GroupId  string
	Path     string
}

type Group struct {
	Id       string
	Name     string
	FullPath string
}

func unmarshalSearchResponse(jsonString string) (*SearchResponse, error) {
	sr := &SearchResponse{}

	err := json.Unmarshal([]byte(jsonString), sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

type EntryInput struct {
	CustomUserFields        map[string]string
	CustomApplicationFields map[string]string
	Tags                    []string
	Name                    string
	Username                string
	Password                string
	Url                     string
	Notes                   string
	GroupId                 string
}

func UnmarshalEntryInput(jsonString string) (*EntryInput, error) {
	ei := &EntryInput{}

	err := json.Unmarshal([]byte(jsonString), ei)
	if err != nil {
		return nil, err
	}

	return ei, nil
}

func MarshalEntryInput(entryInput *EntryInput) (string, error) {
	b, err := json.Marshal(entryInput)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func PathAndNameMatching(resourcePath, name string) bool {
	s := strings.Split(resourcePath, "/")
	return s[len(s)-1] == name
}
