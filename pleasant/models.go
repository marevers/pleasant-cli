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
	Credentials []Credential
	Groups      []Group
}

type Credential struct {
	Id       string
	Name     string
	Username string
	Url      string
	Notes    string
	Tags     []string
	GroupId  string
	Path     string
}

type Group struct {
	Id       string
	Name     string
	FullPath string
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
	Expires                 string
}

type FolderInput struct {
	CustomUserFields        map[string]string
	CustomApplicationFields map[string]string
	Children                []map[string]string
	Credentials             []map[string]string
	Tags                    []string
	Name                    string
	ParentId                string
	Notes                   string
	Expires                 string
}

type FolderOutput struct {
	Credentials []Credential
	Children    []Group
}
