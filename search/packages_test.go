package search

import (
	"io"
	"io/ioutil"

	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/assert"
)

type MockRtServicesManagerPackages struct {
	artifactory.EmptyArtifactoryServicesManager
}

func (esm *MockRtServicesManagerPackages) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "testdata/search_packages.json")

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "prefix-")
	fileContent, _ := ioutil.ReadFile(filePath)
	_, _ = tmpFile.Write(fileContent)
	tmpFile.Close()

	reader := content.NewContentReader(tmpFile.Name(), "results")
	return reader, nil
}

func (esm *MockRtServicesManagerPackages) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	if readPath == "repository/_/b2/4.0.0/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.0.0@_/_",
			"revisions": [{
				"revision": "3c07b6a54477e856d429493d01c85636",
				"time": "2020-09-16T14:05:05.965+0000"
			}, {
				"revision": "5918010f58ef4294511ff176ccc236b0",
				"time": "2020-08-17T15:20:47.871+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.0.1/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.0.1@_/_",
			"revisions": [{
				"revision": "fe103dcc7b9fa2226d82f5fb43af1d09",
				"time": "2020-09-16T14:06:23.885+0000"
			}, {
				"revision": "64a94a3e9fe90b33033ec9e00eb036e6",
				"time": "2020-08-17T15:20:52.616+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.2.0/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.2.0@_/_",
			"revisions": [{
				"revision": "efacbfac6ee3561ff07968a372b940af",
				"time": "2020-09-16T14:08:54.728+0000"
			}, {
				"revision": "7987eb34c600c944d8a30ffe090fd013",
				"time": "2020-08-17T15:21:01.757+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.3.0/_/ec8af29b790f5745890470ce4220ed50/package/46f53f156846659bf39ad6675fa0ee8156e859fe/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.3.0@_/_#ec8af29b790f5745890470ce4220ed50:46f53f156846659bf39ad6675fa0ee8156e859fe",
			"revisions": [{
				"revision": "anotherprev",
				"time": "2020-11-01T01:08:43.496+0000"
			}, {
				"revision": "91521b313ac2e32c6306677464116901",
				"time": "2020-11-08T01:08:43.496+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.2.0/_/efacbfac6ee3561ff07968a372b940af/package/46f53f156846659bf39ad6675fa0ee8156e859fe/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.2.0@_/_#efacbfac6ee3561ff07968a372b940af:46f53f156846659bf39ad6675fa0ee8156e859fe",
			"revisions": [{
				"revision": "1d55eb15426c7b4f58fd685a82798f2c",
				"time": "2020-09-16T14:09:30.645+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.2.0/_/efacbfac6ee3561ff07968a372b940af/package/4db1be536558d833e52e862fd84d64d75c2b3656/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.2.0@_/_#efacbfac6ee3561ff07968a372b940af:4db1be536558d833e52e862fd84d64d75c2b3656",
			"revisions": [{
				"revision": "f9d9ecd0a8f306a14ec77b3b14f7284a",
				"time": "2020-09-16T14:09:32.532+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.3.0/_/ec8af29b790f5745890470ce4220ed50/package/4db1be536558d833e52e862fd84d64d75c2b3656/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.3.0@_/_#ec8af29b790f5745890470ce4220ed50:4db1be536558d833e52e862fd84d64d75c2b3656",
			"revisions": [{
				"revision": "675b3df28a8ad03689634e1b4f46187f",
				"time": "2020-11-08T01:08:49.352+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.0.1/_/fe103dcc7b9fa2226d82f5fb43af1d09/package/ca33edce272a279b24f87dc0d4cf5bbdcffbc187/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.0.1@_/_#fe103dcc7b9fa2226d82f5fb43af1d09:ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
			"revisions": [{
				"revision": "a462982b96300ac531d82ce34c84ab60",
				"time": "2020-09-16T14:06:24.595+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.2.0/_/efacbfac6ee3561ff07968a372b940af/package/ca33edce272a279b24f87dc0d4cf5bbdcffbc187/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference": "b2/4.2.0@_/_#efacbfac6ee3561ff07968a372b940af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
			"revisions": [{
				"revision": "6ea9badc4dd235e150326d1460ca61b0",
				"time": "2020-09-16T14:08:55.456+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.0.0/_/3c07b6a54477e856d429493d01c85636/package/46f53f156846659bf39ad6675fa0ee8156e859fe/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.0.0@_/_#3c07b6a54477e856d429493d01c85636:46f53f156846659bf39ad6675fa0ee8156e859fe",
			"revisions" : [ {
			  "revision" : "f62ce7a872642c6b5beb8ae1fed2131b",
			  "time" : "2020-09-16T14:05:49.061+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.0.1/_/fe103dcc7b9fa2226d82f5fb43af1d09/package/46f53f156846659bf39ad6675fa0ee8156e859fe/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.0.1@_/_#fe103dcc7b9fa2226d82f5fb43af1d09:46f53f156846659bf39ad6675fa0ee8156e859fe",
			"revisions" : [ {
			  "revision" : "b9345d42018a28b312bdfcb37fc32f7f",
			  "time" : "2020-09-16T14:07:05.542+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.1.0/_/151655c3ac57c4adcc3681a2bf44e0af/package/46f53f156846659bf39ad6675fa0ee8156e859fe/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.1.0@_/_#151655c3ac57c4adcc3681a2bf44e0af:46f53f156846659bf39ad6675fa0ee8156e859fe",
			"revisions" : [ {
			  "revision" : "1bbc47cdda1c74e1ca04119559845127",
			  "time" : "2020-08-17T15:21:40.975+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.0.0/_/3c07b6a54477e856d429493d01c85636/package/4db1be536558d833e52e862fd84d64d75c2b3656/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.0.0@_/_#3c07b6a54477e856d429493d01c85636:4db1be536558d833e52e862fd84d64d75c2b3656",
			"revisions" : [ {
			  "revision" : "513adef99548254b2b5800a5fc3569c6",
			  "time" : "2020-09-16T14:05:47.123+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.1.0/_/151655c3ac57c4adcc3681a2bf44e0af/package/4db1be536558d833e52e862fd84d64d75c2b3656/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.1.0@_/_#151655c3ac57c4adcc3681a2bf44e0af:4db1be536558d833e52e862fd84d64d75c2b3656",
			"revisions" : [ {
			  "revision" : "482838ec4e86ac4544ba3ce11c0bce59",
			  "time" : "2020-08-17T15:21:27.382+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.1.0/_/151655c3ac57c4adcc3681a2bf44e0af/package/ca33edce272a279b24f87dc0d4cf5bbdcffbc187/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.1.0@_/_#151655c3ac57c4adcc3681a2bf44e0af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
			"revisions" : [ {
			  "revision" : "716c28b4b6c8d29b57c41bc4026f0b09",
			  "time" : "2020-08-17T15:21:05.047+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.3.0/_/ec8af29b790f5745890470ce4220ed50/package/ca33edce272a279b24f87dc0d4cf5bbdcffbc187/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.3.0@_/_#ec8af29b790f5745890470ce4220ed50:ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
			"revisions" : [ {
			  "revision" : "2904158bb9b96db13de732f1c8ca4b64",
			  "time" : "2020-11-08T01:08:54.535+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.0.0/_/3c07b6a54477e856d429493d01c85636/package/ca33edce272a279b24f87dc0d4cf5bbdcffbc187/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.0.0@_/_#3c07b6a54477e856d429493d01c85636:ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
			"revisions" : [ {
			  "revision" : "5f27426c663ab6ef42a368cc5f41be25",
			  "time" : "2020-09-16T14:05:06.655+0000"
			} ]
		  }`)), nil
	} else if readPath == "repository/_/b2/4.0.1/_/fe103dcc7b9fa2226d82f5fb43af1d09/package/4db1be536558d833e52e862fd84d64d75c2b3656/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"packageReference" : "b2/4.0.1@_/_#fe103dcc7b9fa2226d82f5fb43af1d09:4db1be536558d833e52e862fd84d64d75c2b3656",
			"revisions" : [ {
			  "revision" : "61e82452df2c33f90fdaafd84a3fcbb9",
			  "time" : "2020-09-16T14:07:06.263+0000"
			} ]
		  }`)), nil
	}
	log.Info(">>> readPath:", readPath)
	return nil, nil
}

func TestSearchPackages(t *testing.T) {
	servicesManager := MockRtServicesManagerPackages{}
	packages, err := SearchPackages(&servicesManager, "repository", "b2", false, false)
	assert.Nil(t, err)
	assert.Equal(t, 25, len(packages))
}

func TestSearchPackagesLatestRecipes(t *testing.T) {
	servicesManager := MockRtServicesManagerPackages{}
	packages, err := SearchPackages(&servicesManager, "repository", "b2", true, false)
	assert.Nil(t, err)
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].String() < packages[j].String()
	})
	assert.Equal(t, 16, len(packages))
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:46f53f156846659bf39ad6675fa0ee8156e859fe#f62ce7a872642c6b5beb8ae1fed2131b", packages[0].String())
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:4db1be536558d833e52e862fd84d64d75c2b3656#513adef99548254b2b5800a5fc3569c6", packages[1].String())
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#5f27426c663ab6ef42a368cc5f41be25", packages[2].String())

	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:46f53f156846659bf39ad6675fa0ee8156e859fe#b9345d42018a28b312bdfcb37fc32f7f", packages[3].String())
	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:4db1be536558d833e52e862fd84d64d75c2b3656#61e82452df2c33f90fdaafd84a3fcbb9", packages[4].String())
	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#a462982b96300ac531d82ce34c84ab60", packages[5].String())

	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:46f53f156846659bf39ad6675fa0ee8156e859fe#1bbc47cdda1c74e1ca04119559845127", packages[6].String())
	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:4db1be536558d833e52e862fd84d64d75c2b3656#482838ec4e86ac4544ba3ce11c0bce59", packages[7].String())
	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#716c28b4b6c8d29b57c41bc4026f0b09", packages[8].String())

	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:46f53f156846659bf39ad6675fa0ee8156e859fe#1d55eb15426c7b4f58fd685a82798f2c", packages[9].String())
	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:4db1be536558d833e52e862fd84d64d75c2b3656#f9d9ecd0a8f306a14ec77b3b14f7284a", packages[10].String())
	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#6ea9badc4dd235e150326d1460ca61b0", packages[11].String())

	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:46f53f156846659bf39ad6675fa0ee8156e859fe#91521b313ac2e32c6306677464116901", packages[12].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:46f53f156846659bf39ad6675fa0ee8156e859fe#anotherprev", packages[13].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:4db1be536558d833e52e862fd84d64d75c2b3656#675b3df28a8ad03689634e1b4f46187f", packages[14].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#2904158bb9b96db13de732f1c8ca4b64", packages[15].String())
}

func TestSearchPackagesLatestAll(t *testing.T) {
	servicesManager := MockRtServicesManagerPackages{}
	packages, err := SearchPackages(&servicesManager, "repository", "b2", true, true)
	assert.Nil(t, err)
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].String() < packages[j].String()
	})
	assert.Equal(t, 15, len(packages))
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:46f53f156846659bf39ad6675fa0ee8156e859fe#f62ce7a872642c6b5beb8ae1fed2131b", packages[0].String())
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:4db1be536558d833e52e862fd84d64d75c2b3656#513adef99548254b2b5800a5fc3569c6", packages[1].String())
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#5f27426c663ab6ef42a368cc5f41be25", packages[2].String())

	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:46f53f156846659bf39ad6675fa0ee8156e859fe#b9345d42018a28b312bdfcb37fc32f7f", packages[3].String())
	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:4db1be536558d833e52e862fd84d64d75c2b3656#61e82452df2c33f90fdaafd84a3fcbb9", packages[4].String())
	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#a462982b96300ac531d82ce34c84ab60", packages[5].String())

	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:46f53f156846659bf39ad6675fa0ee8156e859fe#1bbc47cdda1c74e1ca04119559845127", packages[6].String())
	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:4db1be536558d833e52e862fd84d64d75c2b3656#482838ec4e86ac4544ba3ce11c0bce59", packages[7].String())
	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#716c28b4b6c8d29b57c41bc4026f0b09", packages[8].String())

	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:46f53f156846659bf39ad6675fa0ee8156e859fe#1d55eb15426c7b4f58fd685a82798f2c", packages[9].String())
	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:4db1be536558d833e52e862fd84d64d75c2b3656#f9d9ecd0a8f306a14ec77b3b14f7284a", packages[10].String())
	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#6ea9badc4dd235e150326d1460ca61b0", packages[11].String())

	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:46f53f156846659bf39ad6675fa0ee8156e859fe#91521b313ac2e32c6306677464116901", packages[12].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:4db1be536558d833e52e862fd84d64d75c2b3656#675b3df28a8ad03689634e1b4f46187f", packages[13].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#2904158bb9b96db13de732f1c8ca4b64", packages[14].String())
}
