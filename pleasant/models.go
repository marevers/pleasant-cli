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
