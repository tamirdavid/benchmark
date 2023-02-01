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
* up
    - Flags:
        - size: Determine the size of the files (in MB) that will be generated for testing the upload process. **[Default: 50]** **[Mandatory: False]**
        - iterations:How many files will be created for testing the upload process. **[Default: 30]** **[Mandatory: False]**
        - repo_name: repository the tests will be executed on. **[Default: benchmark-up-tests]** **[Mandatory: False]**
    - Example:
    ```
  $ jf benchmark up
  $ jf benchmark up --size 50 --iterations 5
  $ jf benchmark up --size 50 --iterations 5 --repo_name mytestrepo
  ```
* dl
    - Flags:
        - size: Determine the size of the files (in MB) that will be generated for testing the download process. **[Default: 50]** **[Mandatory: False]**
        - iterations: How many files will be created for testing the download process. **[Default: 30]** **[Mandatory: False]**
        - repo_name: Repository the tests will be executed on. **[Default: benchmark-up-tests]** **[Mandatory: False]**
    - Example:
    ```
  $ jf benchmark dl  
  $ jf benchmark dl --size 50 --iterations 5
  $ jf benchmark dl --size 50 --iterations 5 --repo_name mytestrepo
  ```


## Release Notes
The release notes are available [here](RELEASE.md).
