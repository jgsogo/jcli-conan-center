// Package indexer contains the types and functions related to the ConanCenter indexer
package indexer

import "strings"

// Package contains information about a single packageID inside the `IndexData`
type Package struct {
	PackageID       string `json:"package_id"`
	Version         string `json:"version"`
	PackageRevision string `json:"package_revision"`
	Settings map[string]string `json:"settings"`
	Requires []string          `json:"requires"`
}

// NewPackage creates a `Package` instance and initializes its members
func NewPackage() *Package {
	return &Package{Settings: make(map[string]string)}
}

// AddSetting add the key-value pair for a settings, taking into account some key transformations.
func (pkg *Package) AddSetting(key string, value string) {
	key = strings.ReplaceAll(key, ".", "_")
	pkg.Settings[key] = value
}

// IndexData is the structure with all the information.
type IndexData struct {
	User           string `json:"user"`
	Channel        string `json:"channel"`
	RecipeRevision string `json:"recipe_revision"`

	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	License     string `json:"license,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
	URL         string `json:"giturl,omitempty"`
	Topics      string `json:"topics,omitempty"`

	Requires []string  `json:"requires"`
	Packages []Package `json:"packages"`

	Force         bool `json:"force"`
	ForceRequires bool `json:"force_requires"`
	ForceSettings bool `json:"force_settings"`
}

// SetForce sets the value of the `force` attribute and related ones.
func (data *IndexData) SetForce(value bool) {
	data.Force = value
	data.ForceRequires = value
	data.ForceSettings = value
}
