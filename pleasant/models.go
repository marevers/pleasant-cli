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

type ConfigFile struct {
	ServerUrl string `yaml:"serverurl"`
}

type TokenFile struct {
	Token *Token `yaml:"bearertoken"`
}

type Token struct {
	AccessToken string `yaml:"accesstoken"`
	ExpiresAt   int64  `yaml:"expiresat"`
}

type SearchOutput struct {
	Credentials []SearchEntry
	Groups      []SearchGroup
}

type SearchEntry struct {
	Id       string
	Name     string
	Username string
	Url      string
	Notes    string
	Tags     []any
	GroupId  string
	Path     string
}

type SearchGroup struct {
	Id       string
	Name     string
	FullPath string
}

type Tag struct {
	Name string
}

type Entry struct {
	Tags     []Tag
	Id       string
	Name     string
	Username string
	Password string
	Url      string
	Notes    string
	GroupId  string
	Expires  string
}

type Folder struct {
	Children []Entry
	Tags     []Tag
	Id       string
	Name     string
	ParentId string
	Notes    string
	Expires  string
}

type FolderOutput struct {
	Credentials []Entry
	Children    []Folder
}
