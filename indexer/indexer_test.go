package indexer

import (
	"encoding/json"
	"testing"
	servicesUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	pkg := NewPackage()
	pkg.PackageID = "pkgID"
	pkg.Version = "version"
	pkg.PackageRevision = "prev"
	pkg.Requires = append(pkg.Requires, "r1")
	pkg.AddSetting("os", "Linux")
	pkg.AddSetting("compiler.version", "gcc")

	b, err := json.Marshal(pkg)
	assert.Nil(t, err)
	assert.Equal(t, `{"package_id":"pkgID","version":"version","package_revision":"prev","settings":{"compiler_version":"gcc","os":"Linux"},"requires":["r1"]}`, string(b))
}

func TestIndexData(t *testing.T) {
	data := IndexData{}
	data.User = "user"
	data.Channel = "channel"
	data.RecipeRevision = "rrev"
	data.Name = "name"
	data.Version = "version"
	data.Description = "the description"
	data.License = "license"
	data.Homepage = "https://homepage.url"
	data.URL = "url"
	data.Topics = "t1,t2"
	data.Requires = append(data.Requires, "r1")
	data.SetForce(true)

	pkg := NewPackage()
	pkg.PackageID = "pkgID"
	pkg.Version = "version"
	pkg.PackageRevision = "prev"
	pkg.Requires = append(pkg.Requires, "r1")
	pkg.AddSetting("os", "Linux")
	pkg.AddSetting("compiler.version", "gcc")
	data.Packages = append(data.Packages, *pkg)

	b, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, `{"user":"user","channel":"channel","recipe_revision":"rrev","name":"name","version":"version","description":"the description","license":"license","homepage":"https://homepage.url","giturl":"url","topics":"t1,t2","requires":["r1"],"packages":[{"package_id":"pkgID","version":"version","package_revision":"prev","settings":{"compiler_version":"gcc","os":"Linux"},"requires":["r1"]}],"force":true,"force_requires":true,"force_settings":true}`, string(b))
}

func TestIndexDataEmptyFields(t *testing.T) {
	data := IndexData{}
	data.User = "user"
	data.Channel = "channel"
	data.RecipeRevision = "rrev"
	data.Name = "name"
	data.Version = "version"
	data.Requires = append(data.Requires, "r1")
	data.SetForce(true)

	b, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, `{"user":"user","channel":"channel","recipe_revision":"rrev","name":"name","version":"version","requires":["r1"],"packages":null,"force":true,"force_requires":true,"force_settings":true}`, string(b))
}

func TestNewFromProperties(t *testing.T) {
	props := []servicesUtils.Property{}
	props = append(props, servicesUtils.Property{Key: "topics", Value: "topic1"})
	props = append(props, servicesUtils.Property{Key: "topics", Value: "topic2"})
	props = append(props, servicesUtils.Property{Key: "settings", Value: "os"})
	props = append(props, servicesUtils.Property{Key: "settings", Value: "arch"})
	props = append(props, servicesUtils.Property{Key: "description", Value: "B2 makes it easy to build C++ projects, everywhere."})
	props = append(props, servicesUtils.Property{Key: "user"})
	props = append(props, servicesUtils.Property{Key: "deprecated"})
	props = append(props, servicesUtils.Property{Key: "options", Value: "toolset"})

	indexData := NewFromProperties(props)
	assert.Equal(t, "aaa", indexData.Name)
}
