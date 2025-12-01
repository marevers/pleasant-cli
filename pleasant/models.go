package pleasant

type ConfigFile struct {
	ServerUrl string `yaml:"serverurl"`
	Timeout   int    `yaml:"timeout"`
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
	Tags     []Tag  `json:"Tags,omitempty"`
	Id       string `json:"Id,omitempty"`
	Name     string `json:"Name,omitempty"`
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
	Url      string `json:"Url,omitempty"`
	Notes    string `json:"Notes,omitempty"`
	GroupId  string `json:"GroupId,omitempty"`
	Expires  string `json:"Expires,omitempty"`
}

type Folder struct {
	Children []Entry `json:"Children,omitempty"`
	Tags     []Tag   `json:"Tags,omitempty"`
	Id       string  `json:"Id,omitempty"`
	Name     string  `json:"Name,omitempty"`
	ParentId string  `json:"ParentId,omitempty"`
	Notes    string  `json:"Notes,omitempty"`
	Expires  string  `json:"Expires,omitempty"`
}

type FolderOutput struct {
	Credentials []Entry
	Children    []Folder
}
