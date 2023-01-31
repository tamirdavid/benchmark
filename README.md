# benchmark

## About this plugin
This plugin is designed to measure upload and download operations to and from your Artifactory instance. It can be used to perform load tests on your instance to ensure that it can handle the desired level of traffic. The plugin provides metrics on upload and download speeds, as well as the number of successful and failed operations. This information can be used to identify and troubleshoot performance issues, and to optimize the configuration of your instance for optimal performance.

## Installation with JFrog CLI
Installing the latest version:

`$ jf plugin install benchmark`

Installing a specific version:

`$ jf plugin install benchmark@version`

Uninstalling a plugin

`$ jf plugin uninstall benchmark`

## Usage
### Commands
* hello
    - Arguments:
        - addressee - The name of the person you would like to greet.
    - Flags:
        - shout: Makes output uppercase **[Default: false]**
    - Example:
    ```
  $ jf hello-frog hello world --shout
  
  NEW GREETING: HELLO WORLD
  ```

### Environment variables
* HELLO_FROG_GREET_PREFIX - Adds a prefix to every greet **[Default: New greeting: ]**

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
