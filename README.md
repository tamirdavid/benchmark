# hello-frog

## About this plugin
This plugin is a template and a functioning example for a basic JFrog CLI plugin. 
This README shows the expected structure of your plugin's README.

## Installation with JFrog CLI
Installing the latest version:

`$ jf plugin install hello-frog`

Installing a specific version:

`$ jf plugin install hello-frog@version`

Uninstalling a plugin

`$ jf plugin uninstall hello-frog`

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
