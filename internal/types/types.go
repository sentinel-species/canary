package types

type Dependency struct {
	Id      int     `json:"id"`
	Package Package `json:"package"`
}

type Package struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Ecosystem    string    `json:"ecosystem"`
	Dependencies []Package `json:"dependencies"`
}

type UpgradeResponse struct {
	Id          string `json:"id" db:"id"`
	PackageName string `json:"package_name" db:"package_name"`
	FixVersion  string `json:"fix_version" db:"fix_version"`
	EnvName     string `json:"environment_name" db:"environment_name"`
}
