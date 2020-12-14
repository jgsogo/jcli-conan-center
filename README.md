# jcli-conan-center

![Test Across Matrix](https://github.com/jgsogo/jcli-conan-center/workflows/Test%20Across%20Matrix/badge.svg?branch=master)
![Lint Go Code](https://github.com/jgsogo/jcli-conan-center/workflows/Lint%20Go%20Code/badge.svg?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/jgsogo/jcli-conan-center.svg)](https://pkg.go.dev/github.com/jgsogo/jcli-conan-center)
[![codecov](https://codecov.io/gh/jgsogo/jcli-conan-center/branch/master/graph/badge.svg)](https://codecov.io/gh/jgsogo/jcli-conan-center)


## About this plugin
This is a JFrog CLI plugin to manage Conan repositories. Some of its commands
expect to find properties associated to packages as of those in Conan Center.


## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg),
it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. CD into the root directory of the cloned project.
5. Run ```go build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.


## Commands

This plugin requires a valid authentication to Artifactory, use JFrog CLI as usual
to provide credentials and configure the connection.

Available commands:

 * Search packages: `search [command options] <repo>`
 * Get properties: `properties [command options] <repo> <reference>`
 * Indexeer JSON call: `index-reference [command options] <repo> <reference>`

**Note.-** Commands are documented using the plugin isolated, to use them within
JFrog CLI just change the `go run main.go` with `jfrog conan-center` after installing
it.

## Search packages: `search [command options] <repo>`

Returns the list of Conan references in a given Artifactory repository

* Arguments:

  * `repo`: Name of the Artifactory repository

* Flags:

  * `--server-id` [Optional]: Artifactory server ID configured using the config
    command. If not specified, the default configured Artifactory server is used.
  * `--ref-name` [Optional]: Name of the Conan reference to search (only the name).
    If not set, it will search for all references.
  * `--packages` [Default: `false`]: If specified, it will retrieve Conan packages
    instead of references.
  * `--only-latest` [Default: `false`]: If specified, it will retrieve only the latest
    revision for packages or recipes.


<details><summary>Example: all references (all revisions) in a repository</summary>
<p>

```
$> go run main.go search conan-center

[Info] Found 171 references.
optional-lite/3.2.0#084d0464901dd0fe38a2bd9ddfb5f1df
optional-lite/3.2.0#54a8db3bf59eda2b62f21758b91473ee
optional-lite/3.2.0#dfe10998d2a51a857e69cc34cb5ff91b
nasm/2.13.01#250720a29c2eaaccf49ea3df06f2772a
nasm/2.13.01#63659723342a256b38af04c8fe0237ce
openssl/1.1.1c#efa6db368062e31e64cad382990970e3
openssl/1.1.1f#cda22c20cbf83946b5313386f88267ab
boost/1.73.0#26a99093ce49eabfb628051ea23f7242
span-lite/new-version-bump#0d114658ddebb6582b1bb1068f9f60d8
...
```
</p>
</details>

<details><summary>Example: all references (latest revision) in a repository</summary>
<p>

```
$> go run main.go search conan-center --only-latest

[Info] Found 95 references.
optional-lite/3.2.0#dfe10998d2a51a857e69cc34cb5ff91b
nasm/2.14#e2ed38224348da1b0e31223aae547690
openssl/1.0.2s#f3ac03b5eb67f428a21444ebdae8d0b6
openssl/1.1.0l#7f3fa5cfcfba31fffa344c71a9795176
...
```
</p>
</details>

<details><summary>Example: Search reference by name (latest revision)</summary>
<p>

```
$> go run main.go search conan-center --ref-name=b2 --only-latest

[Info] Found 5 references.
b2/4.0.0#3c07b6a54477e856d429493d01c85636
b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09
b2/4.2.0#efacbfac6ee3561ff07968a372b940af
b2/4.1.0#87e5e0e1d7eab23643ca941d08aecac7
b2/4.3.0#ec8af29b790f5745890470ce4220ed50
...
```
</p>
</details>

<details><summary>Example: Packages by reference name (latest revision)</summary>
<p>

```
$> go run main.go search conan-center --ref-name=b2 --only-latest --packages

Found 15 packages:
b2/4.0.0#3c07b6a54477e856d429493d01c85636:4db1be536558d833e52e862fd84d64d75c2b3656#513adef99548254b2b5800a5fc3569c6
b2/4.0.0#3c07b6a54477e856d429493d01c85636:46f53f156846659bf39ad6675fa0ee8156e859fe#f62ce7a872642c6b5beb8ae1fed2131b
b2/4.0.0#3c07b6a54477e856d429493d01c85636:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#5f27426c663ab6ef42a368cc5f41be25
b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#a462982b96300ac531d82ce34c84ab60
b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:46f53f156846659bf39ad6675fa0ee8156e859fe#b9345d42018a28b312bdfcb37fc32f7f
b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09:4db1be536558d833e52e862fd84d64d75c2b3656#61e82452df2c33f90fdaafd84a3fcbb9
b2/4.2.0#efacbfac6ee3561ff07968a372b940af:ca33edce272a279b24f87dc0d4cf5bbdcffbc187#6ea9badc4dd235e150326d1460ca61b0
b2/4.2.0#efacbfac6ee3561ff07968a372b940af:46f53f156846659bf39ad6675fa0ee8156e859fe#1d55eb15426c7b4f58fd685a82798f2c
b2/4.2.0#efacbfac6ee3561ff07968a372b940af:4db1be536558d833e52e862fd84d64d75c2b3656#f9d9ecd0a8f306a14ec77b3b14f7284a
...
```
</p>
</details>


## Get properties: `properties [command options] <repo> <reference>`

Returns the properties associated to a given Conan reference in a given Artifactory repository

* Arguments:

  * `repo`: Name of the Artifactory repository
  * `reference`: Conan reference to work with (use v2 style, without trailing @).
    If no revision is given, it will use latest one

* Flags:

  * `--server-id` [Optional]: Artifactory server ID configured using the config
    command. If not specified, the default configured Artifactory server is used.
  * `--packages` [Default: `false`]: If specified, it will retrieve properties
    from packages as well.

<details><summary>Example: Return properties for a given reference</summary>
<p>

```
$> go run main.go properties conan-center b2/4.3.0#ec8af29b790f5745890470ce4220ed50

Reference 'b2/4.3.0#ec8af29b790f5745890470ce4220ed50':
  topics: conan
  topics: builder
  settings: os
  url: https://github.com/conan-io/conan-center-index
  options: use_cxx_env
  license: BSL-1.0
  settings: arch
  topics: boost
  version: 4.3.0
  deprecated: 
  options: toolset
  name: b2
  channel: 
  description: B2 makes it easy to build C++ projects, everywhere.
  topics: installer
  homepage: https://boostorg.github.io/build/
  user:
```
</p>
</details>

<details><summary>Example: Return properties for a reference (latest revision)
and all the packages associated to it</summary>
<p>

```
$> go run main.go properties conan-center b2/4.3.0 --packages

Reference 'b2/4.3.0#ec8af29b790f5745890470ce4220ed50':
  topics: conan
  topics: builder
  settings: os
  url: https://github.com/conan-io/conan-center-index
  options: use_cxx_env
  license: BSL-1.0
  settings: arch
  topics: boost
  version: 4.3.0
  deprecated: 
  options: toolset
  name: b2
  channel: 
  description: B2 makes it easy to build C++ projects, everywhere.
  topics: installer
  homepage: https://boostorg.github.io/build/
  user: 
Package 'b2/4.3.0#ec8af29b790f5745890470ce4220ed50:46f53f156846659bf39ad6675fa0ee8156e859fe#91521b313ac2e32c6306677464116901':
  settings: build_type=Release
  topics: conan
  ...
Package 'b2/4.3.0#ec8af29b790f5745890470ce4220ed50:4db1be536558d833e52e862fd84d64d75c2b3656#675b3df28a8ad03689634e1b4f46187f':
  topics: boost
  settings: os=Linux
  ...
```
</p>
</details>


## Indexeer JSON call: `index-reference [command options] <repo> <reference>`

Returns the properties associated to a given Conan reference in a given Artifactory repository

* Arguments:

  * `repo`: Name of the Artifactory repository
  * `reference`: Conan reference to work with (use v2 style, without trailing @).
    If no revision is given, it will use latest one.

* Flags:

  * `--server-id` [Optional]: Artifactory server ID configured using the config
    command. If not specified, the default configured Artifactory server is used.
  * `--force` [Default: `false`]: Value for argument `force` in the indexer call.

<details><summary>Example: Indexer call for a reference</summary>
<p>

```json
$> go run main.go index-reference conan-center b2/4.3.0 --force

{
        "user": "",
        "channel": "",
        "recipe_revision": "ec8af29b790f5745890470ce4220ed50",
        "name": "b2",
        "version": "4.3.0",
        "description": "B2 makes it easy to build C++ projects, everywhere.",
        "license": "BSL-1.0",
        "homepage": "https://boostorg.github.io/build/",
        "giturl": "https://github.com/conan-io/conan-center-index",
        "topics": "conan,builder,boost,installer",
        "requires": null,
        "packages": [
                {
                        "package_id": "46f53f156846659bf39ad6675fa0ee8156e859fe",
                        "version": "4.3.0",
                        "package_revision": "91521b313ac2e32c6306677464116901",
                        "settings": {
                                "arch": "x86_64",
                                "arch_build": "x86_64",
                                "build_type": "Release",
                                "compiler": "apple-clang",
                                "compiler_libcxx": "libc++",
                                "compiler_version": "10.0",
                                "os": "Macos",
                                "os_build": "Macos"
                        }
                },
                {
                        "package_id": "4db1be536558d833e52e862fd84d64d75c2b3656",
                        "version": "4.3.0",
                        "package_revision": "675b3df28a8ad03689634e1b4f46187f",
                        "settings": {
                                "arch": "x86_64",
                                "arch_build": "x86_64",
                                "build_type": "Release",
                                "compiler": "gcc",
                                "compiler_libcxx": "libstdc++",
                                "compiler_version": "4.9",
                                "os": "Linux",
                                "os_build": "Linux"
                        }
                },
                {
                        "package_id": "ca33edce272a279b24f87dc0d4cf5bbdcffbc187",
                        "version": "4.3.0",
                        "package_revision": "2904158bb9b96db13de732f1c8ca4b64",
                        "settings": {
                                "arch": "x86_64",
                                "arch_build": "x86_64",
                                "build_type": "Release",
                                "compiler": "Visual Studio",
                                "compiler_runtime": "MT",
                                "compiler_version": "14",
                                "os": "Windows",
                                "os_build": "Windows"
                        }
                }
        ],
        "force": true,
        "force_requires": true,
        "force_settings": true
}
```
</p>
</details>


## Additional info
Work in progress.

## Release Notes
The release notes are available [here](RELEASE.md).
