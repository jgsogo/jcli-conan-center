# jcli-conan-center

## About this plugin
This is a JFrog CLI plugin to manage Conan repositories. Some of its commands
expect to find properties associated to packages as of those in Conan Center.

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install conan-center`

Installing a specific version:

`$ jfrog plugin install conan-center@version`


## Usage
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
