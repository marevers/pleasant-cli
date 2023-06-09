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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type prerequisite struct {
	Message         string
	PrerequisiteMet bool
}

func CheckPrerequisites(prereq ...*prerequisite) bool {
	for _, p := range prereq {
		if !p.PrerequisiteMet {
			fmt.Println(p.Message)

			return false
		}
	}

	return true
}

func IsTokenValid() *prerequisite {
	b := time.Now().Unix() <= viper.GetInt64("bearertoken.expiresat")

	pr := &prerequisite{
		Message:         "Token is expired or not present. Please log in (again) with 'pleasant-cli login'.",
		PrerequisiteMet: b,
	}

	return pr
}

func IsServerUrlSet() *prerequisite {
	b := viper.IsSet("serverurl")

	pr := &prerequisite{
		Message:         "Server URL is not set. Please set it with 'pleasant-cli config serverurl <SERVER URL>'.",
		PrerequisiteMet: b,
	}

	return pr
}

func IsOneOfRequiredFlagsSet(cmd *cobra.Command, flags ...string) *prerequisite {
	var count int

	for _, f := range flags {
		if cmd.Flags().Changed(f) {
			count++
		}
	}

	b := count > 0

	pr := &prerequisite{
		Message:         "At least one of the following flags must be set: " + strings.Join(flags, ", "),
		PrerequisiteMet: b,
	}

	return pr
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

func WriteConfigFile(file, serverUrl string) error {
	t := &ConfigFile{
		ServerUrl: serverUrl,
	}

	b, err := yaml.Marshal(t)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, b, 0666)
	if err != nil {
		return err
	}

	return nil
}

func WriteTokenFile(file, accessToken string, expiresAt int64) error {
	t := &TokenFile{
		Token: &Token{
			AccessToken: accessToken,
			ExpiresAt:   expiresAt,
		},
	}

	b, err := yaml.Marshal(t)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, b, 0666)
	if err != nil {
		return err
	}

	return nil
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
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
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
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
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
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, generateError(res.StatusCode)
	}

	return res, nil
}

func unmarshalSearchResponse(jsonString string) (*SearchOutput, error) {
	sr := &SearchOutput{}

	err := json.Unmarshal([]byte(jsonString), sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func UnmarshalEntry(jsonString string) (*Entry, error) {
	ei := &Entry{}

	err := json.Unmarshal([]byte(jsonString), ei)
	if err != nil {
		return nil, err
	}

	return ei, nil
}

func MarshalEntry(entry *Entry) (string, error) {
	b, err := json.Marshal(entry)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func UnmarshalFolder(jsonString string) (*Folder, error) {
	fi := &Folder{}

	err := json.Unmarshal([]byte(jsonString), fi)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

func MarshalFolder(folder *Folder) (string, error) {
	b, err := json.Marshal(folder)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func UnmarshalFolderOutput(jsonString string) (*FolderOutput, error) {
	fo := &FolderOutput{}

	err := json.Unmarshal([]byte(jsonString), fo)
	if err != nil {
		return nil, err
	}

	return fo, nil
}

func PathAndNameMatching(resourcePath, name string) bool {
	s := strings.Split(resourcePath, "/")
	return s[len(s)-1] == name
}
