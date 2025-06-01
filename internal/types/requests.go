package types

type DependencyCheckRequest struct {
	Ecosystem      string `json:"ecosystem"`
	Package        string `json:"package"`
	CurrentVersion string `json:"currentVersion"`
	UpgradeVersion string `json:"upgradeVersion"`
	FileContent    string `json:"fileContent"`
}
