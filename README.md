# jcli-conan-center

![Test Across Matrix](https://github.com/jgsogo/jcli-conan-center/workflows/Test%20Across%20Matrix/badge.svg?branch=master)
![Lint Go Code](https://github.com/jgsogo/jcli-conan-center/workflows/Lint%20Go%20Code/badge.svg?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/jgsogo/jcli-conan-center.svg)](https://pkg.go.dev/github.com/jgsogo/jcli-conan-center)
[![codecov](https://codecov.io/gh/jgsogo/jcli-conan-center/branch/master/graph/badge.svg)](https://codecov.io/gh/jgsogo/jcli-conan-center)

## About this plugin
This is a JFrog CLI plugin to manage Conan repositories. Some of its commands
expect to find properties associated to packages as of those in Conan Center.

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install conan-center`

Installing a specific version:

`$ jfrog plugin install conan-center@version`


## Commands

This plugin requires a valid authentication to Artifactory, use JFrog CLI as usual
to provide credentials and configure the connection.

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



### Commands
* hello
    - Arguments:
        - addressee - The name of the person you would like to greet.
    - Flags:
        - shout: Makes output uppercase **[Default: false]**
        - repeat: Greets multiple times **[Default: 1]**
    - Example:
    ```
  $ jfrog hello-frog hello world --shout --repeat=2
  
  NEW GREETING: HELLO WORLD!
  NEW GREETING: HELLO WORLD!
  ```

### Environment variables
* HELLO_FROG_GREET_PREFIX - Adds a prefix to every greet **[Default: New greeting: ]**

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
